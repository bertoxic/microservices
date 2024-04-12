package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/bertoxic/log-service/data"
	"github.com/bertoxic/log-service/logs"
	"google.golang.org/grpc"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
    Models data.Models
}


func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest)(*logs.LogResponse, error){
    log.Println("inside write log in logger-service")
    input := req.GetLogEntry()

    //writelog
    logEntry := data.LogEntry {
        Name: input.Name,
        Data: input.Data,
    }
    err := l.Models.LogEntry.Insert(logEntry)
    if err != nil {
        res := &logs.LogResponse{
            Result: "failed",
        }
        return res, err
    }
    log.Println("done with write log  in logger-service")
    //return a response
    res := &logs.LogResponse{Result: "successfully logged!"}
    return res, nil
}

func ( app *Config) gRPCListen(){
         lis, err := net.Listen("tcp", fmt.Sprintf(":%s",gRPCPort))
         if err != nil {
            log.Fatalf("failed to listen for grpc: %v", err)
         }
         s := grpc.NewServer()

         logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})
         log.Println("grpc server started on port 80 ")
        
         s.Serve(lis)

}