package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// setupApplication, configures the info and error loggers of the application
// type. It configures all needed general parameters for the application,
// e.g. it seeds a *rand.Rand instance to generate random numbers and it
// establishes a connection with a MongoDB instance (remote or local, as
// established in the environment variable in the .env file).
func (app *application) setupApplication() error {
	if err := app.loadEnv(); err != nil {
		return fmt.Errorf("an error occured while loading the environment variables: %w", err)
	}

	// Create a logger for INFO messages, the prefix "INFO" and a tab will be
	// displayed before each log message. The flags Ldate and Ltime provide the
	// local date and time.
	app.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create an ERROR messages logger, addiotionally use the Lshortfile flag to
	// display the file's name and line number for the error.
	app.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create a new and unique random seed, which can be used throughout the
	// application each time that a random string has to be generated. Use as a
	// seed the actual Unix time with nano seconds precision.
	// For more information, check:
	// https://www.calhoun.io/creating-random-strings-in-go/
	app.seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	app.cbCustomHeader = new(CBCustomHeader)

	// Initialize the CB custom HTTP headers for future requests.
	if err := app.initCBCustomHeaders(); err != nil {
		return fmt.Errorf("unable to initialize CB custom HTTP headers: %w", err)
	}

	// Define the delay values used after HTTP requests (in s).
	// Delays after extracting data from API.
	app.userConfigurations.maxDelayExtract = 420
	app.userConfigurations.minDelayExtract = 120
	// Delays after requesting login credentials.
	app.userConfigurations.maxDelayLogin = 120
	app.userConfigurations.minDelayLogin = 60

	return nil
}

// loadEnv, loads the environment variables with sensitive data store in files
// that do not get push to the repository.
func (app *application) loadEnv() error {
	// Load environment variables found in .env file in the same path where this
	// program is running. The .env file contains sensitive data required to for
	// example access the database.
	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("no .env file found: %w", err)
	}

	// Fetch the name of the db used for the data pipeline.
	app.dbName = os.Getenv("DB_NAME")
	if app.dbName == "" {
		return fmt.Errorf("name of database in .env file is empty or not defined")
	}
	// Fetch the name of the collection used for raw CB data.
	app.collCB = os.Getenv("COLL_CB")
	if app.collCB == "" {
		return fmt.Errorf("name of collection for CB raw data in .env file is empty or not defined")
	}

	// Fetch the username and password of CB account.
	app.cbUsername = os.Getenv("CB_USERNAME")
	if app.cbUsername == "" {
		return fmt.Errorf("username of CB account is empty or not defined.")
	}
	app.cbPassword = os.Getenv("CB_PASSWORD")
	if app.cbPassword == "" {
		return fmt.Errorf("password of CB account is empty or not defined.")
	}

	return nil
}

// loadCookies, loads cookies from external file into app.client's cookiejar.
// trunk-ignore(golangci-lint/unused)
func (app *application) loadCookies() error {
	// Open the file with the previously stored cookies.
	cookiesFile, err := os.Open("./cookies.json")
	if err != nil {
		return fmt.Errorf("unable to open cookies.json file: %w", err)
	}
	defer cookiesFile.Close()

	// Store cookies from the file in a slice of *http.Cookie.
	var cookies []*http.Cookie

	// Read all the data out of the file, output is []byte.
	fileByteValue, err := io.ReadAll(cookiesFile)
	if err != nil {
		return fmt.Errorf("unable to read cookies file: %w", err)
	}

	// Decode json from data the cookies file into []*http.Cookie.
	err = json.Unmarshal(fileByteValue, &cookies)
	if err != nil {
		return fmt.Errorf("unable to unmarshal cookies json file: %w", err)
	}

	// Transform Crunchbase's API URL into *url.URL.
	urlObj, err := url.Parse("https://www.crunchbase.com/v4/")
	if err != nil {
		return fmt.Errorf("unable to parse crunchbase.com string into url structure: %w", err)
	}
	// Add retrieved cookies to cookie jar from app.client object.
	app.client.Jar.SetCookies(urlObj, cookies)

	return nil

}

// unmarshalFile, decodes the JSON data from a file (par. fileData) and returns
// the decoded data in an OrganizationDocument slice.
func unmarshalFile(fileData []byte) ([]OrganizationDocument, error) {
	organizationDocumentSlice := []OrganizationDocument{}

	if err := json.Unmarshal(fileData, &organizationDocumentSlice); err != nil {
		err = fmt.Errorf("unable to decode json data into organizationDocumentSlice: %w", err)
		return nil, err
	}

	return organizationDocumentSlice, nil
}

