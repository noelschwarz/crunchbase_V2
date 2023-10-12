package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/erodrigufer/UVC_data_pipeline/internal/mongodb"
	"github.com/urfave/cli/v2"
)

// application, sets the types/objects which are needed application-wide, like
// for example a client for a MongoDB instance or a seed for a PRNG.
type application struct {
	// tui, Terminal User Interface for the CLI app.
	tui *cli.App
	// errorLog, error log handler.
	errorLog *log.Logger
	// infoLog, info log handler.
	infoLog *log.Logger
	// client, http.Client that deals with HTTP requests to Crunchbase API.
	client *http.Client
	// mongoDB, is a client for a MongoDB instance running locally.
	mongoDB *mongodb.MongoDBInstance
	// mongoDBURI, the URI of the db to which the script will attempt a
	// connection.
	mongoDBURI string
	// cbCustomHeader, CB based custom HTTP headers that changes for each new
	// run of the script, to make API calls less suspicious.
	cbCustomHeader *CBCustomHeader
	// seededRand, is a *rand.Rand instance seeded from a unique source, used to
	// generate random numbers.
	seededRand *rand.Rand
	// userConfigurations is the struct that stores all the user-defined
	// configuration values.
	userConfigurations userConfigurations
	// dbName, name of database used in the local MongoDB instance.
	dbName string
	// collCB, collection used for Crunchbase raw data.
	collCB string
	// cbUsername, username of the Crunchbase account used for requests.
	cbUsername string
	// cbPassword, password of the CB account which is used for requests.
	cbPassword string
}

// userConfigurations, stores the user-defined configurations which can be
// defined at the initialization of the program execution through flags. If the
// values are not defined by the user with flags, some defaults values will be
// used by the program.
type userConfigurations struct {
	// maxDelayLogin, defines the maximal amount of delay (in s) after the login
	// request is sent to the API.
	maxDelayLogin int
	// minDelayLogin, defines the minimal amount of delay (in s) after the login
	// request is sent to the API.
	minDelayLogin int
	// maxDelayExtract, defines the maximal amount of delay (in s) to wait,
	// after extracting data from the API.
	maxDelayExtract int
	// minDelayExtract, defines the minimal amount of delay (in s) to wait,
	// after extracting data from the API.
	minDelayExtract int
}

// DataContainer, type of data container that unpacks Crunchbase JSON output
// with count and entities.
type DataContainer struct {
	Count    int      `json:"count"`
	Entities []Entity `json:"entities"`
}

// Entity, type of entity that represents all organizational data such as uuid
// (identifier) and organization characteristics (properties).
type Entity struct {
	Uuid       string     `json:"uuid"`
	Properties Properties `json:"properties"`
}

// Person, type of person to store personal references (investors, founders,
// employees) of datapoints.
type Person struct {
	Uuid        string `json:"uuid" bson:"uuid"`
	EntityDefId string `json:"entity_def_id" bson:"entity_def_id"`
	Permalink   string `json:"permalink" bson:"permalink"`
	Name        string `json:"value" bson:"value"`
}

// Category, type of category to store industry references of datapoints.
type Category struct {
	Uuid        string `json:"uuid" bson:"uuid"`
	EntityDefId string `json:"entity_def_id" bson:"entity_def_id"`
	Name        string `json:"value" bson:"value"`
}

// Location, type of location to store geographical references of datapoints.
type Location struct {
	Uuid         string `json:"uuid"`
	LocationType string `json:"location_type"`
	Name         string `json:"value"`
}

// Funding, type parses only the USD value as an int.
type Funding struct {
	ValueUSD int `json:"value_usd"`
}

// Properties, type of properties that unpacks all the organization
// characteristics such as IPOStatus, FoundedOn, etc.
type Properties struct {
	FoundedOn             map[string]string `json:"founded_on"`
	Website               map[string]string `json:"website"`
	Identifier            map[string]string `json:"identifier"`
	NumFounders           int               `json:"num_founders"`
	FounderIdentifiers    []Person          `json:"founder_identifiers"`
	Description           string            `json:"description"`
	Linkedin              map[string]string `json:"linkedin"`
	Facebook              map[string]string `json:"facebook"`
	ShortDescription      string            `json:"short_description"`
	NumInvestors          int               `json:"num_investors"`
	OperatingStatus       string            `json:"operating_status"`
	NumEmployeesEnum      string            `json:"num_employees_enum"`
	FundingTotal          Funding           `json:"funding_total"`
	FundingStage          string            `json:"funding_stage"`
	NumFundingRounds      int               `json:"num_funding_rounds"`
	LastEquityFundingType string            `json:"last_equity_funding_type"`
	InvestorIdentifiers   []Person          `json:"investor_identifiers"`
	LastFundingTotal      Funding           `json:"last_funding_total"`
	LastFundingType       string            `json:"last_funding_type"`
	LastFundingAt         string            `json:"last_funding_at"`
	Categories            []Category        `json:"categories"`
	Locations             []Location        `json:"location_identifiers"`
	NumTrademarkReg       int               `json:"ipqwery_num_trademark_registered"`
	NumPatentGrant        int               `json:"ipqwery_num_patent_granted"`
	NumOfTechUsed         int               `json:"builtwith_num_technologies_used"`
	NumOfArticles         int               `json:"num_articles"`
	ContactEmail          string            `json:"contact_email"`
	NumVisitsLastMonth    int               `json:"semrush_visits_latest_month"`
	NumVisitPerPageviews  float32           `json:"semrush_visit_pageviews"`
	BounceRate            float32           `json:"semrush_bounce_rate"`
	VisitDuration         int               `json:"semrush_visit_duration"`
}

