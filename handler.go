package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

func getCustomers(w http.ResponseWriter, r *http.Request) {
	//telling the client that the response sent is in json format
	w.Header().Set("Content-Type", "application/json")

	var customers []Customer

	rows, err := db.Query("select * from customers")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var c Customer
		err := rows.Scan(&c.ID, &c.Name, &c.PhoneNo, &c.Address)
		if err != nil {
			log.Println(err)
		}
		customers = append(customers, c)
	}

	res, _ := json.Marshal(customers)
	w.Write([]byte(res))
}

func getCustomerbyID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var customers []Customer
	params := mux.Vars(r)
	id := params["id"]

	rows, err := db.Query("SELECT ID, Name , phoneNo , Address FROM customers WHERE ID= ?", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var c Customer

		err := rows.Scan(&c.ID, &c.Name, &c.PhoneNo, &c.Address)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		customers = append(customers, c)

	}

	res, _ := json.Marshal(customers)

	if string(res) == "null" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No customer with ID = %s found! ", id)

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func createCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res, err := db.Prepare("INSERT INTO customers (ID,Name,phoneNo,Address) VALUES( ?, ?, ?, ? )")

	if err != nil {
		log.Println(err)
	}

	defer res.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	var c Customer
	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	_, err = res.Exec(c.ID, c.Name, c.PhoneNo, c.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func deleteCustomerbyID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	_, err = db.Exec("DELETE FROM customers WHERE ID = ?", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Customer with ID = %s was deleted", id)

}

//func updateCustomer(w http.ResponseWriter, r *http.Request) {
//	params := mux.Vars(r)
//	id := params["id"]
//	feedData, err := db.Prepare("UPDATE customers SET Name=? , phoneNo=? , Address=? where id=?")
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//	defer feedData.Close()
//
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		log.Println(err.Error())
//
//	}
//
//	var c customer
//	err = json.Unmarshal(body, &c)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//
//	_, err = feedData.Exec(c.Name, c.PhoneNo, c.Address, id)
//	if err != nil {
//		log.Println(err.Error())
//		return
//	}
//
//	fmt.Fprintf(w, "Customer with ID = %s Updated!", id)
//
//}

func createQuery(id string, c Customer) string {
	query := "UPDATE customers SET "
	var q []string
	if c.Name != "" {
		q = append(q, " name = \""+c.Name+"\"")
	}
	if c.PhoneNo != "" {
		q = append(q, " phoneNo = \""+c.PhoneNo+"\"")
	}
	if c.Address != "" {
		q = append(q, " address = \""+c.Address+"\"")
	}

	if q == nil {
		return ""
	}

	query += strings.Join(q, " , ")

	query += " where ID = " + string(id) + " ; "

	return query

}

func updateCustomerbyID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var c Customer
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := createQuery(id, c)
	if query == "" {
		return
	}

	_, err = db.Exec(query)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Customer with ID = %s is updated!", id)

}
