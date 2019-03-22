package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	sqlInterface "dsmartsys/sqlinterface"
	ut "dsmartsys/util"

	mx "github.com/gorilla/mux"
)

type app struct {
	sqlDB     *sqlInterface.MySqldb
	router    *mx.Router
	dbName    string
	tableName string
}

type login struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Organization string `json:"organization"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Website      string `json:"website"`
}

type errorneous struct {
	Error interface{}
}

func main() {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in main function, Error Info: ", errD)
		}
	}()
	a := &app{}
	user := os.Args[1]
	pass := os.Args[2]
	host := os.Args[3]
	port := os.Args[4]
	DNS := user + ":" + pass + "@tcp(" + host + ":" + port + ")/"
	//DNS: user:passwd@tcp(127.0.0.1:3306)/
	a.sqlDB = a.initialiseDB(DNS)
	defer a.sqlDB.Close()

	a.router = mx.NewRouter()
	a.initialiseRouter()
	a.run("localhost:8585")
}

func (a *app) initialiseDB(DNS string) *sqlInterface.MySqldb {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in initialiseDB method, Error Info: ", errD)
		}
	}()
	a.dbName = "dsmartsys"
	createDB := "CREATE DATABASE IF NOT EXISTS " + a.dbName + ";"
	a.tableName = "login"
	table := "CREATE TABLE IF NOT EXISTS " + a.dbName + "." + a.tableName + " (id INT NOT NULL AUTO_INCREMENT, first_name VARCHAR(255), last_name VARCHAR(255), organization VARCHAR(255), phone_number VARCHAR(255), email VARCHAR(255), website VARCHAR(255), PRIMARY KEY (id));"
	sqlDB, err := sqlInterface.CreateDataBase(DNS, createDB, table)
	if err != nil {
		log.Fatal(err)
	}
	return sqlDB
}

func (a *app) run(addr string) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in run method, Error Info: ", errD)
		}
	}()
	log.Fatal(http.ListenAndServe(addr, a.router))
}

func (a *app) initialiseRouter() {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in initialiseRouter method, Error Info: ", errD)
		}
	}()
	//http://localhost:8585/login
	a.router.HandleFunc("/logins", a.getLogins).Methods("GET")
	a.router.HandleFunc("/login/{id:[0-9]+}", a.getLogin).Methods("GET")
	a.router.HandleFunc("/login", a.createLogin).Methods("POST")
	a.router.HandleFunc("/login/{id:[0-9]+}", a.updateLogin).Methods("PUT")
}

func (a *app) getLogin(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in getLogin method, Error Info: ", errD)
		}
	}()
	fmt.Println("Inside getLogin")
	params := mx.Vars(r)
	id := params["id"]
	strQuery := fmt.Sprintf("SELECT first_name, last_name, organization, phone_number, email, website from %s.%s WHERE id=%s;", a.dbName, a.tableName, id)
	rows, err := a.sqlDB.SelectQuery(strQuery)
	if err != nil {
		fmt.Println("db.Exec err:", err)
		respondWithJSON(w, http.StatusInternalServerError, errorneous{err.Error()})
		return
	}
	defer rows.Close()
	var ls []login
	for rows.Next() {
		var l login
		if err := rows.Scan(&l.FirstName, &l.LastName, &l.Organization, &l.PhoneNumber, &l.Email, &l.Website); err != nil {
			fmt.Println("rows.Scan ,Error Info : ", err)
			break
		}
		ls = append(ls, l)
	}
	if ls != nil {
		respondWithJSON(w, http.StatusOK, ls)
	} else {
		respondWithJSON(w, http.StatusBadRequest, errorneous{"Selected id doesn't exist in table"})
	}
}

func (a *app) getLogins(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in getLogins method, Error Info: ", errD)
		}
	}()
	fmt.Println("Inside getLogins")
	strQuery := fmt.Sprintf("SELECT first_name, last_name, organization, phone_number, email, website from %s.%s;", a.dbName, a.tableName)
	rows, err := a.sqlDB.SelectQuery(strQuery)
	if err != nil {
		fmt.Println("db.Exec err:", err)
		respondWithJSON(w, http.StatusInternalServerError, errorneous{err.Error()})
		return
	}
	defer rows.Close()
	var ls []login
	for rows.Next() {
		var l login
		if err := rows.Scan(&l.FirstName, &l.LastName, &l.Organization, &l.PhoneNumber, &l.Email, &l.Website); err != nil {
			fmt.Println("rows.Scan ,Error Info : ", err)
			break
		}
		ls = append(ls, l)
	}
	if ls != nil {
		respondWithJSON(w, http.StatusOK, ls)
	} else {
		respondWithJSON(w, http.StatusNotImplemented, errorneous{"No records exist in table"})
	}
}

func (a *app) updateLogin(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in updateLogin method, Error Info: ", errD)
		}
	}()
	fmt.Println("Inside updateLogin")
	params := mx.Vars(r)
	id := params["id"]
	strQuery := fmt.Sprintf("SELECT id from %s.%s WHERE id=%s;", a.dbName, a.tableName, id)
	rows, err := a.sqlDB.SelectQuery(strQuery)
	if err != nil {
		fmt.Println("db.Exec err:", err)
		respondWithJSON(w, http.StatusInternalServerError, errorneous{err.Error()})
		return
	}
	defer rows.Close()
	if !rows.Next() {
		respondWithJSON(w, http.StatusBadRequest, errorneous{"Selected id doesn't exist in table"})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll error:", err)
	}
	var l login
	if err = json.NewDecoder(strings.NewReader(string(body))).Decode(&l); err != nil {
		fmt.Println("json.NewDecoder error:", err)
		respondWithJSON(w, http.StatusBadRequest, errorneous{err.Error()})
		return
	}
	defer r.Body.Close()

	strQuery = fmt.Sprintf("UPDATE %s.%s SET first_name='%s', last_name='%s', organization='%s', phone_number='%s', email='%s', website='%s' WHERE id=%s;", a.dbName, a.tableName, l.FirstName, l.LastName, l.Organization, l.PhoneNumber, l.Email, l.Website, id)
	err = a.sqlDB.ExecuteQuery(strQuery)
	if err != nil {
		fmt.Println("db.Exec err:", err)
		respondWithJSON(w, http.StatusInternalServerError, errorneous{err.Error()})
		return
	}
	respondWithJSON(w, http.StatusOK, l)
}

func (a *app) createLogin(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in createLogin method, Error Info: ", errD)
		}
	}()
	fmt.Println("Inside createLogin")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll error:", err)
	}
	var l login
	if err = json.NewDecoder(strings.NewReader(string(body))).Decode(&l); err != nil {
		fmt.Println("json.NewDecoder error:", err)
		respondWithJSON(w, http.StatusBadRequest, errorneous{err.Error()})
		return
	}
	defer r.Body.Close()

	strQuery := fmt.Sprintf("INSERT INTO %s.%s (first_name, last_name, organization, phone_number, email, website) VALUES('%s', '%s', '%s', '%s', '%s', '%s');", a.dbName, a.tableName, l.FirstName, l.LastName, l.Organization, l.PhoneNumber, l.Email, l.Website)
	err = a.sqlDB.ExecuteQuery(strQuery)
	if err != nil {
		fmt.Println("db.Exec err:", err)
		respondWithJSON(w, http.StatusInternalServerError, errorneous{err.Error()})
		return
	}
	respondWithJSON(w, http.StatusCreated, l)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	defer func() {
		if errD := recover(); errD != nil {
			fmt.Println("Exception occurred at ", ut.RecoverExceptionDetails(ut.FunctionName()), " and recovered in respondWithJSON function, Error Info: ", errD)
		}
	}()
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