// SemRush, type of SemRush that holds all the SemRush relevant data such as
// DurationVisit or Bouncerate.
type SemRush struct {
	NumVisitsLastMonth   int     `json:"sr_visits_latest_month" bson:"sr_visits_latest_month"`
	VisitDuration        int     `json:"sr_visit_duration" bson:"sr_visit_duration"`
	BounceRate           float32 `json:"sr_bounce_rate" bson:"sr_bounce_rate"`
	NumVisitPerPageviews float32 `json:"sr_visit_pageviews" bson:"sr_visit_pageviews"`
}

// OrgnizationDocument, type of OrganizationDocument that holds all the previous
// data from Crunchbase and that will persisted in our database later.
type OrganizationDocument struct {
	Uuid                  string     `json:"uuid" bson:"uuid"`
	Timestamp             time.Time  `json:"timestamp" bson:"timestamp"`
	EntityDefId           string     `json:"entityDefId" bson:"entityDefId"`
	OrganizationName      string     `json:"organizationName" bson:"organizationName"`
	Description           string     `json:"description" bson:"description"`
	ShortDescription      string     `json:"shortDescription" bson:"shortDescription"`
	FundingStage          string     `json:"fundingStage" bson:"fundingStage"`
	FoundedOn             time.Time  `json:"foundedOn" bson:"foundedOn"`
	OperatingStatus       string     `json:"operatingStatus" bson:"operatingStatus"`
	Website               string     `json:"website" bson:"website"`
	Linkedin              string     `json:"linkedin" bson:"linkedin"`
	Facebook              string     `json:"facebook" bson:"facebook"`
	Industries            []Category `json:"industries" bson:"industries"`
	City                  string     `json:"city" bson:"city"`
	Country               string     `json:"country" bson:"country"`
	ContactEmail          string     `json:"contactEmail" bson:"contactEmail"`
	NumFounders           int        `json:"numFounders" bson:"numFounders"`
	NumEmployeesEnum      string     `json:"num_employees_enum" bson:"num_employees_enum"`
	FounderIdentifiers    []Person   `json:"founderIdentifiers" bson:"founderIdentifiers"`
	NumOfTechUsed         int        `json:"numOfTechUsed" bson:"numOfTechUsed"`
	NumOfArticles         int        `json:"numArticles" bson:"numArticles"`
	NumTrademarkReg       int        `json:"numTrademarkReg" bson:"numTrademarkReg"`
	NumPatentGrant        int        `json:"numPatentGrant" bson:"numPatentGrant"`
	SemRush               SemRush    `json:"semRush" bson:"semRush"`
	NumInvestors          int        `json:"numInvestors" bson:"numInvestors"`
	FundingTotal          int        `json:"fundingTotal" bson:"fundingTotal"`
	NumFundingRounds      int        `json:"numFundingRounds" bson:"numFundingRounds"`
	LastEquityFundingType string     `json:"lastEquityFundingType" bson:"lastEquityFundingType"`
	LastFundingType       string     `json:"lastFundingType" bson:"lastFundingType"`
	LastFundingTotal      int        `json:"lastFundingTotal" bson:"lastFundingTotal"`
	LastFundingAt         time.Time  `json:"lastFundingAt" bson:"lastFundingAt"`
	InvestorIdentifiers   []Person   `json:"investorIdentifiers" bson:"investorIdentifiers"`
}

// CBCustomConfigHeaders, struct with CB custom HTTP headers that holds
// possible values for referers, user-agent and accept-language headers, after
// being parsed from external file.
type CBCustomConfigHeaders struct {
	UrlReferers     []string `json:"cb-referer-url"`
	AcceptLanguages []string `json:"cb-accepted-lg"`
	UserAgents      []string `json:"cb-user-agent"`
}

// CBCustomHeader, struct of CB custom HTTP headers that holds randomly chosen
// configurations for referers, user-agent and accept-language HTTP headers
// based on the CBCustomConfigHeaders.
type CBCustomHeader struct {
	UrlReferer     string
	AcceptLanguage string
	UserAgent      string
}
