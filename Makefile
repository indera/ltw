.PHONY: create list delete

run: build
	./ltw

build:
	go build ./...

create_1:
	curl -s POST http://localhost:8090/create -d '{ "id": "id_1", "labels": ["a","x"], "object": {"url": "https://google.com"}}'
create_2:
	curl -s POST http://localhost:8090/create -d '{ "id": "id_2", "labels": ["a","y"], "object": {}}'


list:
	curl -s GET http://localhost:8090/list | jq '.[] | (.id + " " + (.createdAt|tostring))'
	curl -s GET http://localhost:8090/list?ordering=createdAt:desc | jq '.[] | (.id + " " + (.createdAt|tostring))'

list_a:
	curl -s GET http://localhost:8090/list?filtering=a | jq '.[] | (.id + " " + (.labels|tostring))'
list_x:
	curl -s GET http://localhost:8090/list?filtering=x | jq '.[] | (.id + " " + (.labels|tostring))'
list_y:
	curl -s GET http://localhost:8090/list?filtering=y | jq '.[] | (.id + " " + (.labels|tostring))'

list_id_asc:
	curl -s GET http://localhost:8090/list?ordering=id:asc | jq '.[] | (.id + " " + (.labels|tostring))'
list_id_desc:
	curl -s GET http://localhost:8090/list?ordering=id:desc | jq '.[] | (.id + " " + (.labels|tostring))'

delete_1:
	curl -I -X POST http://localhost:8090/delete/id_1
delete_2:
	curl -I -X POST http://localhost:8090/delete/id_2
