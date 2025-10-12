package main

import (
	"fmt"
	"time"
)

func main(){
	//multi condition switch

	switch time.Now().Weekday(){
		case time.Saturday, time.Sunday:
			fmt.Println("Its weekend")
		default:
			fmt.Println("Its working day")
	}

	//type switch
	whoAmI:= func(i interface{}){
		switch t:= i.(type){
			case int:
				fmt.Println("its integer")

			case string:
				fmt.Println("Its string")

			default:
				fmt.Println("Its other datatype", t)
		}
	}

	whoAmI("Praj")
}
