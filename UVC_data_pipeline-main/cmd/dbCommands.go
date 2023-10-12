package main

import (
	"fmt"
	"os"

	"github.com/erodrigufer/UVC_data_pipeline/internal/mongodb"
)

// formatMongoURI, format the Mongo URI properly. Add an IP address to the
// Mongo URI.
func (app *application) formatMongoURI(IP string) error {
	// Fetch the URI of the MongoDB instance to which the script will attempt
	// a connection. The URI is missing the IP address portion.
	URIIncomplete := os.Getenv("MONGODB_URI_REMOTE")
	if URIIncomplete == "" {
		return fmt.Errorf("MONGODB_URI_REMOTE in .env file is empty or not defined.")
	}

	if IP == "" {
		return fmt.Errorf("IP address passed to function is an empty string.")
	}

	// Add IP address of MongoDB instance to URI.
	// URIIncomplete has the following format:
	// mongodb://${DB_USERNAME}:${DB_PASSWORD}@%s:27017/${DB_NAME}
	// '%s' is replaced by the appropriate IP with fmt.Sprintf.
	app.mongoDBURI = fmt.Sprintf(URIIncomplete, IP)

	return nil
}

// setupDBCommands, setup and configure DB to start a connection.
func (app *application) setupDBCommands(IP string) error {
	if err := app.formatMongoURI(IP); err != nil {
		return fmt.Errorf("error while formatting Mongo URI: %w", err)
	}
	if err := app.configureMongoDB(); err != nil {
		return fmt.Errorf("error while configuring database: %w", err)
	}
	return nil
}

// configureMongoDB, create a new client with the DB using the URI.
func (app *application) configureMongoDB() error {
	var err error
	// Connect to local MongoDB instance and get a client.
	app.mongoDB, err = mongodb.NewClient(app.mongoDBURI)
	if err != nil {
		return fmt.Errorf("connection attempt or pinging the database failed: %w", err)
	}
	app.infoLog.Printf("Successfully established connection to MongoDB instance.")

	return nil
}

// insertDB, inserts a local `file` (path of file) into a MongoDB instance.
func (app *application) insertDB(file string) error {
	// Read data from file (path of file)
	fileData, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("unable to read data from file %s: %w", file, err)
	}
	// Unmarshal .json data into slice.
	// We suppose that the file's data is composed of many mongo documents.
	organizationDocumentSlice, err := unmarshalFile(fileData)
	if err != nil {
		return fmt.Errorf("unable to unmarshal data from file %s: %w", file, err)
	}

	// Documents to be inserted into the db. Create an interface{} slice of the
	// correct size.
	docs := make([]interface{}, len(organizationDocumentSlice))
	// Populate the interface{} with the values.
	for i, u := range organizationDocumentSlice {
		docs[i] = u
	}
	// Insert the documents into the DB.
	if err := app.mongoDB.InsertMultipleDocuments(docs, app.dbName, app.collCB); err != nil {
		return fmt.Errorf("failed to insert multiple documents into DB: %w", err)
	}

	return nil
}
