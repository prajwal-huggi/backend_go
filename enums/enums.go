package main

import "fmt"

// There is no inbiult enum in the go so for that we are using the const
type OrderStatus int

const (
	Received OrderStatus = iota
	Confirmed
	Prepared
)

func changeOrderStatus(status OrderStatus) {
	fmt.Println("changing order status to", status)
}

func main() {
	changeOrderStatus(Received)
}
