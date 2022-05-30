package startup

import (
	"fmt"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/api"
	connection "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/persistence"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	connectionStore := server.initConnectionStore(mongoClient)
	userStore := server.initUserStore(mongoClient)
	connectionService := server.initConnectionService(connectionStore, userStore)
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
	return store
}

func (server *Server) initConnectionService(store domain.ConnectionStore, userStore domain.UserStore) *application.ConnectionService {
	return application.NewConnectionService(store, userStore)
}

func (server *Server) initConnectionHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewCompanyHandler(service)
}

func (server *Server) startGrpcServer(productHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	connection.RegisterConnectionServiceServer(grpcServer, productHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
