package main

import (
	"context"
	"log"
	"time"

	"github.com/bertoxic/log-service/data"
	"go.mongodb.org/mongo-driver/bson"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
    log.Printf("currently in loginfo in loggrer rpc data obtained is %s", payload)
    _ = data.LogEntry{}

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), bson.M{
        "name" : payload.Name,
        "data" : payload.Data,
        "created_at" : time.Now(),
    })
    if err != nil {
        log.Printf("error in adding data to mongodb using rpc : %s", err.Error())
        return err
    }
    *resp = "processed payload via RPC" + payload.Name
    return nil
}