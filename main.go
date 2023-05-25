package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"github.com/leandroveronezi/go-recognizer"
)

const (
	fotosDir = "enrrolados"
	dataDir = "modelos"
)

var rec = recognizer.Recognizer{}

func addFile(rec *recognizer.Recognizer, Path, Id string) error {
	err := rec.AddImageToDataset(Path, Id)
	return err
}

func main() {
	err := rec.Init(dataDir)
	if err != nil {
		log.Println("Error al inicializar biometria con los templates")
		return
	}
	rec.Tolerance = 0.4
	rec.UseGray = true
	rec.UseCNN = false
	defer rec.Close()
	files, err := ioutil.ReadDir("./enrrolados/")
	if err != nil {
		log.Println("Error al inicializar biometria con los rostros enrrolados")
		return
	}
	fmt.Println("Cargando rostros enrrolados...")
	for _, file := range files {
		if addFile(&rec, filepath.Join(fotosDir, file.Name()), file.Name()) != nil {
			log.Println("Imagen no valida:", file.Name())
			return
		} else {
			log.Println("Imagen cargada:", file.Name())
		}
	}
	rec.SetSamples()
	fmt.Println("http://localhost:5001")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "enrrolar.html") })
	http.HandleFunc("/identificar.html", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "identificar.html") })
	http.HandleFunc("/opencv.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "opencv.js") })
	http.HandleFunc("/utils.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "utils.js") })
	http.HandleFunc("/haarcascade_frontalface_default.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "haarcascade_frontalface_default.xml")
	})
	http.Handle("/enrrolados/", http.StripPrefix("/enrrolados", (http.FileServer(http.Dir("./enrrolados")))))
	http.HandleFunc("/enrrolar", enrrolar)
	http.HandleFunc("/identificar", identificar)
	err = http.ListenAndServe(":5001", nil)
	if err != nil {
		log.Fatal("Error de ListenAndServe, Necesites derechos de super usuario y que el puerto 5001 este desocupado sin bloqueos del firewall")
	}
}


func identificar(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error al obtener la imagen", http.StatusBadRequest)
		log.Println("Error al obtener la imagen")
		file.Close()
		return
	}
	defer file.Close()
	tempFile, err := ioutil.TempFile("temp", "*"+filepath.Ext(handler.Filename))
	if err != nil {
		http.Error(w, "Error al nombrar la imagen", http.StatusBadRequest)
		log.Println("Error al nombrar la imagen")
		file.Close()
		tempFile.Close()
		return
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error al leer la imagen", http.StatusBadRequest)
		log.Println("Error al leer la imagen")
		file.Close()
		tempFile.Close()
		return
	}
	tempFile.Write(fileBytes)
	filetype := handler.Header.Get("content-type")
	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/gif" && filetype != "image/png" {
		http.Error(w, "Error de formato en el archivo", http.StatusBadRequest)
		log.Println("Error de formato en el archivo")
		file.Close()
		tempFile.Close()
		return
	}
	face, err := rec.ClassifyMultiples(tempFile.Name())
	if err != nil {
		http.Error(w, "Error al comparar rostro", http.StatusBadRequest)
		log.Println("Error al comparar rostro")
		file.Close()
		tempFile.Close()
		return
	}

	if len(face) > 0 {
		fmt.Fprintf(w, face[0].Data.Id)
		log.Println(face[0].Data.Id)
	} else {
		fmt.Fprintf(w, "NO IDENTIFICADO")
		log.Println("NO IDENTIFICADO")
	}
	file.Close()
	tempFile.Close()
}


func enrrolar(w http.ResponseWriter, r *http.Request) {
	identificador := r.FormValue("id")
	enrrolados, err := os.Stat("./enrrolados/" + identificador)
	if (err == nil) && !enrrolados.IsDir() {
		fmt.Fprintf(w, "Ya existe el identificador: '%s'\n", identificador)
		log.Println("Ya existe el identificador:", identificador)
		return
	}
	if identificador == "" || identificador == "null" {
		http.Error(w, "Error al obtener el identificador", http.StatusBadRequest)
		log.Println("Error al obtener el identificador")
		return
	}
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error al obtener el imagen", http.StatusBadRequest)
		log.Println("Error al obtener el imagen")
		file.Close()
		return
	}
	defer file.Close()
	filetype := handler.Header.Get("content-type")
	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/gif" && filetype != "image/png" {
		http.Error(w, "Error de formato en el archivo", http.StatusBadRequest)
		log.Println("Error de formato en el archivo")
		file.Close()
		return
	}
	dst, err := os.Create("./enrrolados/" + identificador)
	if err != nil {
		http.Error(w, "Error al enrrolar imagen", http.StatusBadRequest)
		log.Println("Error al enrrolar imagen")
		file.Close()
		dst.Close()
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error al almacenar la imagen", http.StatusBadRequest)
		log.Println("Error al almacenar la imagen")
		file.Close()
		dst.Close()
		return
	}
	file.Close()
	dst.Close()
	addFile(&rec, filepath.Join(fotosDir, identificador), identificador)
	rec.SetSamples()
	fmt.Fprintf(w, "Enrrolado Exitosamente '%s'\n", identificador)
	log.Println("Enrrolado Exitosamente ", identificador)
}
