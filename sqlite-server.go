package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	_ "github.com/mattn/go-sqlite3"
)

var (
	writeQueue chan string
)

func sqliteRead(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Println("GET called for /sqliteRead")
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println("Error in request body")
		}
		readQuery := string(requestDump)
		fmt.Println("sspInstaller: Request body is ", string(readQuery))
	case "POST":
		fmt.Println("POST not supported for this URL")
	}
}

func sqliteWrite(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println("Error in request body")
		}
		fmt.Println("POST called for /sqlitewrite", requestDump)
		body, err := ioutil.ReadAll(req.Body)
		writeQuery := string(body)
		fmt.Println("Write statement: ", string(writeQuery))
		writeQueue <- writeQuery
	case "GET":
		fmt.Println("GET not supported for this URL")
	}
}

///Basic Requirements
///1. Setup simple file-server
///2. Expose a URL with GET handling reads. This URL will handle concurrent reads cos of WAL
///3. Expose second URL for SQLite writes. We use Go routine to serialize Writes to the database
func main() {
	fmt.Println("Starting SQLite server")
	writeQueue = make(chan string, 100)
	go writeQueueHandler()
	http.HandleFunc("/sqlitewrite", sqliteWrite)
	http.HandleFunc("/sqliteread", sqliteRead)
	http.ListenAndServe(":8099", nil)
}

func writeQueueHandler() {
	for true {
		stm, okay := <-writeQueue
		fmt.Println("write statement", stm)
		if okay {
			fmt.Println("In write queue")
			db, err := sql.Open("sqlite3", "test.db")
			if err != nil {
				fmt.Println("Error in database opening")
			}
			statement, error := db.Prepare(stm)
			if error != nil {
				fmt.Println("Error in database statement")
				continue
			}
			statement.Exec()
			db.Close()
		}
	}
}
