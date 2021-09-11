package main

import (
	"net/http"

	"github.com/bmizerany/pat"
)

func routes() *pat.PatternServeMux {
	m := pat.New()
	http.Handle("/check", http.HandlerFunc(handlerLog(check)))

	m.Get("/:collectionName/:id", http.HandlerFunc(handlerLog(collectionDelGetPatchById)))
	m.Get("/:collectionName", http.HandlerFunc(handlerLog(getColecaoAll)))
	m.Del("/:collectionName/:id", http.HandlerFunc(handlerLog(collectionDelGetPatchById)))
	m.Patch("/:collectionName/:id", http.HandlerFunc(handlerLog(collectionDelGetPatchById)))
	m.Post("/:collectionName", http.HandlerFunc(handlerLog(createCollection)))

	return m
}
