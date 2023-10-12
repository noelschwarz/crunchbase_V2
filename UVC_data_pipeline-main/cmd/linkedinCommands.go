package main

import (
	"fmt"
	"time"

	"github.com/erodrigufer/UVC_data_pipeline/internal/mongodb"
)

// updateUniqueCompanies, finds all companies after a given date and checks if
// all the companies have already been added to the collection with all the
// companies' LinkedIn URLs (the LinkedIn targets).
func (app *application) updateUniqueCompanies(date string) error {
	// Convert date parameter to time.Time type.
	dateParsed, err := parseDate(date)
	if err != nil {
		return err
	}
	// Find all companies with a timestamp after 'date'.
	// TODO: check if this function really returns an error, if no companies
	// were found after a particular date, e.g. because there were no companies
	// after a particular date.
	companies, err := mongodb.FindAfterDate(app.mongoDB.Client, "datapipeline", "crunchbaseTest", dateParsed)
	if err != nil {
		return fmt.Errorf("could not find companies after date: %w", err)
	}

	// Check if every company is already in the collection with the LinkedIn
	// targets.
	// The slice newCompanies collects all the new companies that will be added
	// in this round.
	newCompanies := make([]mongodb.LinkedInTargetCompany, 0, 20)
	for _, company := range companies {
		ok, err := mongodb.IsCompanyInColl(app.mongoDB.Client, "datapipeline", "linkedinCompanyTargets", company.UUID)
		if err != nil {
			return fmt.Errorf("could not check if company is already present in collection: %w", err)
		}
		if ok {
			// Company is already in the collection. Check next company.
			app.infoLog.Printf("%s is already in collection.", company.OrganizationName)
			continue
		} else {
			// Company is not in the collection yet. Append it to slice of new
			// companies, if the company has a valid LinkedIn URL (some
			// companies do not provide a LinkedIn URL).
			if company.Linkedin != "" {
				newCompanies = append(newCompanies, company)
				app.infoLog.Printf("%s will be added to the collection.", company.OrganizationName)
			}
		}

	}

	for i, r := range newCompanies {
		fmt.Printf("%d. %s -- %s\n", i, r.OrganizationName, r.Linkedin)
	}

	// If newCompanies slice is not empty, add new companies to collection.
	if len(newCompanies) != 0 {
		// Documents to be inserted into the db. Create an interface{} slice of the
		// correct size.
		docs := make([]interface{}, len(newCompanies))
		// Populate the interface{} with the values.
		for i, u := range newCompanies {
			docs[i] = u
		}
		app.infoLog.Print("Inserting new companies into 'linkedinCompanyTargets' collection.")
		// Insert the documents into the DB.
		if err := app.mongoDB.InsertMultipleDocuments(docs, "datapipeline", "linkedinCompanyTargets"); err != nil {
			return fmt.Errorf("failed to insert multiple documents into DB: %w", err)
		}
	}

	return nil
}

// companyIsInColl, checks if a company (UUID) is present in a collection.
func (app *application) companyIsInColl(uuid string) error {
	ok, err := mongodb.IsCompanyInColl(app.mongoDB.Client, "datapipeline", "linkedinCompanyTargets", uuid)
	if ok {
		fmt.Printf("Company (UUID: %s) is present in the collection.\n", uuid)
	} else {
		fmt.Printf("Company (UUID: %s) is NOT present in the collection.\n", uuid)
	}
	return err
}

func (app *application) findCompaniesTimestamp(date string) error {
	// Convert date parameter to time.Time type.
	dateParsed, err := parseDate(date)
	if err != nil {
		return err
	}
	// Find all companies with a timestamp after 'date'.
	results, err := mongodb.FindAfterDate(app.mongoDB.Client, "datapipeline", "crunchbaseRaw", dateParsed)

	for i, r := range results {
		if r.Linkedin != "" {
			fmt.Printf("%d. %s -- %s\n", i, r.OrganizationName, r.Linkedin)
		}
	}
	// fmt.Println(results)

	return err
}

// parseDate, parse a date string into a time.Time type. If the given date
// paremeter is not parseble the function returns an error.
func parseDate(date string) (time.Time, error) {
	// Layout used by function to parse dates.
	// Year-Month-Day (Month= 3 letter representation in English).
	const dateLayout = "2006-Jan-02"
	t, err := time.Parse(dateLayout, date)
	if err != nil {
		return t, fmt.Errorf("could not parse date: %w", err)
	}

	return t, nil
}
