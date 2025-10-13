package main

import (
	"fmt"
	"time"
)

// order struct
type order struct {
	id        string
	amount    float32
	status    string
	createdAt time.Time // nanosecond precision
}

// To build the constructor we use the trick. For the naming conventino we are using the new keyword
func newOrder(id string, amount float32, status string) *order{
	myOrder := order{
		id:        id,
		amount:    amount,
		status:    status,
	}

	return &myOrder
}

// reciever type
func (o *order) changeStatus(status string){
	o.status= status
}

func (o order)getAmount()float32{
	return o.amount
}

/*
The standard convention in the struct is that
- If you want to modify the value of the struct then use *
- If you want to just get the value of the struct then don't use *
*/

func main() {
	// If you don't set any field, default value is zero value
	// int-> 0, float-> 0, string-> "", bool-> false
	myOrder := order{
		id:        "1",
		amount:    50,
		status:    "pending",
		createdAt: time.Now(),
	}

	fmt.Println(myOrder)
	myOrder.changeStatus("paid")
	fmt.Println(myOrder)

	fmt.Println(myOrder.getAmount())

	o1:= newOrder("10", 30.50, "Successful")

	fmt.Println(o1)

	// If you want to use the struct only for the first time then in that case we can use the inline struct
	language:= struct{
		name string
		isGood bool
	}{"golang", true}

	fmt.Println(language)
}
