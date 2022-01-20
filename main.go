package main

import (
	"github.com/umang01-hash/CustomerAPI/handler"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/customer", handler.GetCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customer/{id}", handler.GetCustomerByID).Methods(http.MethodGet)
	router.HandleFunc("/customer", handler.CreateCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customer/{id}", handler.UpdateCustomerByID).Methods(http.MethodPut)
	router.HandleFunc("/customer/{id}", handler.DeleteCustomerByID).Methods(http.MethodDelete)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(srv.ListenAndServe())
}
