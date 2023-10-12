package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBInstance, exported type to interact with a MongoDB instance. Methods
// for CRUD operations act upon this type.
type MongoDBInstance struct {
	// Client, exported MongoDB client used within and outside the package for
	// CRUD operations.
	Client *mongo.Client
}
