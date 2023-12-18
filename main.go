package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gofr.dev/pkg/gofr"
)

// STRUCTURE FOR THE LAPTOP REPAIR ENTRY
type Repair struct {
	CustomerName string `json:"customerName"`
	RepairID     string `json:"repairID"`
	DeviceName   string `json:"deviceName"`
	Status       string `json:"status"`
}

// CREATING DATABASE
var (
	mongoURI    = "mongodb://localhost:27017"
	database    = "repair_shop"
	collection  = "repairs"
	mongoClient *mongo.Client
)

func initMongoDB() error {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}
	fmt.Println("Connected to MongoDB!")
	mongoClient = client
	return nil
}

func main() {
	err := initMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	// CREATING GoFr SERVER
	app := gofr.New()

	// CREATE OPERATION
	app.POST("/repair/submit", func(ctx *gofr.Context) (interface{}, error) {
		var repair Repair
		err := json.NewDecoder(ctx.Request().Body).Decode(&repair)
		if err != nil {
			return nil, err
		}
		repairsCollection := mongoClient.Database(database).Collection(collection)
		_, err = repairsCollection.InsertOne(context.Background(), repair)
		if err != nil {
			return nil, err
		}
		return "Laptop repair submitted successfully!", nil
	})

	// READ OPERATION
	app.GET("/repairs/list", func(ctx *gofr.Context) (interface{}, error) {
		var repairs []Repair
		repairsCollection := mongoClient.Database(database).Collection(collection)
		cursor, err := repairsCollection.Find(context.Background(), bson.M{})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var repair Repair
			if err := cursor.Decode(&repair); err != nil {
				return nil, err
			}
			repairs = append(repairs, repair)
		}

		if err := cursor.Err(); err != nil {
			return nil, err
		}
		return repairs, nil
	})

	// UPDATE OPERATION
	app.PUT("/repairs/list/{repairid}", func(ctx *gofr.Context) (interface{}, error) {
		rawPath := ctx.Request().URL.Path
		pathParts := strings.Split(rawPath, "/")
		repairid := pathParts[3]

		fmt.Println("Handling PUT request for ISBN:", repairid)
		var updatedRepair Repair
		err := json.NewDecoder(ctx.Request().Body).Decode(&updatedRepair)
		if err != nil {
			return nil, err
		}
		RepairsCollection := mongoClient.Database(database).Collection(collection)
		filter := bson.M{"repairid": repairid}
		update := bson.M{"$set": updatedRepair}

		result, err := RepairsCollection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			return nil, err
		}
		if result.ModifiedCount == 0 {
			return nil, fmt.Errorf("No Item found with Repairid: %s", repairid)
		}
		return "Laptop updated successfully", nil
	})

	// DELETE OPERATION
	app.DELETE("/repair/delete/{repairiD}", func(ctx *gofr.Context) (interface{}, error) {
		rawPath := ctx.Request().URL.Path
		pathParts := strings.Split(rawPath, "/")
		repairid := pathParts[3]

		repairsCollection := mongoClient.Database(database).Collection(collection)
		result, err := repairsCollection.DeleteOne(context.Background(), bson.M{"repairid": repairid})
		if err != nil {
			return nil, err
		}
		if result.DeletedCount == 0 {
			return nil, fmt.Errorf("No item found with REPAIR_ID: %s", repairid)
		}
		return "Laptop with Repair_Id is deleted successfully", nil
	})

	// STARTING THE SERVER
	app.Start()
}
