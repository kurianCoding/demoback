package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
)

const DB_USERNAME = "admin"
const DB_PWD = "project2846"
const AWS_RDS_URL = "mysql.ce66hgxuxi0k.ap-south-1.rds.amazonaws.com"
const AWS_PORT = "3306"
const DB_NAME = "project"

func connect() (*sql.DB, error) {
	//connect to mysql database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DB_USERNAME, DB_PWD, AWS_RDS_URL, AWS_PORT, DB_NAME))
	return db, err
}
func getRows(n string, conn *sql.DB) (error, map[string]interface{}) {
	result, err := conn.Query("select * from webcam  order by timestamp desc limit ?", n)
	if err != nil {
		return err, nil
	}
	var timearr []string
	var countarr []int32
	for result.Next() {
		var time string
		var count int32
		err = result.Scan(&time, &count)
		if err != nil {
			panic(err)
		}
		timearr = append(timearr, time)
		countarr = append(countarr, count)
	}
	res := map[string]interface{}{"time": timearr, "count": countarr}
	return nil, res
}

func main() {
	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}
	sqlconn, err := connect() // returns database connection
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// call db function
			(*&w).Header().Set("Access-Control-Allow-Origin", "*")
			rval := r.URL.Query()["rangeVal"][0]
			err, res := getRows(rval, sqlconn)
			if err != nil {
				fmt.Println(err)
			}
			// return JSONified response
			err = json.NewEncoder(w).Encode(res)
			return
		}
	})

	log.Printf("Starting server on port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
