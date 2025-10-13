package main

import "fmt"
/*
	- for the interface we have to add the er at the suffix.
	- It is acting as the contractor which tells us that it should have to function metioned inside it.
	- In other langauges we generally use the implements keyword
	- but here it is checking that if there is any signature function mentioned in any other struct then it will implicitly make changes.
*/
type paymenter interface{
	pay(amount float32)
}

type payment struct{
	gateway paymenter
}

func (p payment) makePayment(amount float32){
	p.gateway.pay(amount)
}

type razorpay struct{}

func (r razorpay) pay(amount float32){
	//logic to make payment
	fmt.Println("making payment using razorpay", amount)
}

type stripe struct{}

func (s stripe) pay(amount float32){
	fmt.Println("making payment using stript", amount)
}

func main(){
	stripeGw:= stripe{}
	// razorpayGw:= razorpay{}

	newPayment:= payment{
		// gateway: razorpayGw,
		gateway: stripeGw,
	}
	newPayment.makePayment(100)
}
