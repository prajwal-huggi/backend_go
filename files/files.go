package main

import (
	"fmt"
	"os"
)

func main(){
	f, err:= os.Open("files/example.txt")

	if err!= nil{
		panic(err)
	}
	
	fileInfo, err:= f.Stat()
	if err!= nil{
		panic(err)
	}

	fmt.Println(fileInfo.Name())
	fmt.Println(fileInfo.Size())

	f, err= os.Open("files/example.txt")
	if err!= nil{
		panic(err)
	}

	defer f.Close()

	buf:= make([]byte, fileInfo.Size())

	d, err:= f.Read(buf)
	if err!= nil{
		panic(err)
	}

	println("data", d, string(buf))
}
