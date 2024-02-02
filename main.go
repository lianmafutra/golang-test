package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"


	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}
func loaderioHandler(w http.ResponseWriter, r *http.Request) {
	// Replace "YOUR_LOADERIO_TOKEN" with the actual token provided by Loader.io
	verificationToken := "loaderio-e1a6aec71495d1efb7865f9cf35b0f71"

	// Respond with the verification token
	fmt.Fprintf(w, verificationToken)
}

func main() {
	// Set MongoDB connection parameters with username, password, and database
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/MongoTestingDb").
		SetAuth(options.Credential{
			Username: "admin",
			Password: "admin",
		})

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Close the connection when the main function exits
	defer func() {
		if err = client.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// Get a handle to the "example" database and the "people" collection
	database := client.Database("MongoTestingDb")
	collection := database.Collection("people")

// Create a new HTTP server and register the handler
	http.HandleFunc("/loaderio-e1a6aec71495d1efb7865f9cf35b0f71.txt", loaderioHandler)

	// Define an HTTP handler to insert a person into MongoDB
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse the request body into a Person struct
		var person Person
		err := json.NewDecoder(r.Body).Decode(&person)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Insert the person into the "people" collection
		insertResult, err := collection.InsertOne(context.Background(), person)
		if err != nil {
			http.Error(w, "Failed to insert data into MongoDB", http.StatusInternalServerError)
			return
		}

		// Respond with the ID of the inserted document
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Inserted document with ID: %v\n", insertResult.InsertedID)
	})

	// Start the HTTP server on port 8080
	log.Println("Server listening on :8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
