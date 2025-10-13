package main

import (
	"fmt"
	"prajwal_huggi/backend_go/auth"
	"prajwal_huggi/backend_go/user"
)

// go mod init urlOfTheGithub

func main(){
	auth.LoginWithCredentials("prajwal", "secret")
	session:= auth.GetSession()

	fmt.Println(session)

	user:= user.User{
		Email: "user@gmail.com",
		Name: "user",
	}

	fmt.Println(user)
}
