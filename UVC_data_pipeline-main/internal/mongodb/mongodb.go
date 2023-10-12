package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo

// NewClient, function returns a *MongoDBInstance, which can then be used as an
// abstracted client for a local MongoDB instance. It connects to a local
// MongoDB instance with the URI found in the local .env file. If the connection
// attempt fails or pinging the database fails, this method returns an error.
func NewClient(mongodbURI string) (*MongoDBInstance, error) {
	dbClient := new(MongoDBInstance)
	if err := dbClient.connectDB(mongodbURI); err != nil {
		err = fmt.Errorf("unable to connect to MongoDB instance: %w", err)
		return nil, err
	}

	return dbClient, nil
}

// connectDB, this method uses an URI for a local MongoDB instance to connect
// to the DB and store the DB client as a field in a MongoDBInstance object.
func (db *MongoDBInstance) connectDB(mongodbURI string) error {
	if mongodbURI == "" {
		return fmt.Errorf("MONGODB_URI_REMOTE variable was not set in the .env file")
	}
	// Configure a timeout for establishing a DB connection.
	timeoutDB, err := time.ParseDuration("20s")
	if err != nil {
		return fmt.Errorf("error while parsing time duration: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	// Port 27017 is the default port for a local MongoDB daemon.
	db.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI))
	if err != nil {
		return fmt.Errorf("unable to establish a connection with a MongoDB instance: %w", err)
	}

	// mongo.Connect() does not return an error if the deployment is down. To
	// check for that it is better to call Ping().
	if err := db.Client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("unable to ping MongoDB instance: %w", err)
	}

	return nil
}

// connectDBSSL, establish a connection to a MongoDB instance (also remotely)
// through SSL.
// IMPORTANT: due to configuration issues, this method still has not properly
// worked, while connecting remotely to the db.
func (db *MongoDBInstance) connectDBSSL() error {
	// Get the value stored in the environment variable, which should have been
	// previously loaded from the .env file.
	uri := os.Getenv("MONGODB_URI_SSL")
	if uri == "" {
		return fmt.Errorf("MONGODB_URI_SSL variable was not set in the .env file")
	}
	// SSL Authentication.
	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
	}
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)

	// Configure a timeout for establishing a DB connection.
	timeoutDB := time.Duration(20) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	var err error
	db.Client, err = mongo.Connect(ctx, clientOpts)
	if err != nil {
		return fmt.Errorf("unable to establish a connection with a MongoDB instance: %w", err)
	}

	// mongo.Connect() does not return an error if the deployment is down. To
	// check for that it is better to call Ping().
	if err := db.Client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("unable to ping MongoDB instance: %w", err)
	}

	return nil

}

// InsertMultipleDocuments, inserts the documents (parameter: docs) in the
// database (par: dbName) inside the collection (par: coll).
func (db *MongoDBInstance) InsertMultipleDocuments(docs []interface{}, dbName, coll string) error {
	collection := db.Client.Database(dbName).Collection(coll)

	// Configure a timeout for inserting documents.
	timeoutDB, err := time.ParseDuration("120s")
	if err != nil {
		return fmt.Errorf("could not parse time duration for ctx timeout: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDB)
	defer cancel()

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("could not insert (many) documents to collection in db: %w", err)
	}

	return nil
}
