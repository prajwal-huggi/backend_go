package main

import (
	"fmt"
	"maps"
)

func main(){
	//creating the map	
	m:= make(map[string]string)

	// setting an element
	m["name"]= "golang"

	fmt.Println(m["name"])

	fmt.Println(len(m))

	// 2nd way to create the map
	mp:= map[string]bool{"intelligent":true, "ambitious":true}

	fmt.Println(mp)

	k, ok:= mp["ambitious"]// k returns the value of the key and ok returns whether the key is present or not.
	fmt.Println(k)

	if ok{
		fmt.Println("all ok")
	}else{
		fmt.Println("not ok")
	}

	// To check the equality between 2 maps
	m1:= map[string]int{"one":1, "two": 2}
	m2:= map[string]int{"one":1, "two": 2}

	fmt.Println(maps.Equal(m1, m2))
}
