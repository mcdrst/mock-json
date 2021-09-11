package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var data = sliceStringInterface{
	{"name": "peter", "age": 58., "id": 1., "valid": false},
	{"name": "jane", "age": 43., "id": 2., "valid": false},
	{"name": "mary", "age": 43., "id": 3., "valid": true},
}

//pat routes
var m = routes()

//copia json do diretorio tesdata
func setupDB() {
	cfg.JSONFile = "dbtest.json"
	readDBJSONFile("./testdata/dbtest.json")
	saveJSONFile(cfg.JSONFile, db)
}

func Test_getFilterData(t *testing.T) {

	type args struct {
		dataSlice sliceStringInterface
		key       string
		value     string
	}
	tests := []struct {
		name    string
		args    args
		want    sliceStringInterface
		wantErr bool
	}{
		{
			name: "find by age, return two records",
			args: args{
				dataSlice: data,
				key:       "age",
				value:     "43",
			},
			want: sliceStringInterface{
				{"name": "jane", "age": 43., "id": 2., "valid": false},
				{"name": "mary", "age": 43., "id": 3., "valid": true},
			},
			wantErr: false,
		},
		{
			name: "find bey name, return one record",
			args: args{
				dataSlice: data,
				key:       "name",
				value:     "jane",
			},
			want: sliceStringInterface{
				{"name": "jane", "age": 43., "id": 2., "valid": false},
			},
			wantErr: false,
		},
		{
			name: "find by valid, return one reocord",
			args: args{
				dataSlice: data,
				key:       "valid",
				value:     "true",
			},
			want: sliceStringInterface{
				{"name": "mary", "age": 43., "id": 3., "valid": true},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFilterData(tt.args.dataSlice, tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFilterData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilterData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findById(t *testing.T) {

	type args struct {
		dataSlice sliceStringInterface
		id        float64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "nonexistent ID",
			args: args{
				dataSlice: data,
				id:        1.,
			},
			want: 0,
		},
		{
			name: "existent ID",
			args: args{
				dataSlice: data,
				id:        22.,
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findById(tt.args.dataSlice, tt.args.id); got != tt.want {
				t.Errorf("findById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getColletcion(t *testing.T) {
	setupDB()

	type args struct {
		nameCollection string
	}
	tests := []struct {
		name string
		args args
		want sliceStringInterface
	}{
		{
			name: "find a category",
			args: args{nameCollection: "category"},
			want: sliceStringInterface{
				{
					"id":    1.,
					"name":  "burguer",
					"valid": true,
				},
				{
					"id":    2.,
					"name":  "Water",
					"valid": false,
				},
				{
					"id":    3.,
					"name":  "cheese",
					"valid": true,
				},
			},
		},
		{
			name: "nonexistent category",
			args: args{nameCollection: "nonexistent"},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCollection(tt.args.nameCollection); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_check(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/check", nil)
	t.Run("check test", func(t *testing.T) {
		check(w, r)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
		}
		if !strings.Contains(string(data), `"status":"available"`) {
			t.Errorf(`expected {"status":"avalilable"} got %v`, string(data))
		}
	})
}

func Test_getCollectionAll(t *testing.T) {
	setupDB()
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		sequence int
		name     string
		args     args
	}{
		{
			sequence: 1,
			name:     "01 - all",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category", nil),
			},
		},
		{
			sequence: 2,
			name:     "02 - nothing",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/inexistente", nil),
			},
		},
		{
			sequence: 3,
			name:     "03 - sort by -id and valid=true",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category?sort=-id&valid=true", nil),
			},
		},
		{
			sequence: 4,
			name:     "04 - page=2 limit=2",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category?page=2&limit=2", nil),
			},
		},
		{
			sequence: 5,
			name:     "05 - products with price preço=10.5",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product?price=10.5", nil),
			},
		},
		{
			sequence: 6,
			name:     "06 - product with name=burguer x",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product?name=burguer%20x", nil),
			},
		},
		{
			sequence: 7,
			name:     "07 - product by name descending",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product?sort=-name", nil),
			},
		},
		{
			sequence: 8,
			name:     "08 - product by name ascending",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product?sort=name", nil),
			},
		},
		{
			sequence: 9,
			name:     "09 - page greater than max",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product?page=10", nil),
			},
		},
		{
			sequence: 10,
			name:     "10 - invalid page",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product?page=abc", nil),
			},
		},
		{
			sequence: 11,
			name:     "11 - invalid limit",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product?limit=abc", nil),
			},
		},
		{
			sequence: 12,
			name:     "12 - sort by bool - no effect",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category?sort=valid", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.ServeHTTP(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			if tt.sequence == 1 {
				if !strings.Contains(string(data), `"id":1`) || !strings.Contains(string(data), `"id":2`) || !strings.Contains(string(data), `"id":3`) {
					t.Errorf(`expected ids 1, 2, and 3,  got %v`, string(data))
				}
			}
			if tt.sequence == 2 {
				if !strings.Contains(string(data), `{"error":"the request resource could not be found"}`) {
					t.Errorf(`expected {"error":"the request resource could not be found"}, got %v`, string(data))
				}
			}
			if tt.sequence == 3 {
				if !strings.Contains(string(data), `[{"id":3,"name":"cheese","valid":true},{"id":1,"name":"burguer","valid":true}]`) {
					t.Errorf(`expected [{"id":3,"name":"cheese","valid":true},{"id":1,"name":"burguer","valid":true}] , got %v`, string(data))
				}
			}
			if tt.sequence == 4 {
				if !strings.Contains(string(data), `[{"id":3,"name":"cheese","valid":true}]`) {
					t.Errorf(`expected [{"id":3,"name":"cheese","valid":true}] , got %v`, string(data))
				}
			}
			if tt.sequence == 5 {
				if !strings.Contains(string(data), `"categoryId":2,"id":2,"name":"burguer x","price":10.5`) {
					t.Errorf(`expected "categoryId":2,"id":2,"name":"burguer x","price":10.5, got %v`, string(data))
				}
			}
			if tt.sequence == 6 {
				if !strings.Contains(string(data), `categoryId":2,"id":2,"name":"burguer x","price":10.5`) {
					t.Errorf(`expected categoryId":2,"id":2,"name":"burguer x","price":10.5, got %v`, string(data))
				}
			}
			if tt.sequence == 7 {
				if !strings.Contains(string(data), `"categoryId":2,"id":2,"name":"burguer x","price":10.5},{"categoryId":1,"id":1,"name":"burguer","price":8.5`) {
					t.Errorf(`expected ywo records by name descending, got %v`, string(data))
				}
			}
			if tt.sequence == 8 {
				if !strings.Contains(string(data), `"categoryId":1,"id":1,"name":"burguer","price":8.5},{"categoryId":2,"id":2,"name":"burguer x","price":10.5}`) {
					t.Errorf(`expected ywo records by name ascending, got %v`, string(data))
				}
			}
			if tt.sequence == 9 {
				if !strings.Contains(string(data), `"categoryId":1,"id":1,"name":"burguer","price":8.5},{"categoryId":2,"id":2,"name":"burguer x","price":10.5`) {
					t.Errorf(`expected two records by name ascending, got %v`, string(data))
				}
			}
			if tt.sequence == 10 {
				if !strings.Contains(string(data), `"error":{"page"`) {
					t.Errorf(`expected invalid page error, got %v`, string(data))
				}
			}
			if tt.sequence == 11 {
				if !strings.Contains(string(data), `"error":{"limit"`) {
					t.Errorf(`expected invalid limit error, got %v`, string(data))
				}
			}
			if tt.sequence == 12 {
				if !strings.Contains(string(data), `[{"id":1,"name":"burguer","valid":true},{"id":2,"name":"Water","valid":false},{"id":3,"name":"cheese","valid":true}]`) {
					t.Errorf(`[{"id":1,"name":"burguer","valid":true},{"id":2,"name":"Water","valid":false},{"id":3,"name":"cheese","valid":true}], got %v`, string(data))
				}
			}
		})
	}
}
func Test_collectionDelGetPatchById(t *testing.T) {
	setupDB()
	var jsonData = []byte(`{
		"name": "name 1",
		"price": 100.15
	}`)
	var jsonDataError = []byte(`{
		name: "name 1",
		"price": 100.15,
	}`)
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		sequence int
		name     string
		args     args
	}{
		{
			sequence: 1,
			name:     "01 - invelid",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category/abc", nil),
			},
		},
		{
			sequence: 2,
			name:     "02 - nonexistent id",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category/100", nil),
			},
		},
		{
			sequence: 3,
			name:     "03 - nonexistent collection",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/nonexistent/100", nil),
			},
		},
		{
			sequence: 4,
			name:     "04 - patch sem dados, nao faz nada",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category/1", nil),
			},
		},
		{
			sequence: 5,
			name:     "05 - patch product 2",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPatch, "/product/2", bytes.NewBuffer(jsonData)),
			},
		},
		{
			sequence: 6,
			name:     "06 - json with error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPatch, "/product/2", bytes.NewBuffer(jsonDataError)),
			},
		},
		{
			sequence: 7,
			name:     "07 - delete category 2",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodDelete, "/category/2", nil),
			},
		},
		{
			sequence: 8,
			name:     "08 - get category 2 with error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/category/2", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sequence == 5 || tt.sequence == 6 {
				tt.args.r.Header.Set("Content-Type", "application/json; charset=UTF-8")
			}
			m.ServeHTTP(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			if tt.sequence == 1 {
				if !strings.Contains(string(data), `"error":"numeric parameter conversion error"`) {
					t.Errorf(`expected erro ao converter paramtro para float,  got %v`, string(data))
				}
			}
			if tt.sequence == 2 {
				if !strings.Contains(string(data), `"error":"the request resource could not be found"`) {
					t.Errorf(`expected "erro":"the request resource could not be found",  got %v`, string(data))
				}
			}
			if tt.sequence == 3 {
				if !strings.Contains(string(data), `"error":"the request resource could not be found"`) {
					t.Errorf(`expected "erro":"the request resource could not be found"t,  got %v`, string(data))
				}
			}
			if tt.sequence == 4 {
				if !strings.Contains(string(data), `{"id":1,"name":"burguer","valid":true}`) {
					t.Errorf(`expected {"id":1,"name":"burguer","valid":true},  got %v`, string(data))
				}
			}
			if tt.sequence == 5 {
				if !strings.Contains(string(data), `"product":{"categoryId":2,"id":2,"name":"name 1","price":100.15`) {
					t.Errorf(`expected "product":{"categoryId":2,"id":2,"name":"name 1","price":100.15,  got %v`, string(data))
				}
			}
			if tt.sequence == 6 {
				if !strings.Contains(string(data), `"error":"request body data error"`) {
					t.Errorf(`expected "error":"request body data error",  got %v`, string(data))
				}
			}
			if tt.sequence == 7 {
				if !strings.Contains(string(data), `{"id":2,"name":"Water","valid":false}`) {
					t.Errorf(`expected {"id":2,"name":"Water","valid":false},  got %v`, string(data))
				}
			}
			if tt.sequence == 8 {
				if !strings.Contains(string(data), `"error":"the request resource could not be found"`) {
					t.Errorf(`expected "erro":"the request resource could not be found",  got %v`, string(data))
				}
			}
		})
	}
}
func Test_createCollection(t *testing.T) {
	setupDB()
	var jsonData = []byte(`{
		"name": "name 1",
		"price": 100.15
	}`)
	var jsonDataError = []byte(`{
		name: "name 1",
		"price": 100.15,
	}`)
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		sequence int
		name     string
		args     args
	}{
		{
			sequence: 1,
			name:     "post new product",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(jsonData)),
			},
		},
		{
			sequence: 2,
			name:     "json with error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(jsonDataError)),
			},
		},
		{
			sequence: 3,
			name:     "new collection",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodPost, "/newcollection", bytes.NewBuffer(jsonData)),
			},
		},
		{
			sequence: 4,
			name:     "get product 3",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/product/3", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sequence == 1 || tt.sequence == 2 || tt.sequence == 3 {
				tt.args.r.Header.Set("Content-Type", "application/json; charset=UTF-8")
			}
			m.ServeHTTP(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil, got %v", err)
			}
			if tt.sequence == 1 {
				if !strings.Contains(string(data), `"product":{"id":3,"name":"name 1","price":100.15}`) {
					t.Errorf(`expected "product":{"id":3,"name":"name 1","price":100.15}, got %v`, string(data))
				}
			}
			if tt.sequence == 2 {
				if !strings.Contains(string(data), `"error":"request body data error"`) {
					t.Errorf(`expected "error":"request body data error",  got %v`, string(data))
				}
			}
			if tt.sequence == 3 {
				if !strings.Contains(string(data), `{"newcollection":{"id":1,"name":"name 1","price":100.15}`) {
					t.Errorf(`expected {"newcollection":{"id":1,"name":"name 1","price":100.15},  got %v`, string(data))
				}
			}
			if tt.sequence == 4 {
				if !strings.Contains(string(data), `{"product":{"id":3,"name":"name 1","price":100.15}`) {
					t.Errorf(`expected {"product":{"id":3,"name":"name 1","price":100.15},  got %v`, string(data))
				}
			}
		})
	}
}
