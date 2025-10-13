package main

import (
	"fmt"
	"sync"
)

/*
goroutines are the light weight thread which are used for multi threading and running the multiple things concurrently
*/

func task(id int, w *sync.WaitGroup){
	defer w.Done()
	fmt.Println("doing task", id)
}

func main() {
	var wg sync.WaitGroup

	for i:= 0; i<= 10; i++{
		wg.Add(1)
		go task(i, &wg)
	}

	wg.Wait()
}
// add-> done -> wait

/*
when the wg is 0 then and only then the entire program will be executed.
so when we are diong the wg.Add(1) we are adding the one and in the function
when we are doing the defere then in that case we are decrementing the 1 from it
suppose the entire program is went to wg.Wait() then it will check that if the value of the
wg is 0 then and only then the main function will be executed otherwise it will be not be executed.
*/
