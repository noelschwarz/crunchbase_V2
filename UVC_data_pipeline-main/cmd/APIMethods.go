package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// extract, extract data from the Crunchbase API. Parameters: lastUUID, if empty
// string start the request from the beginning, if not use UUID of last element.
// limit, if true limit elements in server response to 1, if false limit is
// set to 1000.
// Output []byte with response body.
func (app *application) extract(lastUUID string, limit bool) ([]byte, error) {

	// Custom CB's Searched Save URL that got randomly selected
	// and serves as basis for CB specific payload queries.
	url := app.cbCustomHeader.UrlReferer

	// Number of entities being requested.
	limitElements := 1000
	if limit {
		limitElements = 1
	}
	// UVC data team search (CB) with dynamic lastUUID parameter.
	payloadString := fmt.Sprintf(`{"field_ids":["identifier","operating_status","founded_on","ipo_status","diversity_spotlights","location_identifiers","categories","description","last_funding_type","investor_identifiers","last_funding_at","funding_total","funding_stage","investor_type","last_equity_funding_type","last_funding_total","num_funding_rounds","num_lead_investors","num_investors","semrush_visits_latest_month","semrush_visits_latest_6_months_avg","semrush_visits_mom_pct","semrush_visit_duration","semrush_visit_duration_mom_pct","semrush_visit_pageviews","semrush_visit_pageview_mom_pct","semrush_bounce_rate","semrush_bounce_rate_mom_pct","semrush_global_rank","semrush_global_rank_mom","semrush_global_rank_mom_pct","apptopia_total_apps","apptopia_total_downloads","num_founders","founder_identifiers","num_employees_enum","investor_stage","website","linkedin","num_articles","hub_tags","twitter","facebook","short_description","contact_email","last_key_employee_change_date","last_layoff_date","num_event_appearances","rank_org_company","num_contacts","num_private_contacts","builtwith_num_technologies_used","siftery_num_products","ipqwery_num_patent_granted","ipqwery_num_trademark_registered","private_tags","num_private_notes"],"order":[{"field_id":"founded_on","sort":"desc"}],"query":[{"type":"predicate","field_id":"founded_on","operator_id":"gte","include_nulls":false,"values":["2020"]},{"type":"predicate","field_id":"operating_status","operator_id":"includes","include_nulls":false},{"type":"predicate","field_id":"location_identifiers","operator_id":"includes","values":["b25caef9-a1b8-3a5d-6232-93b2dfb6a1d1","6106f5dc-823e-5da8-40d7-51612c0b2c4e"]},{"type":"predicate","field_id":"funding_stage","operator_id":"includes","values":["seed","early_stage_venture","late_stage_venture"]}],"field_aggregators":[],"collection_id":"organization.companies","limit":%d, "after_id": "%s"}`, limitElements, lastUUID)

	// Create io.Reader.
	bigPayload := strings.NewReader(payloadString)
	// Regions: NA and Europe, no Filters for industries, early stage, after 2020
	// bigPayload := strings.NewReader(`{"field_ids":["identifier","operating_status","founded_on","ipo_status","diversity_spotlights","location_identifiers","categories","description","last_funding_type","investor_identifiers","last_funding_at","funding_total","funding_stage","investor_type","last_equity_funding_type","last_funding_total","num_funding_rounds","num_lead_investors","num_investors","semrush_visits_latest_month","semrush_visits_latest_6_months_avg","semrush_visits_mom_pct","semrush_visit_duration","semrush_visit_duration_mom_pct","semrush_visit_pageviews","semrush_visit_pageview_mom_pct","semrush_bounce_rate","semrush_bounce_rate_mom_pct","semrush_global_rank","semrush_global_rank_mom","semrush_global_rank_mom_pct","apptopia_total_apps","apptopia_total_downloads","num_founders","founder_identifiers","num_employees_enum","investor_stage","website","linkedin","num_articles","hub_tags","twitter","facebook","short_description","contact_email","last_key_employee_change_date","last_layoff_date","num_event_appearances","rank_org_company","num_contacts","num_private_contacts","builtwith_num_technologies_used","siftery_num_products","ipqwery_num_patent_granted","ipqwery_num_trademark_registered","private_tags","num_private_notes"],"order":[{"field_id":"founded_on","sort":"desc"}],"query":[{"type":"predicate","field_id":"founded_on","operator_id":"gte","include_nulls":false,"values":["2020"]},{"type":"predicate","field_id":"operating_status","operator_id":"includes","include_nulls":false},{"type":"predicate","field_id":"location_identifiers","operator_id":"includes","values":["b25caef9-a1b8-3a5d-6232-93b2dfb6a1d1","6106f5dc-823e-5da8-40d7-51612c0b2c4e"]},{"type":"predicate","field_id":"funding_stage","operator_id":"includes","values":["seed","early_stage_venture","late_stage_venture"]}],"field_aggregators":[],"collection_id":"organization.companies","limit":1000, "after_id": "${variable}"}`)
	// Remove after_id, use after_id for next request

	// DACH, specialized industries, early stage, after 2020
	// payload := strings.NewReader(`{"field_ids":["identifier","operating_status","founded_on","ipo_status","diversity_spotlights","location_identifiers","categories","description","last_funding_type","investor_identifiers","last_funding_at","funding_total","funding_stage","investor_type","last_equity_funding_type","last_funding_total","num_funding_rounds","num_lead_investors","num_investors","semrush_visits_latest_month","semrush_visits_latest_6_months_avg","semrush_visits_mom_pct","semrush_visit_duration","semrush_visit_duration_mom_pct","semrush_visit_pageviews","semrush_visit_pageview_mom_pct","semrush_bounce_rate","semrush_bounce_rate_mom_pct","semrush_global_rank","semrush_global_rank_mom","semrush_global_rank_mom_pct","apptopia_total_apps","apptopia_total_downloads","num_founders","founder_identifiers","num_employees_enum","investor_stage","website","linkedin","num_articles","hub_tags","twitter","facebook","short_description","contact_email","last_key_employee_change_date","last_layoff_date","num_event_appearances","rank_org_company","num_contacts","num_private_contacts","builtwith_num_technologies_used","siftery_num_products","ipqwery_num_patent_granted","ipqwery_num_trademark_registered","private_tags","num_private_notes"],"order":[{"field_id":"founded_on","sort":"desc"}],"query":[{"type":"predicate","field_id":"founded_on","operator_id":"gte","include_nulls":false,"values":["2020"]},{"type":"predicate","field_id":"operating_status","operator_id":"includes","include_nulls":false},{"type":"predicate","field_id":"location_identifiers","operator_id":"includes","values":["6085b4bf-b18a-1763-a04e-fdde3f6aba94","6d705437-ce74-b061-9864-0079d15fb639","078d9679-a862-02a2-57c8-8337e9a1eec8"]},{"type":"predicate","field_id":"category_groups","operator_id":"includes","values":["85b6bca9-930a-11bc-a608-a513b76fb637","4fe3f3ac-e522-5889-7477-c1b6d6663710","d1079d33-97d7-1f5a-7e6c-b80d5373a3e0","e5514a50-8200-7f6b-de87-b07990670800","26833aa6-0585-2aa7-8c69-63b4b14727c5","ec09d1af-e88f-6a8d-1db8-1dd5e3d49ea0","adc31356-a675-00a1-305e-8becd771319e","133d294c-e5d0-2c4f-9acc-aed0ada1fa8a","701eef4f-18c1-4aff-b550-caf732cd575f","285e29fc-8f70-bf00-1749-9e94158f64f4","2e6eafef-f310-ba60-d932-62f866a87779"]},{"type":"predicate","field_id":"funding_stage","operator_id":"includes","values":["seed","early_stage_venture","late_stage_venture"]}],"field_aggregators":[],"collection_id":"organization.companies","limit":1000}`)

	// Configure a timeout for the client's HTTP request. If the request takes
	// more than this time duration, then it should be cancelled.
	dataRequestDuration := time.Duration(time.Minute * 2)
	ctx, cancel := context.WithTimeout(context.Background(), dataRequestDuration)
	// Cancelling a context releases resources associated with it,
	// cancel should be call as soon as the operations running in a context
	// complete.
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bigPayload)
	if err != nil {
		err = fmt.Errorf("unable to create a new POST request (extract) with timeout context: %w", err)
		return nil, err
	}
	// The general and custom headers are required to trick the
	// Crunchbase API to think that we are not a bot.
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", app.cbCustomHeader.UserAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", app.cbCustomHeader.AcceptLanguage)
	req.Header.Add("Referer", url)

	// Send HTTP request.
	res, err := app.client.Do(req)
	if err != nil {
		err = fmt.Errorf("unable to send HTTP request (extract): %w", err)
		return nil, err
	}
	defer res.Body.Close()

	// Read the response's body.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("unable to read response's body after HTTP req. (extract): %w", err)
		return nil, err
	}
	// Check if we got a 200 Code response, if not return error with
	// status code.
	if res.StatusCode != 200 {
		err = fmt.Errorf("HTTP request (extract) failed with status code %d (%s).\n", res.StatusCode, res.Status)
		return nil, err
	}

	return body, nil
}

