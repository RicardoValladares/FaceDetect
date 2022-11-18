package main

import (
	//"io/ioutil"
	"net/http"
	"fmt"
	"os"
	"io"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/identificar.html")})
	http.HandleFunc("/enrrolar.html", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/enrrolar.html")})
	http.HandleFunc("/opencv.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/opencv.js")})
	http.HandleFunc("/utils.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/utils.js")})
	http.HandleFunc("/haarcascade_frontalface_default.xml", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "web/haarcascade_frontalface_default.xml")})
	http.Handle("/enrrolados/", http.StripPrefix("/enrrolados", (http.FileServer(http.Dir("./enrrolados"))) ))
	http.HandleFunc("/enrrolar", func(w http.ResponseWriter, r *http.Request) { 
		identificador := r.FormValue("id")
		if identificador == "" || identificador == "null" {
			http.Error(w, "Error al obtener el identificador", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error al obtener el imagen", http.StatusBadRequest)
			file.Close()
			return
		}
		defer file.Close()
		dst, err := os.Create("./enrrolados/"+identificador)
		if err != nil {
			http.Error(w, "Error al enrrolar imagen", http.StatusBadRequest)
			file.Close()
			return
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Error al almacenar la imagen", http.StatusBadRequest)
			file.Close()
			return
		}
		file.Close()
		fmt.Fprintf(w, "Enrrolado Exitosamente '%s'\n",identificador)
	})
	http.HandleFunc("/identificar", func(w http.ResponseWriter, r *http.Request) { 
		fmt.Fprintf(w, "%s\n","Ricardo")
		//http.Error(w, "Error al identificar", http.StatusBadRequest)
	})
	http.ListenAndServe(":5000", nil)
}
