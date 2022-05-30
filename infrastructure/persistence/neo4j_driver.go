package persistence

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"os"
)

func GetNeo4jDriver() (neo4j.Driver, error) {
	neo4jUri, found := os.LookupEnv("NEO4J_URI")
	if !found {
		panic("NEO4J_URI not set")
	}
	neo4jUsername, found := os.LookupEnv("NEO4J_USERNAME")
	if !found {
		panic("NEO4J_USERNAME not set")
	}
	neo4jPassword, found := os.LookupEnv("NEO4J_PASSWORD")
	if !found {
		panic("NEO4J_PASSWORD not set")
	}

	driver, err := neo4j.NewDriver(neo4jUri, neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""))
	return driver, err
}
