package startup

import (
	"fmt"
	saga "github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_APIGateway/saga/messaging/nats"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/application"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/domain"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/api"
	connection "github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/grpc/proto"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/infrastructure/persistence"
	"github.com/XWS-BSEP-Tim-13/Dislinkt_ConnectionService/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
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
	QueueGroup     = "block_user"
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
	eventStore := server.initEventStore(mongoClient)
	neo4jDriver := server.initNeo4jDriver()
	neo4jConnectionStore := server.initNeo4jConnectionStore(neo4jDriver)
	seedConnectionStore(neo4jConnectionStore, userStore)

	commandPublisher := server.initPublisher(server.config.BlockUserCommandSubject)
	replySubscriber := server.initSubscriber(server.config.BlockUserReplySubject, QueueGroup)
	createBlockOrchestrator := server.initCreateOrderOrchestrator(commandPublisher, replySubscriber)
	commandSubscriber := server.initSubscriber(server.config.BlockUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.BlockUserReplySubject)

	eventService := server.initEventService(eventStore)
	server.initEventHandler(eventService)

	connectionService := server.initConnectionService(connectionStore, userStore, neo4jConnectionStore, createBlockOrchestrator, eventStore)
	server.initBlockUserHandler(connectionService, replyPublisher, commandSubscriber)

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

	for _, connection := range connectionRequests {
		err := store.Insert(connection)
		if err != nil {
			log.Fatal(err)
		}
	}

	return store
}

func (server *Server) initEventStore(client *mongo.Client) domain.EventStore {
	store := persistence.NewEventMongoDBStore(client)
	store.DeleteAll()
	for _, event := range events {
		err := store.Insert(event)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initCreateOrderOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *application.BlockUserOrchestrator {
	orchestrator, err := application.NewBlockUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initBlockUserHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewBlockUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
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

	connStore.CreateConnectionBetweenUsers(userSrki, userAna)
	connStore.CreateConnectionBetweenUsers(userAna, userSrki)
	connStore.CreateConnectionBetweenUsers(userSrki, userLjuba)
	connStore.CreateConnectionBetweenUsers(userLenka, userMarija)
	connStore.CreateConnectionBetweenUsers(userMarija, userAna)
	connStore.CreateConnectionBetweenUsers(userAna, userLenka)

	connStore.AddSkillToUser(userAna, "Java")
	connStore.AddSkillToUser(userAna, "Docker")
	connStore.AddSkillToUser(userMarija, "AWS")
	connStore.AddSkillToUser(userMarija, "Docker")

	connStore.AddExperienceToUser(userMarija, &userMarija.Experiences[0])
	connStore.AddExperienceToUser(userMarija, &userMarija.Experiences[1])

	connStore.AddJobOfferFromCompany(companies[0], jobs[1])
	connStore.AddJobOfferFromCompany(companies[1], jobs[0])

	connStore.AddRequiredSkillToJobOffer("AWS", jobs[0])
	connStore.AddRequiredSkillToJobOffer("Docker", jobs[0])
	connStore.AddRequiredSkillToJobOffer("Java", jobs[1])
}

func (server *Server) initConnectionService(store domain.ConnectionStore, userStore domain.UserStore, neo4jStore persistence.ConnectionNeo4jStore, orchestrator *application.BlockUserOrchestrator, eventStore domain.EventStore) *application.ConnectionService {
	return application.NewConnectionService(store, userStore, neo4jStore, orchestrator, eventStore)
}

func (server *Server) initConnectionHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewCompanyHandler(service)
}

func (server *Server) initEventService(store domain.EventStore) *application.EventService {
	return application.NewEventService(store)
}

func (server *Server) initEventHandler(service *application.EventService) *api.EventsHandler {
	return api.NewEventsHandler(service)
}

func (server *Server) startGrpcServer(productHandler *api.ConnectionHandler) {
	/*cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
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
	}*/

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
