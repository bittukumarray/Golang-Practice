package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetCustomers(t *testing.T) {

	testCases := []struct {
		inp string
		out []Customer
	}{
		{"?name=Bittu%20Ray", []Customer{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}}},
		{"", []Customer{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}, {Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}, {Id: 3, Name: "Rupesh", Dob: "05-07-2002", Address: Address{Id: 3, StreetName: "rajabazar", City: "Patna", State: "Bihar", CusId: 3}}}},
		{"name=", []Customer{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}, {Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}, {Id: 3, Name: "Rupesh", Dob: "05-07-2002", Address: Address{Id: 3, StreetName: "rajabazar", City: "Patna", State: "Bihar", CusId: 3}}}},
	}

	for i := range testCases {
		req := httptest.NewRequest("GET", "http://localhost:8080/customer"+testCases[i].inp, nil)
		w := httptest.NewRecorder()
		GetCustomersHandler(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		var cust []Customer
		err = json.Unmarshal(body, &cust)

		if err != nil {
			log.Fatal(err)
		}
		if http.StatusOK != resp.StatusCode || !reflect.DeepEqual(cust, testCases[i].out) {
			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
		}
	}

}

func TestGetCustomer(t *testing.T) {
	testCases := []struct {
		inp string
		out Customer
	}{
		{"1", Customer{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}},
		{"2", Customer{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
		{"4", Customer{}},
	}

	for i := range testCases {
		req := httptest.NewRequest("GET", "http://localhost:8080/customer/"+testCases[i].inp, nil)
		w := httptest.NewRecorder()
		GetCustomerHandler(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		var cust Customer
		err = json.Unmarshal(body, &cust)

		if err != nil {
			log.Fatal(err)
		}
		if http.StatusOK != resp.StatusCode || !reflect.DeepEqual(cust, testCases[i].out) {
			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
		}
	}

}

//func TestPostCustomer(t *testing.T) {
//	testCases := []struct {
//		inp []byte
//		out Customer
//	}{
//		{[]byte(`{"name":"Pintu","dob":"05-12-2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), Customer{Id: 3, Name: "Pintu", Dob: "05-12-2000", Address: Address{Id: 3, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 3}}},
//		//{[]byte(`{}`), Customer{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
//		//{[]byte(`{}`), Customer{}},
//	}
//
//	for i := range testCases {
//		req := httptest.NewRequest("POST", "http://localhost:8080/customer/", bytes.NewBuffer(testCases[i].inp))
//		w := httptest.NewRecorder()
//		PostCustomerHandler(w, req)
//		resp := w.Result()
//		body, err := ioutil.ReadAll(resp.Body)
//		var cust Customer
//		err = json.Unmarshal(body, &cust)
//
//		if err != nil {
//			log.Fatal(err)
//		}
//		if http.StatusCreated != resp.StatusCode || !reflect.DeepEqual(cust, testCases[i].out) {
//			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
//		}
//	}
//
//}

func TestDeleteCustomer(t *testing.T) {
	testCases := []struct {
		inp string
		out Customer
	}{
		{"3", Customer{Id: 3, Name: "Rupesh", Dob: "05-07-2002", Address: Address{Id: 3, StreetName: "rajabazar", City: "Patna", State: "Bihar", CusId: 3}}},
		//{"2", Customer{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
		{"4", Customer{}},
	}

	for i := range testCases {
		req := httptest.NewRequest("DELETE", "http://localhost:8080/customer/"+testCases[i].inp, nil)
		w := httptest.NewRecorder()
		DeleteCustomerHandler(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		var cust Customer
		err = json.Unmarshal(body, &cust)

		if err != nil {
			log.Fatal(err)
		}
		if http.StatusOK != resp.StatusCode || !reflect.DeepEqual(cust, testCases[i].out) {
			t.Errorf("FAILED!! expected %v got %v\n", testCases[i].out, cust)
		}
	}

}
