package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Chat Service Started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
