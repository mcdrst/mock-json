
# mock-json

Restfull api for json files - Inspired by [server-json](https://github.com/typicode/json-server.git)

## Install

clone repository

```bash
$ make build

==> Building mock-json...
2021-09-10T02:33:47Z
v.1.0.1-0-gdfe427a
go build -ldflags='-s -X main.buildTime="2021-09-10T02:33:47Z"" -X main.version=v.1.0.1-0-gdfe427a' -o=./bin/mock-json ./src
==> Building main to linux...
GOOS=linux GOARCH=amd64 go build -ldflags='-s -X main.buildTime="2021-09-10T02:33:47Z"" -X main.version=v.1.0.1-0-gdfe427a' -o=./bin/linux_amd64/mock-json ./src
```

Binary wil be saved in ./bin/mock-json

## Command line paramters

```bash
$ mock-json --version
Version:        v.1.0.1-0-gdfe427a
Build time:     "2021-09-10T02:33:47Z"

$ mock-json --help
Usage of mock-json:
  -JSONFile string
        JSON File to serve (default "db.json")
  -displayflags
        Display flags
  -port int
        API server port (default 4000)
  -version
        Display version and exit
```

## Start server

```bash
$ ./bin/mock-json --JSONFile=dbtest.json
2021/09/10 09:47:37 Reading json file dbtest.json
2021/09/10 09:47:37 Starting http://localhost:4000
2021/09/10 09:47:37 Resources
2021/09/10 09:47:37     http://localhost:4000/product
2021/09/10 09:47:37     http://localhost:4000/category
```

json structure

```json
{
    "product": [
        {
            "id": 1,
            "name": "burguer",
            "price": 8.5,
            "categoryId": 1
        },
        {
            "id": 2,
            "name": "burguer x",
            "price": 10.5,
            "categoryId": 2
        }
    ],
    "category": [
        {
            "id": 1,
            "name": "burguer",
            "valid": true
        },
        {
            "id": 2,
            "name": "Water",
            "valid": false
        },
        {
            "id": 3,
            "name": "cheese",
            "valid": true
        }
    ]
}
```

"id" key is mandatory and must be an unique and positve integer (it is the primary index). API find the greater id and sum 1 to the next record that will be created.

"values" must be a number, string or a bool. Numbers are always float64

## Endpoints

mock-json works with collections, each collection generates a Resource with below list endpoints.

* **GET** - _/:nameCollection_ Get all collection records

* **GET** (id) - _/:nameCollection/:id_ id must be an positive integer

* **POST** - _/:nameColletcion_ Create a new colletiton record, if collection does not exist, a new colletion will be created with a new record

* **PATCH** (id) - _/:nameColletion/:id_ Update a record (there is no PUT, send all data to PATCH)

**DELETE** (id) - _/:nameColletion/:id_ Deleta a record

## URL query string

**GET** - _/nameColletion?param1=value&param2=value..._

**sort** _:/collectionName?sort=key_ (ASC) or :/collectionName?sort=-key (DESC)
key must be a string

**page** _:/collectionName?page=99_ Page number to be recovered

**limit** _:/collectionName?limit=9_ Records per page

**any key** _:/collectionName?key=value_ Return all records that have an equal pair (key, value). Value could be a number, string or bool, null values are unsuported yet, plesse always use zero values (0, "", or (true/false))

## Example

```text
http://localhost:4000/product?page=1&limit=2&categoryId=1&sort=-price
```

Response

```json
{
  "metadata": {
    "current_page": 1,
    "page_size": 2,
    "first_page": 1,
    "last_page": 1,
    "total_records": 1
  },
  "product": [
    {
      "categoryId": 1,
      "id": 1,
      "name": "burguer",
      "price": 8.5
    }
  ]
}
```

## limitations

* json null values
* sort, page and limit are reserved wordks, don't use then as keys in json file
* don't try sort param if json file does not have the same key in all records
* many more...
