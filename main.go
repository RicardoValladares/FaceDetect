package main

import (
	"io/ioutil"
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/index.html")})
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/index.html")})
	http.HandleFunc("/opencv.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/opencv.js")})
	http.HandleFunc("/utils.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/utils.js")})
	http.HandleFunc("/haarcascade_frontalface_default.xml", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/haarcascade_frontalface_default.xml")})
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) { 
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error al obtener el archivo el servidor", http.StatusBadRequest)
			file.Close()
			return
		}
		defer file.Close()
		tempFile, err := ioutil.TempFile("upload", "*"+".png" )
		if err != nil {
			http.Error(w, "La carpeta subidos no Existe o No se logro ramdomizar el nombre", http.StatusBadRequest)
			file.Close()
			return
		}
		defer tempFile.Close()
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "No se logro leer el archivo", http.StatusBadRequest)
			file.Close()
			tempFile.Close()
			return
		}
		tempFile.Write(fileBytes)
		file.Close()
		tempFile.Close()
		fmt.Fprintf(w, "Subida Completada '%s'\n",tempFile.Name())
	})
	http.ListenAndServe(":5000", nil)
}

