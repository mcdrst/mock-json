package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
)

//helpers
func readDBJSONFile(fileName string) {
	jsonFile, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(jsonFile, &db)
	if err != nil {
		log.Fatal(err)
	}
}

func saveJSONFile(fileName string, fileData interface{}) error {
	js, err := json.MarshalIndent(fileData, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	err = os.WriteFile(fileName, js, 0644)
	if err != nil {
		return err
	}
	return nil
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

//sort collection slice by string
func sortSliceString(dataSlice sliceStringInterface, keyOrder, sortDirection string) sliceStringInterface {
	index := make(PairList, len(dataSlice))
	i := 0
	for k, v := range dataSlice {
		index[i] = Pair{v[keyOrder].(string), k}
		i++
	}
	sort.Sort(index)

	orderedDataSlice := make(sliceStringInterface, len(dataSlice))
	if sortDirection == "asc" {
		for k, v := range index {
			orderedDataSlice[k] = dataSlice[v.Value]
		}
	} else {
		for k, v := range index {
			orderedDataSlice[len(dataSlice)-k-1] = dataSlice[v.Value]
		}
	}
	return orderedDataSlice
}

//sort collection slice by float - numbers from json are parsed to floats
func sortSliceFloat(dataSlice sliceStringInterface, keyOrder, sortDirection string) sliceStringInterface {
	index := make(PairFloatList, len(dataSlice))
	i := 0
	for k, v := range dataSlice {
		index[i] = PairFloat{v[keyOrder].(float64), k}
		i++
	}
	sort.Sort(index)

	orderedDataSlice := make(sliceStringInterface, len(dataSlice))
	if sortDirection == "asc" {
		for k, v := range index {
			orderedDataSlice[k] = dataSlice[v.Value]
		}
	} else {
		for k, v := range index {
			orderedDataSlice[len(dataSlice)-k-1] = dataSlice[v.Value]
		}
	}
	return orderedDataSlice
}

//Create headers and body response
func writeJSON(w http.ResponseWriter, status int, data mapStringInterface) error {
	js, err := json.Marshal(data)
	//js, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

//read int URL query string  (example => ID)
func readInt(qs url.Values, key string, defaultValue int) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil || i < 1 {
		return -1
	}
	return i
}

//read int URL query string  (example => name)
func readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}
