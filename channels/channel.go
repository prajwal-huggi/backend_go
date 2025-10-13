package main

import (
	"fmt"
	"time"
)

/*
When there are multiple goroutines are executing then we have to exchange the
data between various goroutines then we are using the goroutines.
*/

// send
func processNum(numChan chan int) {
	for num := range numChan {
		fmt.Println("Processing number", num)
		time.Sleep(time.Second)
	}
}

// receive
func sum(result chan int, num1 int, num2 int) {
	res := num1 + num2

	result <- res
}

// receiving the channel: goroutines synchronizer
func task(done chan bool) {
	defer func() { done <- true }()

	fmt.Println("processing...")
}

func emailSender(emailChan chan string, done chan bool) {
	defer func() { done <- true }()

	for email := range emailChan {
		fmt.Println("Sending email to ", email)
		time.Sleep(time.Second)
	}
}

func main() {
	// The below is the demonstration for the transferring the data from one goroutines to another
	numChan := make(chan int)

	go processNum(numChan)

	// for{
	// 	numChan<- rand.Intn(100)
	// }
	//------------------------------------------------------------------------

	// result:= make(chan int)
	// go sum(result, 4, 5)
	// res:= <- result // blocking (sending)
	// fmt.Println(res)

	//------------------------------------------------------------------------

	// done:= make(chan bool)

	// go task(done)
	// <- done //block until the value is recieved from the line 30

	//------------------------------------------------------------------------
	// Buffered Channel
	emailChan := make(chan string, 100)
	done := make(chan bool)

	go emailSender(emailChan, done)
	for i:= 0; i< 100; i++{
		emailChan<- fmt.Sprintf("%d@gmail.com", i+ 1)
	}
	fmt.Println("done sending.")
	// emailChan <- "1@gmail.com"
	// emailChan <- "2@gmail.com"

	// fmt.Println(<-emailChan)
	// fmt.Println(<-emailChan)

	close(emailChan) // This is used to close the channel to avoid the deadlock. Very important
	<- done

	//------------------------------------------------------------------------

	/*
		messageChan:= make(chan string)

		// inserting the value in the channel or we are sending the data in the channel
		messageChan<- "ping" // blocking

		// to receive the data from the channel
		msg:= <- messageChan

		fmt.Println(msg)
	*/
}