// decodeBody, decodes the body received from the Crunchbase API and returns
// a slice with all parsed entities from the payload.
func decodeBody(payload []byte, organizationDocumentSlice *[]OrganizationDocument) error {
	dataContainer := DataContainer{}
	if err := json.Unmarshal(payload, &dataContainer); err != nil {
		return fmt.Errorf("unable to decode payload []byte into dataContainer type: %w", err)
	}

	organizationDocument := new(OrganizationDocument)

	for _, entity := range dataContainer.Entities {
		// Parse data from one entity into an organizationDocument object.
		if err := organizationDocument.parseRawData(entity); err != nil {
			return fmt.Errorf("unable to parse CB data into organizationDocument type: %w", err)
		}
		// Append the parsed organizationDocument into a slice with all entities
		// contained in the payload from the Crunchbase API.
		*organizationDocumentSlice = append(*organizationDocumentSlice, *organizationDocument)
	}

	return nil

}

// parseRawData, is a method on an OrganizationDocument pointer. It parses the
// CB raw data into the OrganizationDocument, its only parameter is the Entity
// which gets parsed into the *OrganizationDocument.
func (document *OrganizationDocument) parseRawData(entity Entity) error {
	document.Uuid = entity.Uuid
	document.Timestamp = time.Now()
	document.EntityDefId = entity.Properties.Identifier["entity_def_id"]
	document.OrganizationName = entity.Properties.Identifier["value"]
	document.Description = entity.Properties.Description
	document.ShortDescription = entity.Properties.ShortDescription
	document.FundingStage = entity.Properties.FundingStage
	FoundedOnDate, err := time.Parse("2006-01-02", entity.Properties.FoundedOn["value"])
	if err != nil {
		return fmt.Errorf("unable to parse FoundedOn field string into time.Time value: %w", err)
	}
	document.FoundedOn = FoundedOnDate
	document.OperatingStatus = entity.Properties.OperatingStatus
	document.Website = entity.Properties.Website["value"]
	document.Linkedin = entity.Properties.Linkedin["value"]
	document.Facebook = entity.Properties.Facebook["value"]
	document.Industries = entity.Properties.Categories
	document.City = FilterLocation(entity.Properties.Locations, func(location Location) bool {
		return location.LocationType == "city"
	})[0].Name
	document.Country = FilterLocation(entity.Properties.Locations, func(location Location) bool {
		return location.LocationType == "country"
	})[0].Name
	document.ContactEmail = entity.Properties.ContactEmail
	document.NumFounders = entity.Properties.NumFounders
	document.FounderIdentifiers = entity.Properties.FounderIdentifiers
	document.NumOfTechUsed = entity.Properties.NumOfTechUsed
	document.NumPatentGrant = entity.Properties.NumPatentGrant
	document.NumOfArticles = entity.Properties.NumOfArticles
	document.NumTrademarkReg = entity.Properties.NumTrademarkReg
	document.SemRush.NumVisitsLastMonth = entity.Properties.NumVisitsLastMonth
	document.SemRush.BounceRate = entity.Properties.BounceRate
	document.SemRush.NumVisitPerPageviews = entity.Properties.NumVisitPerPageviews
	document.SemRush.VisitDuration = entity.Properties.VisitDuration
	document.NumInvestors = entity.Properties.NumInvestors
	document.FundingTotal = entity.Properties.FundingTotal.ValueUSD
	document.NumFundingRounds = entity.Properties.NumFundingRounds
	document.LastEquityFundingType = entity.Properties.LastEquityFundingType
	document.LastFundingType = entity.Properties.LastFundingType
	document.LastFundingTotal = entity.Properties.LastFundingTotal.ValueUSD
	LastFundingAtDate, err := time.Parse("2006-01-02", entity.Properties.LastFundingAt)
	if err != nil {
		return fmt.Errorf("unable to parse LastFundingAtDate field string into time.Time value: %w", err)
	}
	document.LastFundingAt = LastFundingAtDate
	document.InvestorIdentifiers = entity.Properties.InvestorIdentifiers
	switch entity.Properties.NumEmployeesEnum {
	case "c_00001_00010":
		document.NumEmployeesEnum = "1-10"
	case "c_00011_00050":
		document.NumEmployeesEnum = "11-50"
	case "c_00051_00100":
		document.NumEmployeesEnum = "51-100"
	case "c_00101_00250":
		document.NumEmployeesEnum = "101-250"
	case "c_00251_00500":
		document.NumEmployeesEnum = "251-500"
	case "c_00501_01000":
		document.NumEmployeesEnum = "501-1000"
	case "c_01001_05000":
		document.NumEmployeesEnum = "1001-5000"
	case "c_05001_10000":
		document.NumEmployeesEnum = "5001-10000"
	case "c_10001_max":
		document.NumEmployeesEnum = "10001+"
	}

	return nil
}

