package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/bertoxic/listener-service/events"
	ampq "github.com/rabbitmq/amqp091-go"
)

func main (){

    // try to connect to tabbitmq
    rabbitConn, err := connect()
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
   
    defer rabbitConn.Close()


    //start listening for messages

    log.Println("listening and consuming Rabbitmq messages....")


    //create consumer

    consumer, err := events.NewConsumer(rabbitConn)
    if err != nil {
        panic(err)
    }


    //watch the que and consume events
    err = consumer.Listen([]string{"log.INFO","log.WARNING", "log.ERROR"})
    if err != nil {
        log.Println(err)
    }
}



func connect()(*ampq.Connection, error) {
    var count int64
    var backOff = 1* time.Second
    var connection  *ampq.Connection

            //do not connect until rabbit is ready

            for {
                c, err := ampq.Dial("amqp://guest:guest@rabbitmq")
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
                log.Printf("backing off by... : %d seconds",math.Pow(float64(count) ,2))
                time.Sleep(backOff)
                continue
            }
            return connection, nil
}