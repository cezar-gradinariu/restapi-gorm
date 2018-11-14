package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var db *gorm.DB
var err error

var (
	server   = "localhost"
	port     = 1433
	user     = "sa"
	password = "<YourStrong!Passw0rd>"
	database = "GORM_TEST_DB"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

func InitialMigration() {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to the database")
	}
	defer db.Close()

	db.AutoMigrate(&User{})
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to the database")
	}
	defer db.Close()

	var users []User
	//db.Offset(10).Limit(10).Find(&users) // TODO : IN mssql it doesn't work, there is a fix: https://github.com/jinzhu/gorm/issues/1205 but i don't understand it yet....
	db.Find(&users)

	json, err := json.MarshalIndent(users, "", "	")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to json pretty the response")
	}
	fmt.Fprintf(w, string(json))
	// json.NewEncoder(w).Encode(users)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to the database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	db.Create(&User{Name: name, Email: email})
	fmt.Fprintf(w, "New user successfully created")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to the database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	db.Where("name = ?", name).Delete(&User{})

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	connectionString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s",
		server, user, password, port, database)
	db, err := gorm.Open("mssql", connectionString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to the database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	email := vars["email"]

	var user User
	db.Where("name = ?", name).Find(&user)
	if &user != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("400 - Not Found"))
	} else {
		user.Email = email
		db.Save(&user)
	}
}

// https://github.com/jinzhu/gorm/issues/1205 but i don't understand it yet.... what is the (mssql) thing after func ???
// func (mssql) LimitAndOffsetSQL(limit, offset interface{}) (sql string) {
// 	if offset != nil {
// 		if parsedOffset, err := strconv.ParseInt(fmt.Sprint(offset), 0, 0); err == nil && parsedOffset > 0 {
// 			sql += fmt.Sprintf(" OFFSET %d ROWS", parsedOffset)
// 		}
// 	}
// 	if limit != nil {
// 		if parsedLimit, err := strconv.ParseInt(fmt.Sprint(limit), 0, 0); err == nil && parsedLimit > 0 {
// 			if sql == "" {
// 				// add default zero offset
// 				sql += " OFFSET 0 ROWS"
// 			}
// 			sql += fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", parsedLimit)
// 		}
// 	}
// 	return
// }
