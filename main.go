package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", handleRegister).Methods("POST")
	r.HandleFunc("/login", handleLogin).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
