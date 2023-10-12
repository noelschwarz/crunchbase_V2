package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LinkedInTargetCompany struct {
	OrganizationName string `json:"organizationName" bson:"organizationName"`
	// Timestamp        time.Time `json:"timestamp" bson:"timestamp"`
	UUID string `json:"uuid" bson:"uuid"`
	// Country          string    `json:"country" bson:"country"`
	Linkedin string `json:"linkedin" bson:"linkedin"`
}

// IsCompanyInColl, if a company (represented by its UUID) is present in a
// collection, this functions returns true.
func IsCompanyInColl(client *mongo.Client, database, collection, uuid string) (bool, error) {
	filter := bson.D{{"uuid", uuid}}

	results, err := find(client, database, collection, filter)
	if err != nil {
		return false, fmt.Errorf("could not find company in db: %w", err)
	}
	// Company is not present.
	if len(results) == 0 {
		return false, nil
	}
	return true, nil
}

// FindAfterDate, find all companies in a collection that have a timestamp with
// a date after the date used as parameter.
// This function returns the 'organizationName', 'uuid', 'timestamp' and
// 'linkedin' fields of the documents found.
func FindAfterDate(client *mongo.Client, database, collection string, date time.Time) ([]LinkedInTargetCompany, error) {
	// Find all organizations with a timestamp after 'date'.
	filter := bson.D{{"timestamp", bson.D{{"$gt", date}}}}
	// Project/send back only the fields initialized with a 1.
	projectionOpt := options.Find().SetProjection(bson.D{
		{"organizationName", 1},
		{"uuid", 1},
		// {"country", 1},
		// {"timestamp", 1},
		{"linkedin", 1},
	})
	// Sort by timestamp, starting with newest.
	sortOpt := options.Find().SetSort(bson.D{{"timestamp", -1}})

	results, err := find(client, database, collection, filter, projectionOpt, sortOpt)
	if err != nil {
		return nil, fmt.Errorf("could not find companies after a certain date in db: %w", err)
	}

	return results, nil
}

// func unmarshalLinkedCompanies(companies []bson.D) ([]LinkedInTargetCompany, error) {
// 	results := make([]LinkedInTargetCompany, 0, 20)

// 	for _, company := range companies {
// 		var unmarshallResult LinkedInTargetCompany
// 		if err := bson.Unmarshal(company, &unmarshallResult); err != nil {
// 			return nil, fmt.Errorf("unable to unmarshal a bson slice: %w", err)
// 		}
// 		results = append(results, unmarshallResult)
// 	}

// 	return results, nil
// }

// find, perform a find query on a MongoDB client with a given filter and
// find options. This function returns a slice of bson.D with the results found
// and a non-nil error if a problem happened while searching for results.
func find(client *mongo.Client, database, collection string, filter interface{}, opts ...*options.FindOptions) ([]LinkedInTargetCompany, error) {
	results := make([]LinkedInTargetCompany, 0, 20)

	// TODO: actually add a timeout context.
	ctx := context.TODO()
	coll := client.Database(database).Collection(collection)
	cursor, err := coll.Find(ctx, filter, opts...)
	if err != nil {
		return nil, fmt.Errorf("error could not perform a find query on db: %w", err)
	}
	defer cursor.Close(context.TODO())

	// If an error occurs, the context expires or there are no more documents
	// to be returned from the cursor, it returns false.
	for cursor.Next(context.TODO()) {
		var result LinkedInTargetCompany
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error could not decode result from cursor: %w", err)
		}
		results = append(results, result)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor returned error when looping through db results: %w", err)
	}

	return results, nil
}
