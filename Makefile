
push:
	git status
	git add .
	git commit -m "$$(date)"
	git pull origin Python 
	git push origin Python

run:
	python main.py
