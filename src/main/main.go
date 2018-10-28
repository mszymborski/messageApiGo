package main

import (
	"Utils"
	"log"
	"message"
)


func main() {
	session := message.Session
	defer session.Close()
	createDB()
	handleRequest()
}

func createDB() {
	keyspace := Utils.GetConfig(Utils.KEYSPACE_KEY)
	tableName := Utils.TABLE_NAME
	err := message.Session.Query("CREATE KEYSPACE IF NOT EXISTS " + keyspace +
		" WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}").Exec()
	if err != nil{
		log.Fatal(err)
	}
	err = message.Session.Query("CREATE TABLE IF NOT EXISTS " + keyspace + "." + tableName +
		" (id UUID, email_address text, title text, content text, magic_number int," +
		" PRIMARY KEY(id, email_address, magic_number))").Exec()
	if err != nil{
		log.Fatal(err)
	}
}