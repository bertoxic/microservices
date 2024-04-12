package events

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

 type Emitter struct {
    connection *amqp.Connection
 }


 func (e *Emitter) SetUP() error{
    channel, err := e.connection.Channel()
    if err != nil {
        fmt.Printf("error in broker emitter function setup*: %s", err.Error())
        return err
    }

        defer channel.Close()

        return declareExchange(channel)
 }


 func (e *Emitter) Push (event, severity string) error {

    channel, err := e.connection.Channel()
    if err != nil {
        fmt.Printf("error in broker emitter function push* : %s", err.Error())
        return err
    }

   defer channel.Close()

      err=  channel.Publish(
            "logs_topic",
            severity,
            false,
            false,
            *&amqp.Publishing{
                ContentType: "text/plain",
                Body: []byte(event),
            },
        )

        if err != nil {
            fmt.Printf("error in broker emitter function push* : %s", err.Error())
            return err
        }

        return nil
 }

 func NewEmittter (conn *amqp.Connection) (Emitter , error) {
            emitter := Emitter{
                connection:  conn,
            }
             err := emitter.SetUP()
             if err != nil {
                fmt.Printf("error in broker emitter function NewEmitter* : %s", err.Error())
                return Emitter{}, err
             }

             return emitter, nil
 }