package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestCalculateRandomDelay(t *testing.T) {
	app := new(application)
	app.seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	tests := []struct {
		subTestName  string
		minDelay     int
		maxDelay     int
		errorPresent bool
	}{
		{
			subTestName:  "equal parameters",
			minDelay:     0,
			maxDelay:     0,
			errorPresent: true,
		},
		{
			subTestName:  "minDelay bigger than maxDelay",
			minDelay:     2,
			maxDelay:     1,
			errorPresent: true,
		},
		{
			subTestName:  "minDelay smaller than maxDelay",
			minDelay:     1,
			maxDelay:     2,
			errorPresent: false,
		},
		{
			subTestName:  "minDelay smaller than 0",
			minDelay:     -1,
			maxDelay:     2,
			errorPresent: true,
		},
		{
			subTestName:  "maxDelay smaller than 0",
			minDelay:     1,
			maxDelay:     -2,
			errorPresent: true,
		},
		{
			subTestName:  "minDelay and maxDelay smaller than 0",
			minDelay:     -1,
			maxDelay:     -2,
			errorPresent: true,
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		t.Run(tt.subTestName, func(t *testing.T) {
			_, err := app.calculateRandomDelay(tt.minDelay, tt.maxDelay)
			// No error happened, but an error should have happened.
			if err == nil && tt.errorPresent {
				t.Errorf("No error happened, but an error should have happened.")
			}
			// An error happened, but no error should have happened.
			if err != nil && !(tt.errorPresent) {
				t.Errorf("An error happened, but no error should have happened.")
			}
		})

	}
}

func TestParseRawData(t *testing.T) {

	tests := []struct {
		subTestName string
		entity      Entity
		want        *OrganizationDocument
	}{
		{
			subTestName: "Parse entity and compare with expected OrganizationDocument struct",
			entity: Entity{Uuid: "1", Properties: Properties{
				FoundedOn:             map[string]string{"value": "2022-08-03", "precision": "day"},
				Website:               map[string]string{"value": "www.blub1.ch"},
				Identifier:            map[string]string{"entity_def_id": "organization", "value": "Blub.ai"},
				NumFounders:           2,
				FounderIdentifiers:    []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "linus-torvald", Name: "Linus Torvald"}, {Uuid: "2a", EntityDefId: "person", Permalink: "james-bond", Name: "James Bond"}},
				Description:           "long - blub1 does blub2",
				Linkedin:              map[string]string{"value": "https://www.linkedin.com/company/blub/"},
				Facebook:              map[string]string{"value": "https://www.facebook.com/company/blub/"},
				ShortDescription:      "short - blub1 does blub2",
				NumInvestors:          1,
				OperatingStatus:       "active",
				NumEmployeesEnum:      "c_00001_00010",
				FundingTotal:          Funding{ValueUSD: 1000},
				FundingStage:          "seed",
				NumFundingRounds:      1,
				LastEquityFundingType: "seed",
				InvestorIdentifiers:   []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "bill-gates", Name: "Bill Gates"}, {Uuid: "2a", EntityDefId: "person", Permalink: "steve-jobs", Name: "Steve Jobs"}},
				LastFundingTotal:      Funding{ValueUSD: 1000},
				LastFundingType:       "seed",
				LastFundingAt:         "2022-01-21",
				Categories:            []Category{{Uuid: "2a", EntityDefId: "category", Name: "Machine Learning"}, {Uuid: "2b", EntityDefId: "category", Name: "Online Grocery"}},
				Locations:             []Location{{Uuid: "1a", LocationType: "city", Name: "Basel"}, {Uuid: "2a", LocationType: "region", Name: "Basel City"}, {Uuid: "3a", LocationType: "country", Name: "Switzerland"}},
				NumTrademarkReg:       1,
				NumPatentGrant:        2,
				NumOfTechUsed:         3,
				NumOfArticles:         4,
				ContactEmail:          "hello@blub.ch",
				NumVisitsLastMonth:    69,
				NumVisitPerPageviews:  1.4,
				BounceRate:            1.2,
				VisitDuration:         420,
			}},
			want: &OrganizationDocument{
				Uuid:                  "1",
				Timestamp:             time.Now(),
				EntityDefId:           "organization",
				OrganizationName:      "Blub.ai",
				Description:           "long - blub1 does blub2",
				ShortDescription:      "short - blub1 does blub2",
				FundingStage:          "seed",
				FoundedOn:             time.Date(2022, time.Month(8), 3, 0, 0, 0, 0, time.UTC),
				OperatingStatus:       "active",
				Website:               "www.blub1.ch",
				Linkedin:              "https://www.linkedin.com/company/blub/",
				Facebook:              "https://www.facebook.com/company/blub/",
				Industries:            []Category{{Uuid: "2a", EntityDefId: "category", Name: "Machine Learning"}, {Uuid: "2b", EntityDefId: "category", Name: "Online Grocery"}},
				City:                  "Basel",
				Country:               "Switzerland",
				ContactEmail:          "hello@blub.ch",
				NumFounders:           2,
				NumEmployeesEnum:      "1-10",
				FounderIdentifiers:    []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "linus-torvald", Name: "Linus Torvald"}, {Uuid: "2a", EntityDefId: "person", Permalink: "james-bond", Name: "James Bond"}},
				NumTrademarkReg:       1,
				NumPatentGrant:        2,
				NumOfTechUsed:         3,
				NumOfArticles:         4,
				SemRush:               SemRush{NumVisitsLastMonth: 69, VisitDuration: 420, NumVisitPerPageviews: 1.4, BounceRate: 1.2},
				NumInvestors:          1,
				FundingTotal:          1000,
				NumFundingRounds:      1,
				LastEquityFundingType: "seed",
				LastFundingType:       "seed",
				LastFundingTotal:      1000,
				LastFundingAt:         time.Date(2022, time.Month(1), 21, 0, 0, 0, 0, time.UTC),
				InvestorIdentifiers:   []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "bill-gates", Name: "Bill Gates"}, {Uuid: "2a", EntityDefId: "person", Permalink: "steve-jobs", Name: "Steve Jobs"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.subTestName, func(t *testing.T) {
			initialOrganizationDocument := &OrganizationDocument{}

			initialOrganizationDocument.parseRawData(tt.entity)

			// adds the current timestamp to tt.want to not deal with mocking the time.Now().
			tt.want.Timestamp = initialOrganizationDocument.Timestamp

			// compares the 2 structs after parsing for deep equality.
			if !(reflect.DeepEqual(*initialOrganizationDocument, *tt.want)) {
				t.Errorf("An error: happened. %v and %v are not the same struct, but an error should have happened.", initialOrganizationDocument, tt.want)
			}

		})

	}
}