// login, requests new login/auth credentials to the Crunchbase API, sets the
// new cookies to the app.Client and if the parameter 'storeCookies' is true, it
// stores the new cookies on a persistent external file.
func (app *application) login(storeCookies bool) error {
	// Create the request to get new session cookies.
	urlSessions := "https://www.crunchbase.com/v4/cb/sessions"
	payloadString := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, app.cbUsername, app.cbPassword)
	payload := strings.NewReader(payloadString)
	// Configure a timeout for the client's HTTP request. If the request takes
	// more than this time duration, then it should be cancelled.
	loginRequestDuration := time.Duration(time.Second * 45)
	ctx, cancel := context.WithTimeout(context.Background(), loginRequestDuration)
	// Canceling a context releases resources associated with it, so code should
	// call cancel as soon as the operations running in a context complete.
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", urlSessions, payload)
	if err != nil {
		return fmt.Errorf("unable to create a new POST request (login) with timeout context: %w", err)
	}

	// The general and custom headers are required to trick the
	// Crunchbase API to think that we are not a bot.
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", app.cbCustomHeader.UserAgent)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Language", app.cbCustomHeader.AcceptLanguage)
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Referer", "https://www.crunchbase.com/login")

	app.infoLog.Print("Requesting new auth credentials from CB API.")
	// Send request through app.client (HTTP Client).
	res, err := app.client.Do(req)
	if err != nil {
		return fmt.Errorf("unable to send HTTP request (login): %w", err)
	}
	defer res.Body.Close()

	// Check if we got a 201 Code response (201 Created), if not return error
	// with status code. The API returns 201 when creating new auth cookies.
	if res.StatusCode != 201 {
		err = fmt.Errorf("HTTP request (login) failed with status code %d (%s).\n", res.StatusCode, res.Status)
		return err
	}

	// Transform Crunchbase's API URL into *url.URL.
	urlObj, err := url.Parse("https://www.crunchbase.com/v4/")
	if err != nil {
		return fmt.Errorf("unable to parse CB url: %w", err)
	}
	// Add retrieved cookies to cookie jar from app.client object.
	app.client.Jar.SetCookies(urlObj, res.Cookies())

	if storeCookies {
		// Store newly created cookies in an external file to guarantee
		// persistency of cookies between different program executions.
		if err = app.storeCookies(res); err != nil {
			return fmt.Errorf("unable to store cookies in external file: %w", err)
		}
	}
	// After a sucessful cookie retrieval, sleep for a randomly-generated amount
	// of seconds, in order to simulate a more human-like online behaviour.
	// A human does not send an API request nanoseconds after logging in.
	delay, err := app.calculateRandomDelay(app.userConfigurations.minDelayLogin, app.userConfigurations.maxDelayLogin)
	if err != nil {
		return fmt.Errorf("unable to create a random time delay after retrieval of the cookies: %w", err)
	}
	app.infoLog.Printf("Delay after a successful cookie retrieval: %ds.\n", delay)
	time.Sleep(time.Duration(delay) * time.Second)

	return nil

}

