
push:
	git status
	git add .
	git commit -m "$$(date)"
	git pull origin Golang 
	git push origin Golang

gorun:
	go run main.go

run:
	./main

compile:
	go build main.go

