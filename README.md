# CRUD server

# Testing

## Run the server
```
make run
```

## View empty store

```
make list_id_asc
curl -s GET http://localhost:8090/list?ordering=id:asc | jq '.[] | (.id + " " + (.labels|tostring))'
```

## Add first record + list

```
make create_1
curl -s POST http://localhost:8090/create -d '{ "id": "id_1", "labels": ["a","x"], "object": {"url": "https://google.com"}}'
record saved%

make list_id_asc
curl -s GET http://localhost:8090/list?ordering=id:asc | jq '.[] | (.id + " " + (.labels|tostring))'
"id_1 [\"a\",\"x\"]"
```

## Add second record + list

```
make create_2
curl -s POST http://localhost:8090/create -d '{ "id": "id_2", "labels": ["a","y"], "object": {}}'
record saved%

make list_id_desc
curl -s GET http://localhost:8090/list?ordering=id:desc | jq '.[] | (.id + " " + (.labels|tostring))'
"id_2 [\"a\",\"y\"]"
"id_1 [\"a\",\"x\"]"
```

## List with filter

```
make list_x
curl -s GET http://localhost:8090/list?filtering=x | jq '.[] | (.id + " " + (.labels|tostring))'
"id_1 [\"a\",\"x\"]"

make list_y
curl -s GET http://localhost:8090/list?filtering=y | jq '.[] | (.id + " " + (.labels|tostring))'
"id_2 [\"a\",\"y\"]"
```

## Delete
```
make delete_1
curl -I -X POST http://localhost:8090/delete/id_1
HTTP/1.1 200 OK
Date: Sat, 23 Mar 2024 04:23:28 GMT
Content-Length: 14
Content-Type: text/plain; charset=utf-8

make delete_1
curl -I -X POST http://localhost:8090/delete/id_1
HTTP/1.1 404 Not Found
Date: Sat, 23 Mar 2024 04:23:33 GMT
Content-Length: 0
```

# List remaining records
```
curl -s GET http://localhost:8090/list
[{"id":"id_2","createdAt":"2024-03-23T00:20:03.541105-04:00","labels":["a","y"],"object":{"tag":"","url":""}}]
```
