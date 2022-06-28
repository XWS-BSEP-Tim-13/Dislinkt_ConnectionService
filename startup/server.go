package startup

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/api"
	connection "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/persistence"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

const (
	serverCertFile = "cert/cert.pem"
	serverKeyFile  = "cert/key.pem"
	clientCertFile = "cert/client-cert.pem"
)

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	connectionStore := server.initConnectionStore(mongoClient)
	userStore := server.initUserStore(mongoClient)
	neo4jDriver := server.initNeo4jDriver()
	neo4jConnectionStore := server.initNeo4jConnectionStore(neo4jDriver)
	seedConnectionStore(neo4jConnectionStore, userStore)
	connectionService := server.initConnectionService(connectionStore, userStore, neo4jConnectionStore)
	connectionHandler := server.initConnectionHandler(connectionService)
	server.startGrpcServer(connectionHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.ConnectionDBHost, server.config.ConnectionDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initConnectionStore(client *mongo.Client) domain.ConnectionStore {
	store := persistence.NewConnectionMongoDBStore(client)
	store.DeleteAll()

	for _, connection := range connections {
		err := store.Insert(connection)
		if err != nil {
			log.Fatal(err)
		}
	}

	return store
}

func (server *Server) initUserStore(client *mongo.Client) domain.UserStore {
	store := persistence.NewUserMongoDBStore(client)
	store.DeleteAll()
	for _, user := range users {
		err := store.Insert(user)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initNeo4jDriver() neo4j.Driver {
	driver, _ := persistence.GetNeo4jDriver()
	return driver
}

func (server *Server) initNeo4jConnectionStore(driver neo4j.Driver) persistence.ConnectionNeo4jStore {
	neo4jConnectionStore := persistence.NewConnectionNeo4jStore(driver)
	return neo4jConnectionStore
}

func seedConnectionStore(connStore persistence.ConnectionNeo4jStore, userStore domain.UserStore) {
	userAna, _ := userStore.GetActiveByUsername("anagavrilovic")
	userSrki, _ := userStore.GetActiveByUsername("srdjansukovic")
	userLjuba, _ := userStore.GetActiveByUsername("stefanljubovic")
	userMarija, _ := userStore.GetActiveByUsername("marijakljestan")
	userLenka, _ := userStore.GetActiveByUsername("lenka")

	connStore.CreateConnection(userSrki, userAna)
	connStore.CreateConnection(userAna, userSrki)
	connStore.CreateConnection(userAna, userLjuba)
	connStore.CreateConnection(userSrki, userLjuba)
	connStore.CreateConnection(userLjuba, userMarija)
	connStore.CreateConnection(userMarija, userAna)
	connStore.CreateConnection(userLenka, userAna)

	connStore.AddSkillToUser(userMarija, "AWS")
	fmt.Println("123")
	connStore.AddSkillToUser(userMarija, "Docker")
}

func (server *Server) initConnectionService(store domain.ConnectionStore, userStore domain.UserStore, neo4jStore persistence.ConnectionNeo4jStore) *application.ConnectionService {
	return application.NewConnectionService(store, userStore, neo4jStore)
}

func (server *Server) initConnectionHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewCompanyHandler(service)
}

func (server *Server) startGrpcServer(productHandler *api.ConnectionHandler) {
	cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	pemClientCA, err := ioutil.ReadFile(clientCertFile)
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		log.Fatal(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequestClientCert,
		ClientCAs:    certPool,
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(config)),
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(opts...)
	connection.RegisterConnectionServiceServer(grpcServer, productHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
