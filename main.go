package main

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"io/ioutil"

	"log"
	"github.com/leandroveronezi/go-recognizer"
	"path/filepath"
)


const fotosDir = "enrrolados"
const dataDir = "models"

func addFile(rec *recognizer.Recognizer, Path, Id string) {
	err := rec.AddImageToDataset(Path, Id)
	if err != nil {
		fmt.Println(err)
		return
	}
}


var rec = recognizer.Recognizer{}

func main() {
		err := rec.Init(dataDir)
		if err != nil {
			fmt.Println(err)
			return
		}
		rec.Tolerance = 0.4
		rec.UseGray = true
		rec.UseCNN = false
		defer rec.Close()


		files, err := ioutil.ReadDir("./enrrolados/")
	    if err != nil {
	        log.Fatal(err)
	    }
	
	    for _, file := range files {
	    	addFile(&rec, filepath.Join(fotosDir, file.Name()), file.Name())
	        fmt.Println(file.Name())
	    }
		
		rec.SetSamples()

		
	


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

		
		//fmt.Fprintf(w, "%s\n","Ricardo")
		//http.Error(w, "Error al identificar", http.StatusBadRequest)

		file, handler, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error al obtener la imagen", http.StatusBadRequest)
			return
		}
		defer file.Close()
		tempFile, err := ioutil.TempFile("temp", "*"+filepath.Ext(handler.Filename))
		if err != nil {
			http.Error(w, "Error al nombrar la imagen", http.StatusBadRequest)
			return
		}
		defer tempFile.Close()
		
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Error al leer la imagen", http.StatusBadRequest)
			return
		}
		tempFile.Write(fileBytes)
			//fmt.Fprintf(w, "Subida Completada '%s'\n",)
		

		face, err := rec.ClassifyMultiples(tempFile.Name())
		//_, err = rec.Classify(tempFile.Name())
		if err != nil {
	        log.Fatal(err)
	        return
	    }

		if len(face) > 0 {
			fmt.Println(face[0].Data.Id)
			fmt.Fprintf(w, face[0].Data.Id)
		} else {
			fmt.Println("NO IDENTIFICADO")
			fmt.Fprintf(w, "NO IDENTIFICADO")
		}
		




		
	})
	http.ListenAndServe(":5000", nil)
}
