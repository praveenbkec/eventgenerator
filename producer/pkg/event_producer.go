package pkg

import "time"

type Event struct {
	Name string
	Dept string
	EmpID string
	Time string
}
// EmpID and Event Details
var EmpDB = make(map[string]*Event)
var TrainEventDB = make(map[string]*Event)

const (
	EmployeeAccessEventConst = "EmployeeAccessEvent"
	TrainTickBookingEventConst = "TrainTickBookingEvent"
)


func init() {
	EmpDB["12345"] = &Event{EmpID: "12345", Name: "Praveen BK", Dept: "IT"}
	EmpDB["23456"] = &Event{EmpID: "23456", Name: "Vinaya", Dept: "Telecom"}
	EmpDB["34567"] = &Event{EmpID: "34567", Name: "Adithya", Dept: "Cloud"}
	EmpDB["45678"] = &Event{EmpID: "45678", Name: "Gireesh", Dept: "Platform"}
}

type EventProducerIf interface {
	ProduceEvent()
}

type EmployeeAccessEvent struct {
	EventType string
	EmpID     string
}

func (e EmployeeAccessEvent) ProduceEvent() (*Event, error) {
	if e.EventType == EmployeeAccessEventConst {
		event := EmpDB[e.EmpID]
		event.Time = time.Now().Format(time.RFC850)
		return event, nil
	} else if e.EventType == TrainTickBookingEventConst {
	}

	return nil, nil
}

func main() {
	eae := &EmployeeAccessEvent{EventType: EmployeeAccessEventConst, EmpID: "12345"}
	e, _ := eae.ProduceEvent()
	println(e)

}