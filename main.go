package main

import (
	"database/sql"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"

	_ "github.com/lib/pq"
)



func main() {
	config,err := util.LoadConfig(".")
	if err!=nil{
		log.Fatal("load config err: ",err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect database:", err)
	}
	store :=db.NewStore(conn)
	server,err := api.NewServer(config,store)
	if err!=nil{
		log.Fatal("cannot create sever: ",err)
	}
	err = server.Start(config.ServerAddress)
	if err !=nil{
		log.Fatal("cannot start server:",server)
	}
}