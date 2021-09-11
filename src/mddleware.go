package main

import (
	"log"
	"net/http"
)

func handlerLog(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s", r.Method, r.URL.Path)
		//start := time.Now()
		next.ServeHTTP(w, r)
		// elapsed := time.Since(start)
		// log.Printf("elapsed time %v", elapsed)
	})
}
