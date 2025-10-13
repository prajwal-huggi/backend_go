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
func processNum(numChan chan int){
	for num:= range numChan{
		fmt.Println("Processing number", num)
		time.Sleep(time.Second)
	}
}

// receive
func sum(result chan int, num1 int, num2 int){
	res:= num1+ num2

	result<- res
}

func main(){
	// The below is the demonstration for the transferring the data from one goroutines to another
	numChan:= make(chan int)

	go processNum(numChan)

	// for{
	// 	numChan<- rand.Intn(100)
	// }
//------------------------------------------------------------------------

	result:= make(chan int)
	go sum(result, 4, 5)
	res:= <- result // blocking (sending)
	fmt.Println(res)


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
