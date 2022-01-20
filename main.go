package main

import (
	"log"
	"net/http"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB
var err error

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	PhoneNo string `json:"phone"`
	Address string `json:"address"`
}

func main() {

	db, err = sql.Open("mysql", "root:secret123@tcp(127.0.0.1:3306)/customer")
	if err != nil {
		log.Println(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/customer", getCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customer/{id}", getCustomerbyID).Methods(http.MethodGet)
	router.HandleFunc("/customer", createCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customer/{id}", updateCustomerbyID).Methods(http.MethodPut)
	router.HandleFunc("/customer/{id}", deleteCustomerbyID).Methods(http.MethodDelete)
	http.ListenAndServe(":8000", router)
}
