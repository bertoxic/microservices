package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)



type Consumer struct {
    conn *amqp.Connection
    queueName string
}


func NewConsumer(conn *amqp.Connection)(Consumer, error){

    consumer := Consumer{
        conn:conn,
    }
    err := consumer.setup()
    if err != nil {
        return Consumer{}, err
    }

    return consumer, nil
}


func (consumer *Consumer) setup()error {
    channel, err := consumer.conn.Channel()
    if err != nil {
        return err
    }

    return declareExchange(channel)
    
}

type PayLoad struct {
    Name string `json:"name"`
    Data string `json:"data"`
}

func (consumer *Consumer) Listen (topics []string) error {
    ch, err := consumer.conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()
    q, err := declareRandomQueue(ch)
    if err != nil {
        return err
    }

    for _, s := range topics {
            ch.QueueBind(
                q.Name,
                s,
                "logs_topic",
                false,
                nil,
            )
            if err != nil {
                return err
            }
    }

    messages , err := ch.Consume(q.Name, "", true, false, false, false, nil)
    if err != nil {
        return err
    }

    forever := make(chan bool)
    go func(){
        for d := range messages {
            var payLoad PayLoad
            _ = json.Unmarshal(d.Body,&payLoad)
            go handlePayLoad(payLoad)
            fmt.Printf("handling payload for message on , %s]\n", payLoad)
        }
    }()
    fmt.Printf("waiting for message on [Exchange, Queue] [logs_topic, %s]\n", q.Name)
    <- forever

    return nil
}

func handlePayLoad(payload PayLoad) {
    switch payload.Name {
    case "log", "event":
        //log what is being sent
        err := logEvent(payload)
        if err != nil {
            log.Println(err)
                }
    case "auth" :
        //authenticate


    default :
    err := logEvent(payload)
    if err != nil {
        log.Println(err)
            }
    }
  
}

func logEvent( entry PayLoad) error {
   
        //createjson to send to auth microservice
        log.Println("data gotten is", entry.Name, "and...", entry.Data)
        jsonData, err := json.MarshalIndent(entry, "", "\t")
        if err != nil {
            println("Error marshalling data zzzzzzzzzzzzzz:", err)
            return err
        }
        logServiceURL := "http://logger-service/log"
        println("gotten to request now wooow", jsonData)
        //call the service
        request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
        if err != nil {
            println("eroooosss in making request in boker handler vvv:")
        }
        request.Header.Set("Content-Type", "application/json")
        client := &http.Client{}
    
        response, err := client.Do(request)
        
    
        if err != nil {
            println("eroooosss in sending through client oh respone vvv:", err.Error())
            return err
        }
        body, err := ioutil.ReadAll(response.Body)
        if err != nil {
            println("Error reading response body:", body, err)
            // Handle the error appropriately (e.g., return an error to the caller)
            return err
        }
        if response != nil {
            defer response.Body.Close()
        }
        if response.StatusCode != http.StatusAccepted {
            
            return err
        }
    
        
    
        return nil    
    
}