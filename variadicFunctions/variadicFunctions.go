package main

import "fmt"

// The below function is the example of the variadicFunctions.
// Here if we want to set the data type as any then use interface{}
func sum(nums ...int) int{
	total:= 0
	for _, num:= range nums{
		total+= num
	}

	return total
}

func main(){
	// Here fmt is the variadic function as we can send n number of elements in it
	fmt.Println(1, 2, 3, 4, "hello")

	nums:= []int{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Println(sum(nums...))
}
