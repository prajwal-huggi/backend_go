package main

import (
	"fmt"
	"time"
)

/*
goroutines are the light weight thread which are used for multi threading and running the multiple things concurrently
*/

func task(id int){
	fmt.Println("doing task", id)
}

func main() {
	for i:= 0; i<= 10; i++{
		go task(i)
	}

	time.Sleep(time.Second* 2)
}
