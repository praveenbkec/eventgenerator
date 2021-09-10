package main

import (
	"fmt"
	"time"
)

func main() {
	produceEvents()
}

func produceEvents() {
	for(true)  {
		fmt.Println("User looged into station at "+time.Now().String())
		time.Sleep(5*time.Second)
	}
}
