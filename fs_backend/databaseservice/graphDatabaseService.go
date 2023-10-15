package databaseservice

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphDatabaseService struct {
	ctx    context.Context
	driver neo4j.DriverWithContext
}

func (gds *GraphDatabaseService) Connect() {
	// Creating a context
	gds.ctx = context.Background()

	// Reading required env variables
	uri := os.Getenv("NEO4J_URI")
	if uri == "" {
		log.Default().Println("Needed a Neo4j URI")
	}
	user := os.Getenv("NEO4J_USER")
	if user == "" {
		log.Default().Println("Needed a Neo4j user")
	}
	password := os.Getenv("NEO4J_PASSWORD")
	if password == "" {
		log.Default().Println("Needed a Neo4j password")
	}

	log.Default().Println("Connecting to Neo4j server at", uri)

	// Connecting with server and verifying connectivity
	var err error
	gds.driver, err = neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(user, password, ""))
	if err != nil {
		log.Default().Println("Error connecting Neo4j server : ", err.Error())
	}
	err = gds.driver.VerifyConnectivity(gds.ctx)
	if err != nil {
		log.Fatalln("Error verifying Neo4j server : ", err.Error())
		return
	}
	log.Default().Println("Connected to Neo4j server at", uri)
}

func (gds *GraphDatabaseService) Close() {
	gds.driver.Close(gds.ctx)
	gds.ctx.Done()
}
