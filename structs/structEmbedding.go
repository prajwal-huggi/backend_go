package main

import (
	"fmt"
	"time"
)

type customer struct {
	name  string
	phone string
}

type order struct {
	id        string
	amount    float32
	status    string
	createdAt time.Time
	customer
}

func main() {
	newCustomer := customer{
		name:  "Prajwal",
		phone: "9427063889",
	}

	newOrder := order{
		id:       "1",
		amount:   30,
		status:   "paid",
		customer: newCustomer,
	}

	fmt.Println(newOrder)

}
