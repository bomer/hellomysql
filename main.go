package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"time"
)

var (
	data []SquareNumber
)

type SquareNumber struct {
	Number int
	Square int
}

func read(db *sql.DB, num int) {
	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()
	var squareNum int
	// db.Begin().
	// var rows sql.Rows
	// rows, err = stmtOut.Query(num).Scan(&squareNum) // WHERE number = 13
	err = stmtOut.QueryRow(num).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of %d is: %d \n", num, squareNum)

}

func readAllRows(db *sql.DB) {
	// Prepare statement for reading data
	start := time.Now()

	rows, err := db.Query("SELECT number,squareNumber FROM squarenum")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()

	// Create new temp Dataset
	tempDataSet := []SquareNumber{}

	for rows.Next() {
		var num int
		var square int
		if err := rows.Scan(&num, &square); err != nil {
			panic(err.Error())
		}
		newData := SquareNumber{Number: num, Square: square}
		tempDataSet = append(tempDataSet, newData)
		// fmt.Printf("%d square is %d\n", newData.number, newData.square)

	}
	if err := rows.Err(); err != nil {
		panic(err.Error())
	}
	data = tempDataSet
	tempDataSet = nil

	fmt.Printf("Sucessfully loaded %d rows in ", len(data))
	elapsed := time.Since(start)
	fmt.Println(elapsed)

}

func wipeTable(db *sql.DB) {

	stmtOut, err := db.Prepare("truncate squarenum")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	_, err = stmtOut.Exec()

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
func insert(db *sql.DB, limit int) {
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squarenum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Insert square numbers for 0-24 in the database
	for i := 1; i <= limit; i++ {

		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there. Hit '/data' on this server to get the latest results in JSON.")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(data)
	// js, err := json.Marshal(data)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// w.Header().Set("Content-Type", "application/json")
	// w.Write(js)

}

func pollDatabase(db *sql.DB) {
	for {
		time.Sleep(2 * time.Second)
		go readAllRows(db)
		fmt.Printf("There is %d rows", len(data))
	}
}
func main() {
	//Check for OS env variable. We use this on heorku, so if it's present use that else use local
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		fmt.Println("No db env variable found, using local")
		dburl = "root:password@tcp(127.0.0.1:3306)/testgo"
		// dburl = "bcbb9db7811db6:1e501e1c@tcp(us-cdbr-iron-east-02.cleardb.net:3306)/heroku_2c0d1e682389720" // testing heroku db
	}
	db, err := sql.Open("mysql", dburl)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	wipeTable(db)
	insert(db, 5000)
	// Query the square-number of some test numbers
	read(db, 13)
	read(db, 55)
	read(db, 155)
	read(db, 1555)

	// Query upfront so we dont need a null record
	readAllRows(db)

	http.HandleFunc("/", handler)
	http.HandleFunc("/data", dataHandler)

	go pollDatabase(db)

	http.ListenAndServe(":8008", nil)

}
