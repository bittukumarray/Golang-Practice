package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Address struct {
	Id         int    `json:id`
	StreetName string `json:streetName`
	City       string `json:city`
	State      string `json:state`
	CusId      int    `json:cus_id`
}

type Customer struct {
	Id      int     `json:id`
	Name    string  `json:name`
	Dob     string  `json:dob`
	Address Address `json:address`
}

var db, dbErr = sql.Open("mysql", "root:1118209@/Customer_service")

func GetCustomersData(db *sql.DB, name string) []Customer {
	query := "select * from cust inner join addrs on cust.id=addrs.cus_id order by cust.id, addrs.id"
	var data []interface{}
	//fmt.Println("name is ", name)
	if name != "" {
		query = "select * from cust inner join addrs on cust.id=addrs.cus_id where cust.name=? order by cust.id, addrs.id"
		data = append(data, name)
	}

	rows, err := db.Query(query, data...)

	if err != nil {
		return []Customer{}
	}

	var customer []Customer

	for rows.Next() {
		var c Customer
		rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId)
		customer = append(customer, c)
	}

	return customer
}

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	//fmt.Println(params)
	name, ok := params["name"]
	var C []Customer
	if ok && len(name) > 0 {
		C = GetCustomersData(db, params["name"][0])
	} else {
		C = GetCustomersData(db, "")
	}

	json.NewEncoder(w).Encode(C)

}

func GetCustomerData(db *sql.DB, id int) Customer {
	rows, err := db.Query("select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=? order by cust.id, addrs.id", id)

	if err != nil {
		return Customer{}
	}

	var c Customer

	for rows.Next() {
		rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId)
	}

	return c
}

func GetCustomerHandler(w http.ResponseWriter, r *http.Request) {
	//pathparams := mux.Vars(r)
	//
	//var c Customer
	//id, err := strconv.Atoi(pathparams["id"])

	pathParams, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	var c Customer
	if err != nil {
		log.Fatal(err)
	}
	c = GetCustomerData(db, pathParams)

	json.NewEncoder(w).Encode(c)
}

func InsertCustomerData(db *sql.DB, obj Customer) Customer {
	_, err := db.Query("insert into Atable (name) values (?)", "Ray")

	if err != nil {
		log.Fatal(err)
		return Customer{}
	}
	id, err := db.Query("SELECT LAST_INSERT_ID();")

	if err != nil {
		log.Fatal(err)
		return Customer{}
	}

	var ID int
	for id.Next() {
		id.Scan(&ID)

		fmt.Println(ID)
	}
	fmt.Println("new id is ", ID)
	return Customer{}
}

func DateSubstract(d1 string) int {
	d1_slice := strings.Split(d1, "-")

	newDate := d1_slice[2] + "-" + d1_slice[1] + "-" + d1_slice[0]
	myDate, err := time.Parse("2006-01-02", newDate)

	if err != nil {
		panic(err)
	}

	return int(time.Now().Unix() - myDate.Unix())
}

func PostCustomerHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	x := time.Now().Unix()
	fmt.Println(x)
	if err != nil {
		log.Fatal(err)
	}
	var cust Customer
	err = json.Unmarshal(body, &cust)

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("cust is ", cust)
	timestamp := DateSubstract(cust.Dob)

	if timestamp/(3600*24*12*30) < 18 {
		json.NewEncoder(w).Encode(Customer{})
	} else {
		InsertCustomerData(db, cust)
	}
}

func UpdateData(db *sql.DB, id int, c Customer) {

}

func PutCustomerHandler(w http.ResponseWriter, r *http.Request) {
	pathParams, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(Customer{})
	}
	var customer Customer
	bodyData, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bodyData, &customer)

	if err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(Customer{})
	}
	if customer.Id != 0 || customer.Dob != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id and DOB can't be updated"))
	} else {
		UpdateData(db, pathParams, customer)
	}

}

func DeleteData(db *sql.DB, id int) Customer {
	rows, err := db.Query("select * from cust inner join addrs on addrs.cus_id=cust.id and cust.id=? order by cust.id, addrs.id", id)
	if err != nil {
		log.Fatal(err)
		return Customer{}
	}
	var c Customer
	for rows.Next() {
		rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Address.Id, &c.Address.StreetName, &c.Address.City, &c.Address.State, &c.Address.CusId)
	}
	rows, err = db.Query("delete from cust where id=?", id)
	if err != nil {
		log.Fatal(err)
		return Customer{}
	}
	//fmt.Println("cust is ", c)
	return c
}

func DeleteCustomerHandler(w http.ResponseWriter, r *http.Request) {
	pathParams, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])

	if err != nil {
		log.Fatal(err)
		json.NewEncoder(w).Encode(Customer{})
	} else {
		c := DeleteData(db, pathParams)
		json.NewEncoder(w).Encode(c)
	}

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
	r := mux.NewRouter()
	r.HandleFunc("/customer", GetCustomersHandler).Methods("GET")
	r.HandleFunc("/customer/{id}", GetCustomerHandler).Methods("GET")
	r.HandleFunc("/customer/", PostCustomerHandler).Methods("POST")
	r.HandleFunc("/customer/{id}", PutCustomerHandler).Methods("PUT")
	r.HandleFunc("/customer/{id}", DeleteCustomerHandler).Methods("DELETE")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))

}
