package main

import (
	"bytes"
	"log"

	"github.com/DATA-DOG/go-sqlmock"

	//"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandle(t *testing.T) {
	testCases := []struct {
		req *http.Request
		out string
	}{
		{httptest.NewRequest("GET", "localhost:8080/", nil), "GET - Hello"},
		{httptest.NewRequest("POST", "localhost:8080/", nil), "POST - Hello"},
	}

	for ind := range testCases {
		w := httptest.NewRecorder()
		handle(w, testCases[ind].req)
		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		if http.StatusOK == resp.StatusCode && string(body) != testCases[ind].out {
			t.Errorf("FAILED for %v expected %v got %v\n", testCases[ind].req.Method, testCases[ind].out, string(body))
		}
	}
}

func TestPrintHandle(t *testing.T) {
	testCases := []struct {
		req *http.Request
		out string
	}{
		{httptest.NewRequest("GET", "http://localhost:8080/print/bittu", nil), "Printing bittu"},
		{httptest.NewRequest("GET", "http://localhost:8080/print/helloworld", nil), "Printing helloworld"},
	}

	for ind := range testCases {
		w := httptest.NewRecorder()
		PrinterHandler(w, testCases[ind].req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		if string(body) != testCases[ind].out {
			t.Errorf("FAILED for %v expected %v got %v\n", testCases[ind].req.Method, testCases[ind].out, string(body))
		}
	}

}

func TestCustomerHandle(t *testing.T) {
	testCases := []struct {
		inp []byte
		out string
	}{
		{[]byte(`{"name":"bittu","age":20,"address":"patna"}`), string([]byte(`{"name":"bittu","age":20,"address":"patna"}`))},
		{[]byte(`{"name":"bittu ray","age":25,"address":"     "}`), string([]byte(`{"name":"bittu ray","age":25,"address":"     "}`))},
		{[]byte(`{"name":"bittu ray","age":18,"address":"Bihar Patna"}`), string([]byte(`{"name":"bittu ray","age":18,"address":"Bihar Patna"}`))},
		{[]byte(`{"name":"bittu","age":16,"address":"patna"}`), "You are not eligible to be our customer"},
		{[]byte(`{"name":"bittu","address":"patna"}`), "You are not eligible to be our customer"},
	}

	for ind := range testCases {

		req := httptest.NewRequest("POST", "http://localhost:8080/customer", bytes.NewBuffer(testCases[ind].inp))
		w := httptest.NewRecorder()
		CustomerHandler(w, req)
		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error %v", err)
		}
		if http.StatusOK != resp.StatusCode || string(body) != testCases[ind].out {
			t.Errorf("FAILED for %v expected %v got %v\n", string(testCases[ind].inp), testCases[ind].out, string(body))
		}
	}

}

func TestCustomerAddressHandle(t *testing.T) {

	testCases := []struct {
		inp string
		out []Customers
	}{
		{"1", []Customers{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 9, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}}},
		{"0", []Customers{{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 9, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}, {Id: 2, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 15, StreetName: "ram Nagar", City: "Ara", State: "Bihar", CusId: 2}}, {Id: 4, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 10, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 4}}, {Id: 6, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 3, StreetName: "Raja Bajar", City: "Patna", State: "Bihar", CusId: 6}}, {Id: 6, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 5, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 6}}, {Id: 6, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 7, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 6}}, {Id: 6, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 12, StreetName: "ram Nagar", City: "Ara", State: "Bihar", CusId: 6}}, {Id: 7, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 11, StreetName: "ram Nagar", City: "Ara", State: "Bihar", CusId: 7}}, {Id: 7, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 13, StreetName: "ram Nagar", City: "Ara", State: "Bihar", CusId: 7}}, {Id: 8, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 6, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 8}}, {Id: 8, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 8, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 8}}, {Id: 8, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 14, StreetName: "ram Nagar", City: "Ara", State: "Bihar", CusId: 8}}, {Id: 9, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 4, StreetName: "Raja Bajar", City: "Patna", State: "Bihar", CusId: 9}}, {Id: 10, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 1, StreetName: "Itadhiyan", City: "Bikramganj", State: "Patna", CusId: 10}}, {Id: 10, Name: "Rupesh Raj", Dob: "10-02-1997", Address: Address{Id: 2, StreetName: "Itadhiyan", City: "Bikramganj", State: "Patna", CusId: 10}}}},
		{"222", []Customers(nil)},
		{"1 or 1=1", []Customers(nil)},
	}

	for ind := range testCases {
		data := GetCustomer(db, testCases[ind].inp)
		if !reflect.DeepEqual(data, testCases[ind].out) {
			t.Errorf("FAILED for %v expected %v got %v\n", testCases[ind].inp, testCases[ind].out, data)
		}
	}

}

func TestCreate(t *testing.T) {
	query1 := "drop table addrs"
	query2 := "drop table cust"
	query3 := "create table cust(id int not null auto_increment, name varchar(255), dob varchar(255), primary key(id))"
	query4 := "create table addrs(id int not null auto_increment, streetname varchar(255), city varchar(255), state varchar(255),cus_id int, primary key (id), foreign key (cus_id) references cust(id))"

	_, err := db.Query(query1)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Query(query2)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Query(query3)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Query(query4)
	if err != nil {
		log.Fatal(err)
	}
	testCases := []struct {
		inp Customers
		out Customers
	}{
		{Customers{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}, Customers{Id: 1, Name: "Bittu Ray", Dob: "15-02-1998", Address: Address{Id: 1, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 1}}},
		{Customers{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}, Customers{Id: 2, Name: "Bittu Kumar Ray", Dob: "15-02-1998", Address: Address{Id: 2, StreetName: "Boring road", City: "Patna", State: "Bihar", CusId: 2}}},
	}

	for ind := range testCases {
		data := CreateCustomer(db, testCases[ind].inp)
		if !reflect.DeepEqual(data, testCases[ind].out) {
			t.Errorf("FAILED for %v expected %v got %v\n", testCases[ind].inp, testCases[ind].out, data)
		}
	}

}

func TestCustOnly(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "dob"}).
		AddRow(1, "Bittu Ray", "15-02-1998").
		AddRow(2, "Bittu Ray", "15-02-1999")

	mock.ExpectQuery("^SELECT id, name, dob FROM customers$").WillReturnRows(rows)
	//fmt.Println("a is ", a)
	var c Customers
	if c, err = getCustomerOnly(db, 1); err != nil {

		t.Errorf("got %v: %v", err, c)
	}

}
