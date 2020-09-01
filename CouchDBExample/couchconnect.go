package main

import (
	"fmt"
	"github.com/leesper/couchdb-golang"
	"log"
)

func main() {
	err := connectToCouchDB()

	if err != nil {
		log.Println("failed:", err.Error())
	}
}

func connectToCouchDB() error {
	server, err := couchdb.NewServer("http://10.0.0.10:5984/_utils")

	if err != nil {
		log.Println("Failed to connect to CouchDB:", err.Error())
		return err
	}

	//authToken, err := server.Login("admin", "admin")
	//
	//if err != nil {
	//	log.Println("failed to login:", err.Error())
	//	return err
	//}
	//
	//fmt.Println(authToken)



	db, err := couchdb.NewDatabase("monitutor-results")

	if err != nil {
		log.Println("failed to create new db instance")
		return err
	}
	db1, err := server.Get(db.String())

	if err != nil {
		log.Println("failed to get DB:", err.Error())
		return err
	}

	//dbs, err := server.DBs()
	//
	//if err != nil {
	//	log.Println("failed to get DBs:", err.Error())
	//	return err
	//}
	//
	//for db2 := range dbs {
	//	fmt.Println(db2)
	//}
	err = db1.Contains("mv6889s_itsclient")

	if err == nil {
		fmt.Println("contains")
	} else {
		fmt.Println("not contains")
	}

	//fmt.Println(db1.String())
	return nil
}

//TODO: erst ab CouchDB 2.x!!!