package main

import "fmt"

// below we can also use the [T comparable]
func printSlice[T int | string](items []T){
	for _, item:= range items{
		fmt.Println(item)
	}
}

type stack[T int | string] struct{
	elements []T
}

func main(){
	items:= []int{1,2,3,4,5,6}
	strings:= []string{"hello", "world", "how", "are","you"}

	printSlice(items)
	printSlice(strings)
	// --------------------------------------------------------------
	myStack:= stack[int]{
		elements: []int{1,2,3,4,5,6},
	}

	myStack1:= stack[string]{
		elements: []string{"go","lang"},
	}


	fmt.Println(myStack)
	fmt.Println(myStack1)
}
