package main

import "fmt"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

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
func insert(db *sql.DB) {
	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO squarenum VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Insert square numbers for 0-24 in the database
	for i := 1; i <= 1000; i++ {

		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
}

func main() {
	db, err := sql.Open("mysql", "root:password@/testgo")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	wipeTable(db)

	// Query the square-number of 13 & 55
	read(db, 13)
	read(db, 55)
}
