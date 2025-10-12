package main

import "fmt"

func changeNum(num *int){
	*num= 5
	fmt.Println("In changeNum", *num)
}

func main(){
	// num:= 1
	var num int= 1
	changeNum(&num)

	fmt.Println("Main function", num)
}
