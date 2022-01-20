package main

import (
	"bytes"
	"database/sql"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetCustomer(t *testing.T) {
	db, err = sql.Open("mysql", "root:secret123@tcp(127.0.0.1:3306)/customer")
	if err != nil {
		log.Println(err.Error())
	}

	defer db.Close()
	cases := []struct {
		desc string
		id   string // input
		resp []byte
	}{
		{"id exists in db", "1", []byte(`[{"id":1,"name":"Manav","phone":"9953113063","address":"Shakti Khand 2 Ghaziabad UP"}]`)},
	}

	for i, tc := range cases {
		req := httptest.NewRequest("GET", "/getCustomer", nil)

		r := mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()

		getCustomerbyID(w, r)

		resp := w.Result()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read response %v", err)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, string(body), string(tc.resp))
		}

	}

}

func TestCreateCustomer(t *testing.T) {
	db, err = sql.Open("mysql", "root:secret123@tcp(127.0.0.1:3306)/customer")
	if err != nil {
		log.Println(err.Error())
	}

	defer db.Close()

	cases := []struct {
		desc       string
		customer   []byte
		statusCode int
	}{
		{"New Customer Created", []byte(`{"id":5,"name":"rakshit","phone":"852927395729","address":"Lucknow Up"}`), http.StatusCreated},
	}

	for i, tc := range cases {
		req := httptest.NewRequest("POST", "/getCustomer", bytes.NewReader(tc.customer))
		w := httptest.NewRecorder()

		createCustomer(w, req)

		resp := w.Result()

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}

	}

}

func TestUpdateCustomer(t *testing.T) {

	db, err = sql.Open("mysql", "root:secret123@tcp(127.0.0.1:3306)/customer")
	if err != nil {
		log.Println(err.Error())
	}

	defer db.Close()

	cases := []struct {
		desc       string
		id         string // input
		body       []byte
		resp       []byte
		statusCode int
	}{
		{"id exists in db", "1", []byte(`{"name":"Manav","phone":"9953113063"}`), []byte(`Customer with ID = 1 is updated!`), http.StatusCreated},
	}

	for i, tc := range cases {
		req := httptest.NewRequest("PUT", "/getCustomer", bytes.NewReader(tc.body))
		r := mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()

		updateCustomerbyID(w, r)

		resp := w.Result()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read response %v", err)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, string(body), string(tc.resp))
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}

	}

}

func TestDeleteCustomer(t *testing.T) {
	db, err = sql.Open("mysql", "root:secret123@tcp(127.0.0.1:3306)/customer")
	if err != nil {
		log.Println(err.Error())
	}

	defer db.Close()

	cases := []struct {
		desc       string
		id         string // input
		resp       []byte
		statusCode int
	}{
		{"id exists in db", "2", []byte(`Customer with ID = 2 was deleted`), http.StatusOK},
	}

	for i, tc := range cases {
		req := httptest.NewRequest("PUT", "/getCustomer", nil)
		r := mux.SetURLVars(req, map[string]string{"id": tc.id})
		w := httptest.NewRecorder()

		deleteCustomerbyID(w, r)

		resp := w.Result()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("cannot read response %v", err)
		}

		if !reflect.DeepEqual(body, tc.resp) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, string(body), string(tc.resp))
		}

		if resp.StatusCode != tc.statusCode {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i, resp.StatusCode, tc.statusCode)
		}

	}

}
