package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

var con *sql.DB

const (
	brokerAddress = "messaging-kafka.default.svc.cluster.local:9092"
	topic = "message-event"
	partion = 0
	batchSize = int(10e6)
	createTableQuery = "CREATE TABLE IF NOT EXISTS EVENTS(EmpID varchar(10), Name varchar(50), Dept varchar(10), Time varchar(100));"
	inserteQuery = "INSERT INTO EVENTS(EmpID, Name, Dept, Time) values($1, $2, $3, $4);"
)

func init() {
	con = getDBConnection()
	createTable(con)
}

type Event struct {
	Name string
	Dept string
	EmpID string
	Time string
}

const (
	host     = "db-postgresql.default.svc.cluster.local"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "event_db"
)



func ConsumeEvents() {
	fmt.Println("********* Starting Consumer ************")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerAddress},
		Topic:     topic,
		Partition: partion,
		MinBytes:  batchSize,
		MaxBytes:  batchSize,
	})

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("unable to read message ", err)
			//break
		}
		fmt.Println("\n ===================================== Event received ==========================================")
		fmt.Println(""+string(msg.Key)+ " : "+string(msg.Value))
		//var event map[string] interface{}
		//eventJson := json.Unmarshal(msg.Value, &event)
		//fmt.Println("eventJson",eventJson)
		eventObj:= &Event{}
		json.Unmarshal(msg.Value, eventObj)
		fmt.Println("Name:"+eventObj.Name+", Dept:"+eventObj.Dept+", EmpID:"+eventObj.EmpID+", Time:"+eventObj.Time)
		writeEvent(eventObj, con)
	}

	//errC := r.Close()
	//if errC != nil {
	//	log.Fatal("unable to close reader ", errC)
	//}
}

func createTable(db *sql.DB) {
	_, e := db.Exec(createTableQuery)
	if e != nil {
		log.Fatal("Unable to create database table ", e)
	}
}

func writeEvent(event *Event, db *sql.DB) {
	fmt.Println("=================== writeEvent ==================")
	fmt.Println(event)
	insertStmt := `insert into "events"("empid", "name", "dept", "time") values($1, $2, $3, $4)`
	_, e := db.Exec(insertStmt, event.EmpID, event.Name, event.Dept, event.Time)
	if e != nil {
		log.Fatal("DB error during write event")
		log.Fatal(e)
	}

}

func getDBConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to db!")
	return db
}