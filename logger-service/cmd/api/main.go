package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/bertoxic/log-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRPCPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app:= Config{
		Models: data.New(client),
	}
	log.Panicln("starting listening at rpc conn")
	go app.rpcListen()
	log.Panicln("starting listening at grpc conn")
	go app.gRPCListen()
	//start server
	//go app.serve()
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Println("error in loer service main, rpc regiser")
	}
	
	log.Println("starting service")
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) rpcListen() error {
	log.Println("about listening to rpc connection")
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {

	clientOps := options.Client().ApplyURI(mongoURL)
	clientOps.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	//connect
	c, err := mongo.Connect(context.TODO(), clientOps)
	if err != nil {
		log.Println("Error connecting to mongodb: ", err)
		return nil, err
	}
	log.Println(" connected to mongodb: 100% ")
	return c, nil
}
