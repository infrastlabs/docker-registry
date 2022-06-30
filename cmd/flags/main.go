package main

import (
	"flag"
	"fmt"
)

var (
	user = flag.String("user", "admin", "registry user")
	pass = flag.String("pass", "admin123", "registry pass")
	size = flag.Bool("size", true, "registry size") //or: req's param
)

func main(){
	flag.Parse()

	fmt.Println("user,pass,size: "+*user +" "+*pass)
	// $ go run ./cmd/flags/main.go -user=asdf
	// $ go run ./cmd/flags/main.go --user=asdf
}