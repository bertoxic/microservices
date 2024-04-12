package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	//"path/filepath"
	"strconv"
)

type Config struct {
    Mailer Mail
}

const webport = "80"

func main() {
	app := Config{
        Mailer: createMail(),
    }
	log.Println("stating mail service on port ", webport)

	srv := &http.Server{
        Addr: fmt.Sprintf(":%s",webport),
        Handler: app.routes(),
    }
    err:=srv.ListenAndServe()
    if err != nil {
        log.Panic()
    }
    // dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    // if err != nil {
    //         log.Fatal(err)
    // }
    // fmt.Println("xuiuhaaaaaaaaa",dir)
}


func createMail() Mail{
    port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
    m := Mail{
        Domain: os.Getenv("MAIL_DOMAIN"),
        Host: os.Getenv("MAIL_HOST"),
        Port: port,
        Username: os.Getenv("MAIL_USERNAME"),
        Password: os.Getenv("MAIL_PASSWORD"),
        Encryption: os.Getenv("MAIL_ENCRYPTION"),
        FromName: os.Getenv("FROM_NAME"),
        FromAddress: os.Getenv("FROM_ADDRESS"),
    }

    return m
}