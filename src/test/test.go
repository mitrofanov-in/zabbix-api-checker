package main

import (
	"fmt"
	"getenv"
	//	"os"
)

func main() {

	custUser := getenv.GoDotEnvVariable("customerUser")
	custPass := getenv.GoDotEnvVariable("customerPassword")

	fmt.Println(custPass, custUser)
}
