package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	buildTime string
	version   string
)

type config struct {
	port     int
	JSONFile string
}

type mapStringInterface map[string]interface{}
type sliceStringInterface []map[string]interface{}

var db mapStringInterface

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Key < p[j].Key }

type PairFloat struct {
	Key   float64
	Value int
}
type PairFloatList []PairFloat

func (p PairFloatList) Len() int           { return len(p) }
func (p PairFloatList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairFloatList) Less(i, j int) bool { return p[i].Key < p[j].Key }

var cfg config

func main() {

	cfg.parseConfig()

	//load data from json
	log.Println("Reading json file", cfg.JSONFile)
	readDBJSONFile(cfg.JSONFile)
	endpoint := fmt.Sprintf("http://localhost:%v", cfg.port)
	log.Println(fmt.Sprintf("Starting %s", endpoint))
	log.Printf("Resources\n")
	for k := range db {
		log.Printf("\t%s/%s\n", endpoint, k)
	}
	m := routes()

	// Register this pat with the default serve mux so that other packages
	// may also be exported. (i.e. /debug/pprof/*)
	http.Handle("/", m)
	err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
