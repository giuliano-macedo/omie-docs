package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./bundle"))
	http.Handle("/", http.StripPrefix("/", fs))

	port := 20109

	log.Print("Server ready on ", fmt.Sprint("http://localhost:", port))
	err := http.ListenAndServe(fmt.Sprint(":", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
