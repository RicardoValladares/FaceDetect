
push:
	git status
	git add .
	git commit -m "$$(date)"
	git pull origin facedetection 
	git push origin facedetection

gorun:
	go run main.go

run:
	./main

compile:
	go build main.go

