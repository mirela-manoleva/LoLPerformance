package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Hello World!")
	resp, err := http.Get("http://google.com/")
	if err != nil {
		fmt.Println("Error: ", err)
	} else {
		fmt.Println(resp.StatusCode, resp)
	}
}
