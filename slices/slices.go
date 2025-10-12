// Slices is like the vector of cpp
package main

import (
	"fmt"
	"slices"
)

func main(){
	//uninitialized slice is nil
	// var nums[] int
	// fmt.Println(nums== nil)

	// var nums= make([]int, 2, 5)//2 is the initial size and the 5 is the initial capacity
	//capacity-> maximum number of elements that can be fitted in the slice
	// fmt.Println(cap(nums))
	// fmt.Println(nums)

	// nums= append(nums, 1)
	// fmt.Println(cap(nums))
	// fmt.Println(nums)

	// nums:= []int{}

	//slice operator
	// var nums= []int{1, 2, 3}
	// fmt.Println(nums[0:2])

	//slices package
	var nums1= []int {1, 2, 3}
	var nums2= []int {1, 2, 3}

	fmt.Println(slices.Equal(nums1, nums2))
}
