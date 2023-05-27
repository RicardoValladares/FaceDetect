
push:
	git status
	git add .
	git commit -m "$$(date)"
	git pull origin Python 
	git push origin Python

gorun:
	go run main.go

run:
	./main

compile:
	go build main.go

