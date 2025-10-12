package main

import "fmt"

func add(a, b int)int{
	return a+ b
}

func getLanguages()(string, bool, int){
	return "golang", true, 1
}

func processIt(fn func(a int) int){
	fmt.Println(fn(1))
}

func processingIt() func(a int) int{
	return func(a int) int{
		return a
	}
}

func main(){
	fmt.Println(add(3, 5))

	fmt.Println(getLanguages())

	fn:= func(a int) int{
		return a
	}

	processIt(fn)

	fnt:= processingIt()
	fmt.Println(fnt(9))
}
