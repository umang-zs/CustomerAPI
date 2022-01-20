package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/umang01-hash/CustomerAPI/driver"
	"github.com/umang01-hash/CustomerAPI/model"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

	var customers []model.Customer

	rows, err := db.Query("select * from customers")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}() // todo check rows.Err()

	for rows.Next() {
		var c model.Customer
		err := rows.Scan(&c.ID, &c.Name, &c.PhoneNo, &c.Address)
		if err != nil {
			log.Println(err) // fixme internal Server Error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		customers = append(customers, c)
	}

	res, err := json.Marshal(customers) // fixme handle error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(res))
}

func GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

	var c model.Customer

	params := mux.Vars(r)
	id := params["id"]

	err = db.QueryRow("SELECT ID, Name , phoneNo , Address FROM customers WHERE ID= ?", id).
		Scan(&c.ID, &c.Name, &c.PhoneNo, &c.Address)

	switch err {
	case sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
	case nil:
		res, err := json.Marshal(c)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(res)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// fixme: remove this

	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	var c model.Customer
	err = json.Unmarshal(body, &c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	_, err = db.Exec("INSERT INTO customers (ID,Name,phoneNo,Address) VALUES( ?, ?, ?, ? )", c.ID, c.Name, c.PhoneNo, c.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteCustomerByID(w http.ResponseWriter, r *http.Request) {

	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

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

func createQuery(id string, c model.Customer) (string, []interface{}) {
	var q []string
	var args []interface{}

	if c.Name != "" {
		q = append(q, " name=?")
		args = append(args, c.Name)
	}

	if c.Address != "" {
		q = append(q, " address=?")
		args = append(args, c.Address)
	}

	if c.PhoneNo != "" {
		q = append(q, " phoneNo=?")
		args = append(args, c.PhoneNo)
	}

	if q == nil {
		return "", args
	}

	args = append(args, id)
	query := "UPDATE customers SET" + strings.Join(q, ",") + " WHERE id = ?;"
	return query, args

}

func UpdateCustomerByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := driver.ConnectToSQL()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer db.Close()

	params := mux.Vars(r)
	id := params["id"]

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var c model.Customer
	err = json.Unmarshal(body, &c)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query, args := createQuery(id, c)
	if query == "" {
		return
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Customer with ID = %s is updated!", id)
}