func TestDecodeBody(t *testing.T) {

	// creates test output of API request
	testOutputStruct := DataContainer{Count: 2, Entities: []Entity{{Uuid: "1a", Properties: Properties{
		FoundedOn:             map[string]string{"value": "2022-08-03", "precision": "day"},
		Website:               map[string]string{"value": "www.blub1.ch"},
		Identifier:            map[string]string{"entity_def_id": "organization", "value": "Blub.ai"},
		NumFounders:           2,
		FounderIdentifiers:    []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "linus-torvald", Name: "Linus Torvald"}, {Uuid: "2a", EntityDefId: "person", Permalink: "james-bond", Name: "James Bond"}},
		Description:           "long - blub1 does blub2",
		Linkedin:              map[string]string{"value": "https://www.linkedin.com/company/blub/"},
		Facebook:              map[string]string{"value": "https://www.facebook.com/company/blub/"},
		ShortDescription:      "short - blub1 does blub2",
		NumInvestors:          1,
		OperatingStatus:       "active",
		NumEmployeesEnum:      "c_00001_00010",
		FundingTotal:          Funding{ValueUSD: 1000},
		FundingStage:          "seed",
		NumFundingRounds:      1,
		LastEquityFundingType: "seed",
		InvestorIdentifiers:   []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "bill-gates", Name: "Bill Gates"}, {Uuid: "2a", EntityDefId: "person", Permalink: "steve-jobs", Name: "Steve Jobs"}},
		LastFundingTotal:      Funding{ValueUSD: 1000},
		LastFundingType:       "seed",
		LastFundingAt:         "2022-01-21",
		Categories:            []Category{{Uuid: "2a", EntityDefId: "category", Name: "Machine Learning"}, {Uuid: "2b", EntityDefId: "category", Name: "Online Grocery"}},
		Locations:             []Location{{Uuid: "1a", LocationType: "city", Name: "Basel"}, {Uuid: "2a", LocationType: "region", Name: "Basel City"}, {Uuid: "3a", LocationType: "country", Name: "Switzerland"}},
		NumTrademarkReg:       1,
		NumPatentGrant:        2,
		NumOfTechUsed:         3,
		NumOfArticles:         4,
		ContactEmail:          "hello@blub.ch",
		NumVisitsLastMonth:    69,
		NumVisitPerPageviews:  1.4,
		BounceRate:            1.2,
		VisitDuration:         420,
	}}, {
		Uuid: "2a",
		Properties: Properties{
			FoundedOn:             map[string]string{"value": "2022-08-03", "precision": "day"},
			Website:               map[string]string{"value": "www.blub1.ch"},
			Identifier:            map[string]string{"entity_def_id": "organization", "value": "Blub.ai"},
			NumFounders:           2,
			FounderIdentifiers:    []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "linus-torvald", Name: "Linus Torvald"}, {Uuid: "2a", EntityDefId: "person", Permalink: "james-bond", Name: "James Bond"}},
			Description:           "long - blub1 does blub2",
			Linkedin:              map[string]string{"value": "https://www.linkedin.com/company/blub/"},
			Facebook:              map[string]string{"value": "https://www.facebook.com/company/blub/"},
			ShortDescription:      "short - blub1 does blub2",
			NumInvestors:          1,
			OperatingStatus:       "active",
			NumEmployeesEnum:      "c_00001_00010",
			FundingTotal:          Funding{ValueUSD: 1000},
			FundingStage:          "seed",
			NumFundingRounds:      1,
			LastEquityFundingType: "seed",
			InvestorIdentifiers:   []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "bill-gates", Name: "Bill Gates"}, {Uuid: "2a", EntityDefId: "person", Permalink: "steve-jobs", Name: "Steve Jobs"}},
			LastFundingTotal:      Funding{ValueUSD: 1000},
			LastFundingType:       "seed",
			LastFundingAt:         "2022-01-21",
			Categories:            []Category{{Uuid: "2a", EntityDefId: "category", Name: "Machine Learning"}, {Uuid: "2b", EntityDefId: "category", Name: "Online Grocery"}},
			Locations:             []Location{{Uuid: "1a", LocationType: "city", Name: "Basel"}, {Uuid: "2a", LocationType: "region", Name: "Basel City"}, {Uuid: "3a", LocationType: "country", Name: "Switzerland"}},
			NumTrademarkReg:       1,
			NumPatentGrant:        2,
			NumOfTechUsed:         3,
			NumOfArticles:         4,
			ContactEmail:          "hello@blub.ch",
			NumVisitsLastMonth:    69,
			NumVisitPerPageviews:  1.4,
			BounceRate:            1.2,
			VisitDuration:         420,
		},
	},
	}}

	// casts test output struct to slice of bytes
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(testOutputStruct)
	byteOutput := reqBodyBytes.Bytes()

	tests := []struct {
		subTestName string
		paylod      []byte
		want        *[]OrganizationDocument
	}{{
		subTestName: "Decodes API request body and adds individual entities to structured slice of OrganizationDocument",
		paylod:      []byte{},
		want: &[]OrganizationDocument{{
			Uuid:                  "1a",
			Timestamp:             time.Now(),
			EntityDefId:           "organization",
			OrganizationName:      "Blub.ai",
			Description:           "long - blub1 does blub2",
			ShortDescription:      "short - blub1 does blub2",
			FundingStage:          "seed",
			FoundedOn:             time.Date(2022, time.Month(8), 3, 0, 0, 0, 0, time.UTC),
			OperatingStatus:       "active",
			Website:               "www.blub1.ch",
			Linkedin:              "https://www.linkedin.com/company/blub/",
			Facebook:              "https://www.facebook.com/company/blub/",
			Industries:            []Category{{Uuid: "2a", EntityDefId: "category", Name: "Machine Learning"}, {Uuid: "2b", EntityDefId: "category", Name: "Online Grocery"}},
			City:                  "Basel",
			Country:               "Switzerland",
			ContactEmail:          "hello@blub.ch",
			NumFounders:           2,
			NumEmployeesEnum:      "1-10",
			FounderIdentifiers:    []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "linus-torvald", Name: "Linus Torvald"}, {Uuid: "2a", EntityDefId: "person", Permalink: "james-bond", Name: "James Bond"}},
			NumTrademarkReg:       1,
			NumPatentGrant:        2,
			NumOfTechUsed:         3,
			NumOfArticles:         4,
			SemRush:               SemRush{NumVisitsLastMonth: 69, VisitDuration: 420, NumVisitPerPageviews: 1.4, BounceRate: 1.2},
			NumInvestors:          1,
			FundingTotal:          1000,
			NumFundingRounds:      1,
			LastEquityFundingType: "seed",
			LastFundingType:       "seed",
			LastFundingTotal:      1000,
			LastFundingAt:         time.Date(2022, time.Month(1), 21, 0, 0, 0, 0, time.UTC),
			InvestorIdentifiers:   []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "bill-gates", Name: "Bill Gates"}, {Uuid: "2a", EntityDefId: "person", Permalink: "steve-jobs", Name: "Steve Jobs"}},
		},
			{
				Uuid:                  "2a",
				Timestamp:             time.Now(),
				EntityDefId:           "organization",
				OrganizationName:      "Blub.ai",
				Description:           "long - blub1 does blub2",
				ShortDescription:      "short - blub1 does blub2",
				FundingStage:          "seed",
				FoundedOn:             time.Date(2022, time.Month(8), 3, 0, 0, 0, 0, time.UTC),
				OperatingStatus:       "active",
				Website:               "www.blub1.ch",
				Linkedin:              "https://www.linkedin.com/company/blub/",
				Facebook:              "https://www.facebook.com/company/blub/",
				Industries:            []Category{{Uuid: "2a", EntityDefId: "category", Name: "Machine Learning"}, {Uuid: "2b", EntityDefId: "category", Name: "Online Grocery"}},
				City:                  "Basel",
				Country:               "Switzerland",
				ContactEmail:          "hello@blub.ch",
				NumFounders:           2,
				NumEmployeesEnum:      "1-10",
				FounderIdentifiers:    []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "linus-torvald", Name: "Linus Torvald"}, {Uuid: "2a", EntityDefId: "person", Permalink: "james-bond", Name: "James Bond"}},
				NumTrademarkReg:       1,
				NumPatentGrant:        2,
				NumOfTechUsed:         3,
				NumOfArticles:         4,
				SemRush:               SemRush{NumVisitsLastMonth: 69, VisitDuration: 420, NumVisitPerPageviews: 1.4, BounceRate: 1.2},
				NumInvestors:          1,
				FundingTotal:          1000,
				NumFundingRounds:      1,
				LastEquityFundingType: "seed",
				LastFundingType:       "seed",
				LastFundingTotal:      1000,
				LastFundingAt:         time.Date(2022, time.Month(1), 21, 0, 0, 0, 0, time.UTC),
				InvestorIdentifiers:   []Person{{Uuid: "1a", EntityDefId: "person", Permalink: "bill-gates", Name: "Bill Gates"}, {Uuid: "2a", EntityDefId: "person", Permalink: "steve-jobs", Name: "Steve Jobs"}},
			},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.subTestName, func(t *testing.T) {
			initialOrganizationDocumentSlice := &[]OrganizationDocument{}
			decodeBody(byteOutput, initialOrganizationDocumentSlice)

			// adds the current timestamp to first and second index of pointer tt.want slice to not deal with mocking the time.Now().
			(*tt.want)[0].Timestamp = (*initialOrganizationDocumentSlice)[0].Timestamp
			(*tt.want)[1].Timestamp = (*initialOrganizationDocumentSlice)[1].Timestamp

			// compares the 2 slices of structs after parsing for deep equality.
			if !(reflect.DeepEqual(*initialOrganizationDocumentSlice, *tt.want)) {
				t.Errorf("An error: happened. %v and %v are not the same struct, but an error should have happened.", initialOrganizationDocumentSlice, tt.want)
			}

		})

	}
}