func FilterLocation(vs []Location, f func(Location) bool) []Location {
	filtered := make([]Location, 0)
	for _, v := range vs {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// storeCookies, stores cookies in an external file for persistency between
// API calls. The parameter res, is an *http.Response that contains the cookies
// that will be stored in an external file.
func (app *application) storeCookies(res *http.Response) error {
	// Transform cookies into JSON format.
	data, err := json.Marshal(res.Cookies())
	if err != nil {
		return fmt.Errorf("unable to encode []*http Cookies as JSON object: %w", err)
	}
	app.infoLog.Println("Storing cookies in external file...")
	// Write JSON data into external file, no permissions at all for others.
	err = ioutil.WriteFile("./cookies.json", data, 0640)
	if err != nil {
		return fmt.Errorf("unable to write data into cookies.json file: %w", err)
	}

	return nil
}

// initCBCustomHeaders, initializes the HTTP headers used for the CB API, by
// loading the possible HTTP headers values from an external file. Finally,
// the method randomly picks the HTTP headers that will be used by all further
// requests.
func (app *application) initCBCustomHeaders() error {
	// Open the file with the CB custom header configurations.
	CBHeadersFile, err := os.Open("./cb-custom-http-headers.json")
	if err != nil {
		return fmt.Errorf("unable to open cb-custom-http-headers.json file: %w", err)
	}
	defer CBHeadersFile.Close()

	// Read all the data out of the file, output is []byte.
	fileByteValue, err := io.ReadAll(CBHeadersFile)
	if err != nil {
		return fmt.Errorf("unable to read CBHeadersFile: %w", err)
	}
	// Declare cbCustomConfigHeaders container to parse file output to.
	var cbCustomConfigHeaders CBCustomConfigHeaders

	// Parses JSON data in the form of []byte variables into
	// cbCustomConfigHeaders type.
	err = json.Unmarshal(fileByteValue, &cbCustomConfigHeaders)
	if err != nil {
		return fmt.Errorf("unable to decode JSON from CBHeadersFile: %w", err)
	}
	// Based on cbCustomConfigHeaders type, it randomly selects HTTP headers
	// further used by the app for all CB API requests.
	app.randomizeCBCustomHeaders(cbCustomConfigHeaders)

	return nil
}

// randomizeCBCustomHeaders, randomly picks referer, accept-language,
// and user-agent HTTP headers previously parsed from external file.
func (app *application) randomizeCBCustomHeaders(cbCustomConfigHeaders CBCustomConfigHeaders) {
	app.cbCustomHeader.UrlReferer = cbCustomConfigHeaders.UrlReferers[app.seededRand.Intn(len(cbCustomConfigHeaders.UrlReferers))]
	app.cbCustomHeader.AcceptLanguage = cbCustomConfigHeaders.AcceptLanguages[app.seededRand.Intn(len(cbCustomConfigHeaders.AcceptLanguages))]
	app.cbCustomHeader.UserAgent = cbCustomConfigHeaders.UserAgents[app.seededRand.Intn(len(cbCustomConfigHeaders.UserAgents))]
}

// calculateRandomDelay, returns a random delay value in seconds. It uses a
// pre-seeded random number generator (app.seededRand) to generate different
// random numbers in each program execution.
// Input parameters: minDelay and maxDelay: the minimal and maximal possible
// delay value in seconds.
func (app *application) calculateRandomDelay(minDelay, maxDelay int) (int, error) {
	if minDelay >= maxDelay {
		err := fmt.Errorf("calculateRandomDelay: the minimal delay is not smaller than the maximal delay.")
		return 0, err
	}
	if minDelay <= 0 || maxDelay <= 0 {
		err := fmt.Errorf("calculateRandomDelay: the minimal delay or the maximal delay is smaller than 0s.")
		return 0, err
	}
	return app.seededRand.Intn(maxDelay-minDelay) + minDelay, nil
}

// createTempFile, creates a temp file with a unique and random name each time
// the function is called.
func createTempFile() (*os.File, error) {
	// Create temp file in current directory.
	f, err := os.CreateTemp(".", "CBData_*.json")
	if err != nil {
		err = fmt.Errorf("unable to create a temporary file: %w", err)
		return nil, err
	}
	return f, nil
}
