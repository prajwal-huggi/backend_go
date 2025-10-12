package main

import "fmt"

// iterating over data structure range is used
func main(){
	nums:= []int{6, 7, 8}

	for i:= 0; i< len(nums); i++{
		fmt.Println(nums[i])
	}

	// iterating with the help of range.
	for idx, num := range nums{
		fmt.Println(idx, " ",num)
	}

	// Now iterating over the map using the range
	mp:= make(map[string]int)

	mp["one"]= 1
	mp["two"]= 2
	mp["three"]= 3
	mp["four"]= 4

	for k, v:= range mp{
		fmt.Println(k, v)
	}
	
	// unicode
	// 255-> 1 byte greater than that 2 bytes will be used
	for i, c:= range "golang"{
		fmt.Println(i, string(c))
	}

	fmt.Println(len("golang"))
}