// getTotalCount, sends a request to the API for a single element, and extracts
// the total count of elements for that particular request. It outputs the total
// count of elements as an int and an error type.
func (app *application) getTotalCount() (int, error) {
	// Limit results in response to only 1 element (parameter true). Start with
	// the first element of the results, lastUUID = "".
	payload, err := app.extract("", true)
	if err != nil {
		err = fmt.Errorf("unable to extract data from the Crunchbase API to get the total count of elements: %w", err)
		return 0, err
	}

	// Create a data container, from which the total count can be extracted.
	dataContainer := DataContainer{}
	err = json.Unmarshal(payload, &dataContainer)
	if err != nil {
		err = fmt.Errorf("unable to parse the JSON-encoded payload into the custom dataContainer data model: %w", err)
		return 0, err
	}
	// Extract the total count value from the dataContainer.
	totalCount := dataContainer.Count

	return totalCount, nil

}

// extractCBData, parses and stores data from the Crunchbase API.
func (app *application) extractCBData() error {

	// Get the total count of elements for a particular request.
	totalCount, err := app.getTotalCount()
	if err != nil {
		return fmt.Errorf("unable to retrieve the total count of entries from an API request: %w", err)
	}

	// Initialize the variables needed for every iteration of the for-loop
	// extracting the Crunchbase data. Implicit initializations are always
	// better to convey the value to other people reading the codebase.
	lastUUID := ""
	// Important change: if one initializes the slice with make and the length
	// of all expected entities/companies, then the condition for the for-loop
	// below always will fail, since the length of the slice is no longer 0.
	// But already the size of all expected elements. Appending new elements
	// should not be a costly operations anyways.
	organizationDocumentSlice := []OrganizationDocument{}

	for len(organizationDocumentSlice) < totalCount {
		// In the first iteration, lastUUID equals "", lastUUID is empty.
		payload, err := app.extract(lastUUID, false)
		if err != nil {
			if len(organizationDocumentSlice) > 0 {
				err = fmt.Errorf("unable to extract further data from the API (the API probably identifies the script as a bot), some results were successfully retrieved and will be exported to a JSON file: %w", err)
				app.errorLog.Print(err)
				// There was a problem while extracting data from the Crunchbase
				// API (the API probably blocks all of our requests because it
				// thinks this script is a bot, but some data (entities) were
				// already stored into the organizationDocumentSlice (since its
				// length is larger than 0), therefore break out of the for-loop
				// to stop sending more requests to the API and store the
				// entities into a temporary file.
				break
			}
			// No entities were stored before getting blocked by the API, so
			// just return from this method without storing any data into an
			// external file.
			return fmt.Errorf("unable to extract any data from the API (bot detection), no results can be exported: %w", err)
		}

		// This is the slice that will be stored into a text file, it grows with
		// each iteration of the for-loop.
		err = decodeBody(payload, &organizationDocumentSlice)
		if err != nil {
			return fmt.Errorf("unable to decode the parsed body which was received from the Crunchbase API: %w", err)
		}

		// Parse lastUUID from last API response.
		lastUUID = organizationDocumentSlice[len(organizationDocumentSlice)-1].Uuid
		// Output the total number of entities extracted sofar, important metric
		// to check consistency in number of extractions.
		app.infoLog.Printf("Total # of entities extracted sofar: %d.", len(organizationDocumentSlice))
		// The program has already fetched and parsed all the available data
		// so it can leave the for-loop without going through a last delay.
		if len(organizationDocumentSlice) >= totalCount {
			break
		}
		// Generate a random delay with a max and min delay constraints.
		delay, err := app.calculateRandomDelay(app.userConfigurations.minDelayExtract, app.userConfigurations.maxDelayExtract)
		if err != nil {
			return fmt.Errorf("unable to create a random delay: %w", err)
		}
		app.infoLog.Printf("Delay until next API request: %ds\n", delay)
		time.Sleep(time.Duration(delay) * time.Second)

	}

	// Encode the organizationDocument slice in JSON format, so that it can be
	// exported into an external file.
	organizationDocumentSliceJSON, err := json.Marshal(organizationDocumentSlice)
	if err != nil {
		return fmt.Errorf("unable to marshal []OrganizationDocument to JSON: %w", err)
	}

	// Create temp file in which to store data.
	f, err := createTempFile()
	if err != nil {
		return fmt.Errorf("unable to create temp file to store data: %w", err)
	}
	defer f.Close()
	// Write data to temp file.
	if _, err := f.Write(organizationDocumentSliceJSON); err != nil {
		return fmt.Errorf("unable to write []OrganizationDocument to temp file: %w", err)
	}

	return nil

}
