package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func readBodyJSON(w http.ResponseWriter, r *http.Request, dst *mapStringInterface) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&dst)
	return err
}

//handler POST
func createCollection(w http.ResponseWriter, r *http.Request) {
	collectionName := r.URL.Query().Get(":collectionName")

	//read JSON from request
	var dst mapStringInterface
	err := readBodyJSON(w, r, &dst)
	if err != nil {
		parseJSONError(w, r)
		return
	}
	//create new record in var db
	data := mapStringInterface{}
	for k, v := range dst {
		data[k] = v
	}
	//calculate greater id
	dataSlice := getCollection(collectionName)
	id := 0.
	for _, item := range dataSlice {
		valueId := item["id"].(float64)
		if valueId > id {
			id = valueId
		}
	}

	//write new json file and read new db
	id++
	data["id"] = id
	dataSlice = append(dataSlice, data)
	db[collectionName] = dataSlice
	fileName := cfg.JSONFile
	//tests
	//fileName := "./datafortest/results.json"
	saveJSONFile(fileName, db)
	readDBJSONFile(fileName)

	data = mapStringInterface{collectionName: data}
	err = writeJSON(w, http.StatusOK, data)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}
}

//Handler GET PATCH DELETE
func collectionDelGetPatchById(w http.ResponseWriter, r *http.Request) {
	//find record by id and store in var data
	collectionName := r.URL.Query().Get(":collectionName")
	id := r.URL.Query().Get(":id")
	idFloat, err := strconv.ParseFloat(id, 64)
	if err != nil {
		parseFloatError(w, r)
		return
	}
	dataSlice := getCollection(collectionName)
	if dataSlice == nil {
		notFoundResponse(w, r)
		return
	}
	index := findById(dataSlice, idFloat)
	if index == -1 {
		notFoundResponse(w, r)
		return
	}
	data := mapStringInterface{collectionName: dataSlice[index]}

	fileName := cfg.JSONFile
	//tests
	//fileName := "./datafortest/results.json"

	//delete
	if r.Method == http.MethodDelete {
		dataSlice[index] = dataSlice[len(dataSlice)-1]
		dataSlice = dataSlice[:len(dataSlice)-1]
		db[collectionName] = dataSlice
		saveJSONFile(fileName, db)
		readDBJSONFile(fileName)
	}

	//patch
	if r.Method == http.MethodPatch {

		var dst mapStringInterface
		err := readBodyJSON(w, r, &dst)
		if err != nil {
			parseJSONError(w, r)
			return
		}
		for k, v := range dst {
			dataSlice[index][k] = v
		}
		db[collectionName] = dataSlice
		saveJSONFile(fileName, db)
		readDBJSONFile(fileName)
	}

	err = writeJSON(w, http.StatusOK, data)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}

}

//Handler GET (many)
func getColecaoAll(w http.ResponseWriter, r *http.Request) {

	var err error

	qs := r.URL.Query()

	fieldSort := readString(qs, "sort", "id")
	sortDirection := "asc"
	if strings.HasPrefix(fieldSort, "-") {
		fieldSort = fieldSort[1:]
		sortDirection = "desc"
	}
	page := readInt(qs, "page", 1)
	if page == -1 {
		failedValidationResponse(w, r, mapStringInterface{"page": "must be integer greater than zero"})
		return
	}
	limit := readInt(qs, "limit", 5)
	if limit == -1 {
		failedValidationResponse(w, r, mapStringInterface{"limit": "must be integer greater than zero"})
		return
	}

	nomeColecao := r.URL.Query().Get(":collectionName")

	dataSlice := getCollection(nomeColecao)
	for key, value := range qs {
		if !(key == "limit" || key == "sort" || key == "page" || strings.HasPrefix(key, ":")) {
			dataSlice, err = getFilterData(dataSlice, key, value[0])
			if err != nil {
				parseFloatError(w, r)
				return
			}
		}
	}
	if dataSlice == nil {
		notFoundResponse(w, r)
		return
	}

	var sortedDataSlice sliceStringInterface
	switch t := dataSlice[0][fieldSort].(type) {
	case float64:
		sortedDataSlice = sortSliceFloat(dataSlice, fieldSort, sortDirection)
	case string:
		sortedDataSlice = sortSliceString(dataSlice, fieldSort, sortDirection)
	default:
		_ = t
		sortedDataSlice = dataSlice
	}

	metadata := calculateMetadata(len(sortedDataSlice), page, limit)
	if page > metadata.LastPage {
		page = metadata.LastPage
		metadata.CurrentPage = page
	}

	finish := page * limit
	start := finish - limit
	if finish > len(sortedDataSlice) {
		finish = len(sortedDataSlice)
	}

	sortedDataSlice = sortedDataSlice[start:finish]
	data := mapStringInterface{"metadata": metadata, nomeColecao: sortedDataSlice}
	err = writeJSON(w, http.StatusOK, data)
	if err != nil {
		serverErrorResponse(w, r, err)
		return
	}
}

//Handler check (GET)
func check(w http.ResponseWriter, r *http.Request) {
	data := mapStringInterface{"status": "available"}
	err := writeJSON(w, http.StatusOK, data)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}

//Helper - get collection from var db
func getCollection(collectionName string) sliceStringInterface {
	var dataSlice sliceStringInterface

	colecao := db[collectionName]
	if colecao == nil {
		return nil
	}
	for _, v := range colecao.([]interface{}) {
		dataSlice = append(dataSlice, v.(map[string]interface{}))
	}
	return dataSlice
}

//Helper - find record by id and return index of collection slice
func findById(dataSlice sliceStringInterface, id float64) int {
	index := -1
	for k, item := range dataSlice {
		valueId := item["id"]
		if valueId.(float64) == id {
			return k
		}
	}
	return index
}

//Helper - query data form collection
func getFilterData(dataSlice sliceStringInterface, key, value string) (sliceStringInterface, error) {
	var data sliceStringInterface
	for _, item := range dataSlice {
		itemValue := item[key]
		if itemValue != nil {
			switch t := itemValue.(type) {
			case float64:
				valorFloat, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, err
				}
				if valorFloat == t {
					data = append(data, item)
				}
			case string:
				if strings.EqualFold(value, t) {
					data = append(data, item)
				}
			case bool:
				if value == strconv.FormatBool(t) {
					data = append(data, item)
				}
			}
		}
	}
	return data, nil
}
