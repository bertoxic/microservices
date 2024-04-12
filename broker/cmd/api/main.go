package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"
// const webPort = "8080"

type Config struct{
    Rabbit *amqp.Connection
}

func main() {

        rabbitConn, err := connect()
        if err != nil {
            log.Println("cannot connect to rabbitconn")
        }
		app := Config{
            Rabbit: rabbitConn,
        }

		log.Printf("starting broker service on port %s", webPort)

        srv := http.Server{}

        srv.Addr = fmt.Sprintf(":%s",webPort)
        srv.Handler = app.routes()
       err =  srv.ListenAndServe()
       if err != nil {
        log.Fatalf("unable to start server on port 80 %s",err)
       }

	
}

func connect()(*amqp.Connection, error) {
    var count int64
    var backOff = 1* time.Second
    var connection  *amqp.Connection

            //do not connect until rabbit is ready

            for {
                c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
                if err != nil {
                    fmt.Println("RabbitMQ not yet ready...")
                    count++           
                }else{
                    log.Println("connected to rabbitmq....")
                    connection = c
                    break
                }
                if count>5 {
                    fmt.Println("error connecting to rabbitMQ",err)
                    return nil, err
                }

                backOff = time.Duration(math.Pow(float64(count) ,2 )) * time.Second
                log.Printf("backing off by... : %v seconds",math.Pow(float64(count) ,2))
                time.Sleep(backOff)
                continue
            }
            return connection, nil
}