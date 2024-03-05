package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var PORT = 1954

func setupLogger(logFileName string) *os.File {
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// log to stdout and log file
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Println("Starting new instance of huntleyequine.com")
	return f
}

func loggingMiddleware(next http.Handler) http.Handler {
	hndlr := func(w http.ResponseWriter, r *http.Request) {
		remote_addr := r.Header.Get("X-Real-IP") // nginx header for real client IP
		log.Println(remote_addr + " requesting " + r.Method + " " + r.URL.Path)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(hndlr)
}

func main() {
	f := setupLogger("/home/frank/logs/huntleyequine.com.log")
	defer f.Close() // needs to be closed in this scope

	http.Handle("/", loggingMiddleware(http.FileServer(http.Dir("/home/frank/Code/util-servers/huntleyequine.com/public"))))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
}
