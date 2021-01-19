package main

import "C"
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func handle(w http.ResponseWriter, r *http.Request) {
	//b := []byte("Hello World!")
	//w.Write(b)
	//fmt.Println(r.Host, "in handle")
	if string(r.Method) == "GET" {
		//fmt.Println(string("GET - Hello"))
		io.WriteString(w, string("GET - Hello"))
	} else {
		//fmt.Println(string("POST - Hello"))
		io.WriteString(w, string("POST - Hello"))
	}
}

func PrinterHandler(w http.ResponseWriter, r *http.Request) {
	//b := []byte("Hello World!")
	//w.Write(b)
	//fmt.Println(r.Host, "in printhandle")
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println(string("Printing "+strings.Split(r.URL.Path, "/")[2]))
	io.WriteString(w, string("Printing "+strings.Split(r.URL.Path, "/")[2]))
}

func CustomerHandler(w http.ResponseWriter, r *http.Request) {
	type Customer struct {
		Name    string
		Age     int
		Address string
	}
	var customer Customer

	bodyData, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyData, &customer)
	//fmt.Println(customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Error occurred")
	}
	//fmt.Println(customer)
	if customer.Age < 18 {
		//fmt.Println(customer)
		io.WriteString(w, "You are not eligible to be our customer")
	} else {
		io.WriteString(w, string(bodyData))
	}

}

type Address struct {
	Id         int    `json:id`
	StreetName string `json:streetName`
	City       string `json:city`
	State      string `json:state`
	CusId      int    `json:cus_id`
}

type Customers struct {
	Id      int     `json:id`
	Name    string  `json:name`
	Dob     string  `json:dob`
	Address Address `json:address`
}

var db, dbErr = sql.Open("mysql", "root:1118209@/Customer_service")

func insertCust(db *sql.DB) {
	stmtIns, err := db.Prepare("INSERT INTO customers VALUES( ?, ?, ?)") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	for i := 6; i <= 10; i++ {
		_, err = stmtIns.Exec(i, "Rupesh Raj", "10-02-1997") // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}

func insertAddrs(db *sql.DB) {
	stmtIns, err := db.Prepare("INSERT INTO address VALUES( ?, ?, ?, ?, ?)") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close()

	for i := 11; i <= 15; i++ {
		rand.Seed(int64(i * i))
		num := rand.Intn(10)
		fmt.Println(num)
		_, err = stmtIns.Exec(i, "ram Nagar", "Ara", "Bihar", num) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}

func getCustomerOnly(db *sql.DB, id int) (Customers, error) {
	var c Customers
	rows, err := db.Query("SELECT id, name, dob FROM customers")

	//var data []interface{}
	//
	//data = append(data, id)
	//
	//rows, err := db.Query(query, data...)
	//if err != nil {
	//	return c, err
	//}
	if err != nil {
		return c, err
	}

	for rows.Next() {
		rows.Scan(&c.Id, &c.Name, &c.Dob)
		fmt.Println(c)
	}

	return c, nil
}

func GetCustomer(db *sql.DB, id string) []Customers {

	var customer []Customers
	query := "SELECT * FROM customers inner join address on address.cus_id=customers.id order by customers.id, address.id"
	var ids []interface{}

	if id != "0" {
		query = "SELECT * FROM customers inner join address on address.cus_id=customers.id and customers.id= ? order by customers.id, address.id"
		d, err := strconv.Atoi(id)
		//fmt.Println(err)
		if err != nil {

		}
		ids = append(ids, d)
	}

	rows, err := db.Query(query, ids...)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer rows.Close()

	for rows.Next() {
		var c Customers
		if err = rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId); err != nil {
			log.Fatal(err)
		}
		customer = append(customer, c)
	}
	return customer
}

func CreateCustomer(db *sql.DB, c Customers) Customers {
	query := `insert into cust values (?,?,?)`

	var custData []interface{}

	custData = append(custData, c.Id)
	custData = append(custData, c.Name)
	custData = append(custData, c.Dob)

	_, err := db.Query(query, custData...)

	query = `insert into addrs values (?,?,?,?,?)`

	var addrsData []interface{}

	addrsData = append(addrsData, c.Address.Id)
	addrsData = append(addrsData, c.Address.StreetName)
	addrsData = append(addrsData, c.Address.City)
	addrsData = append(addrsData, c.Address.State)
	addrsData = append(addrsData, c.Address.CusId)

	_, err = db.Query(query, addrsData...)

	if err != nil {
		log.Fatal(err)
	}

	var rows *sql.Rows
	rows, err = db.Query("SELECT * FROM cust inner join addrs on addrs.cus_id=cust.id and cust.id= ? order by cust.id asc, addrs.id asc", c.Id)

	if err != nil {
		log.Fatal(err)
	}
	var customer Customers
	for rows.Next() {
		if err = rows.Scan(&customer.Id, &customer.Name, &customer.Dob, &customer.Address.Id, &customer.Address.StreetName, &customer.Address.City, &customer.Address.State, &customer.Address.CusId); err != nil {

			log.Fatal(err)
		}
	}
	return customer

}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	param := strings.Split(r.URL.Path, "/")[1]
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		//log.Fatal(err)
		//fmt.Println(err)
	}
	if param == "" || paramInt == 0 {
		stmtOut, err := db.Query("SELECT * FROM customers inner join address on address.cus_id=customers.id order by customers.id, address.id")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtOut.Close()
		var customer []Customers
		for stmtOut.Next() {
			var c Customers
			if err := stmtOut.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId); err != nil {
				log.Fatal(err)
			}
			customer = append(customer, c)
		}

		json.NewEncoder(w).Encode(customer)
	} else {
		stmtOut, err := db.Query(`SELECT * FROM customers inner join address on address.cus_id=customers.id and customers.id= ? order by customers.id, address.id`, paramInt)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtOut.Close()
		var customer []Customers
		for stmtOut.Next() {
			var c Customers
			if err := stmtOut.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId); err != nil {
				log.Fatal(err)
			}
			customer = append(customer, c)
		}

		json.NewEncoder(w).Encode(customer)
	}

}

type cust struct {
	Name string `json:name`
	Age  int    `json:age`
	Addr string `json:addr`
}

func (c *cust) UnmarshalJSON(b []byte) error {
	type cust struct {
		Name string `json:name`
		Age  int    `json:age`
	}
	var cr cust
	err := json.Unmarshal(b, &cr)
	if err != nil {
		return err
	}
	c.Name = cr.Name
	c.Age = cr.Age

	return nil
}

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	data := []byte(`{"name":"Bittu Ray", "age":22, "addr":"Patna"}`)
	var c cust
	json.Unmarshal(data, &c)
	json.NewEncoder(w).Encode(c)
}

func main() {
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr.Error()) // proper error handling instead of panic in your app
	}
	//insertCust(db)
	//insertAddrs(db)
	//http.HandleFunc("/", GetHandler)
	http.HandleFunc("/print/", PrinterHandler)
	http.HandleFunc("/customer", CustomerHandler)
	http.HandleFunc("/json", JsonHandler)
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
