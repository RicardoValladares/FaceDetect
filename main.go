package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
//	"path/filepath"
)

func subidor(w http.ResponseWriter, r *http.Request) {

	// maixmo 1MB de subida
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	err := r.ParseMultipartForm(1024*1024)
	if err != nil {
		fmt.Fprintf(w, "El archivo es demasiado grande, excede 1MB o no contiene Multiples parte o se corrompio\n")
		fmt.Println("El archivo es demasiado grande, excede 1MB o no contiene Multiples parte o se corrompio")
		return
	}
	// inputfile de nombre Archivo
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Fprintf(w, "Error al obtener el archivo el servidor\n")
		fmt.Println("Error al obtener el archivo del cliente")
		return
	}
	defer file.Close()
	//fmt.Printf("Archivo: %+v\n", handler.Filename)
	//fmt.Printf("Tamaño: %+v\n", handler.Size)
	//fmt.Printf("MIME Encabezado: %+v\n", handler.Header)

	// generamos un nombre de archivo ramdom
	tempFile, err := ioutil.TempFile("upload", "*"+".png" )
	if err != nil {
		fmt.Fprintf(w, "La carpeta subidos no Existe o No se logro ramdomizar el nombre\n")
		fmt.Println("La carpeta subidos no Existe o No se logro ramdomizar el nombre")
		return
	}
	defer tempFile.Close()

	// obtenemos el archivo desde archivos temporales
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(w, "No se logro leer el archivo\n")
		fmt.Println("No se logro leer el archivo")
		return
	}

	// escribimos el archivo en el servidor
	tempFile.Write(fileBytes)
	fmt.Fprintf(w, "Subida Completada '%s'\n",tempFile.Name())

	//mostrar en consola server
	fmt.Println("Subida Completada (",tempFile.Name(),")")
}


func main() {
	fmt.Println("localhost:5000")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/index.html")})
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/index.html")})
	http.HandleFunc("/opencv.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/opencv.js")})
	http.HandleFunc("/utils.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/utils.js")})
	http.HandleFunc("/haarcascade_frontalface_default.xml", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/haarcascade_frontalface_default.xml")})
	http.HandleFunc("/upload", subidor)
	http.ListenAndServe(":5000", nil)
}
