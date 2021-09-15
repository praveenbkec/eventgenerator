package eventconsumer

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var con *sql.DB

const (
	brokerAddress = "messaging-kafka.default.svc.cluster.local:9092"
	topic = "message-event"
	partion = 0
	batchSize = int(10e6)
	host     = "db-postgresql.default.svc.cluster.local"
	port     = 5432
	user     = "admin"
	password = "admin"
	dbname   = "event_db"
	createTableQuery = "CREATE TABLE IF NOT EXISTS EVENTS(EmpID varchar(10), Name varchar(50), Dept varchar(10), Time varchar(100));"
	inserteQuery = "INSERT INTO EVENTS(EmpID, Name, Dept, Time) values($1, $2, $3, $4);"
)

type EventMgmtIf interface {
	CreateEvent();
	UpdateEvent();
	GetEvent();
	ListEvent();
}

type EventMgmtStruct struct {

}

type EventRequest struct {
	Name string
	Dept string
	EmpID string
	Time string
}

type EventResponse struct {
	status string
	status_message string
}

func init() {
	con = getDBConnection()
	createTable(con)
}

func (e *EventMgmtStruct) CreateEvent(ctx context.Context, req *EventRequest) (*EventResponse, error) {
	fmt.Println("******** CreateEvent starts ********** ")
	resp := &EventResponse{}
	err := writeEventToDB(req, con)
	if err != nil {
		resp.status = "Failure"
		resp.status_message = "Unable to store in database"
		return resp, nil
	}
	resp.status = "Success"
	fmt.Println("******** CreateEvent ends ********** ")
	return resp, nil
}

func (e *EventMgmtStruct) UpdateEvent(ctx context.Context, req *EventRequest) (*EventResponse, error) {
	fmt.Println("******** UpdateEvent starts ********** ")
	resp := &EventResponse{}
	err := updateEventFromDB(req, con)
	if err != nil {
		resp.status = "Failure"
		resp.status_message = "Unable to update in database"
		return resp, nil
	}
	resp.status = "Success"
	fmt.Println("******** UpdateEvent ends ********** ")
	return resp, nil
}

func (e *EventMgmtStruct) GetEvent(ctx context.Context, req *EventRequest) (*EventRequest, error) {
	fmt.Println("******** GetEvent starts ********** ")
	event, err := getEventFromDB(req.EmpID, con)
	if err != nil {
		return nil, err
	}
	fmt.Println("******** GetEvent ends ********** ")
	return event, nil
}

func (e *EventMgmtStruct) ListEvent(ctx context.Context, req *EventRequest) ([]*EventRequest, error) {
	fmt.Println("******** ListEvent starts ********** ")
	events, err := getAllEventFromDB(con)
	if err != nil {
		return nil, err
	}
	fmt.Println("******** ListEvent ends ********** ")
	return events, nil
}


// ALL DB RELATED FUNCTIONS STARTS HERE

func createTable(db *sql.DB) {
	_, e := db.Exec(createTableQuery)
	if e != nil {
		log.Fatal("Unable to create database table ", e)
	}
}

func writeEventToDB(event *EventRequest, db *sql.DB) error {
	fmt.Println("=================== writeEventToDB ==================")
	fmt.Println(event)
	insertStmt := `insert into "events"("empid", "name", "dept", "time") values($1, $2, $3, $4)`
	_, e := db.Exec(insertStmt, event.EmpID, event.Name, event.Dept, event.Time)
	if e != nil {
		log.Fatal("DB error during write event")
		log.Fatal(e)
		return e
	}
	return nil
}

func updateEventFromDB(event *EventRequest, db *sql.DB) error {
	fmt.Println("=================== updateEventFromDB ==================")
	fmt.Println(event)
	updateStmt := `update "events" set "name"=$2, "dept"=$3, "time"=$4 where "empid"=$1`
	res, err := db.Exec(updateStmt, event.EmpID, event.Name, event.Dept, event.Time)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("rows affecdted %v",count)
	return nil
}

func getEventFromDB(empid string, db *sql.DB) (*EventRequest, error) {
	fmt.Println("=================== getEventFromDB ==================")
	resp := &EventRequest{}
	selectStmt := `select "empid", "name", "dept", "time" from "events" where "empid"=$1`
	row := db.QueryRow(selectStmt, empid)
	switch err := row.Scan(&resp.EmpID, &resp.Name, &resp.Dept, &resp.Time); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(resp.EmpID, resp.Name)
	default:
		panic(err)
	}
	return resp, nil
}

func getAllEventFromDB(db *sql.DB) ([]*EventRequest, error) {
	fmt.Println("=================== getAllEventFromDB ==================")
	resp := [] *EventRequest{}
	selectStmt := `select "empid", "name", "dept", "time" from "events"`
	rows, err := db.Query(selectStmt)
	if err != nil {
		log.Fatal("Failed to execute query: ", err)
		return nil, err
	}
	for rows.Next() {
		event := &EventRequest{}
		rows.Scan(&event.EmpID, &event.Name, &event.Dept, &event.Time)
		resp = append(resp, event)
	}
	return resp, nil
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