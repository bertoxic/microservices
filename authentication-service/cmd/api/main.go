package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bertoxic/authentication/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}
var counts int64
func main() {
    log.Println("stating authentication service on port 80")

//connect to database
conn := connectTODB()
if conn==nil {
    log.Panic("cannot connect to database")
}
//set up config
app := Config{
    DB: conn,
    Models: data.New(conn),
}

srv := &http.Server{
    Addr: fmt.Sprintf(":"+webPort),
    Handler: app.routes(),
}

err := srv.ListenAndServe()
if err != nil {
    log.Fatalf("error starting authentication server %s", err)
    return
}

    
}


func openDB(dsn string)(*sql.DB, error){

    db, err := sql.Open("pgx",dsn)
    if err != nil {
        log.Println("not connected to database oh")
        return nil,err
    }
    ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    db.PingContext(ctx)
    err = db.Ping()
    if err != nil {
        log.Fatalf("cannot ping database %v",dsn)
        return nil, err
    }

    return db, nil
}

func connectTODB() *sql.DB {
    dsn := os.Getenv("DSN")
    for {
        conenction, err := openDB(dsn)
       if err != nil {
        log.Println("postgress not ready")
        counts++
       }else{
        log.Printf("connected to database %v",dsn)
        return conenction
       }
       if counts >10 { 
       log.Println(err)
       return nil}

       log.Println("backing off for 2 seconds")
       time.Sleep(7 *time.Second)
       continue
    } 
}




//docker pull alpine:latest
//  sudo chown -R postgres:postgres /Desktop/Projects/microServices/project/db-data/postgres
// rm  ~/.docker/config.json 
// 4C4C4544-0050-5410-8036-B4C04F444332
// [guid]::NewGuid()
// wmic csproduct get uuid
