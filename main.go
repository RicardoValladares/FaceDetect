package main

import (
	"net/http"
	"fmt"
	"log"
	"os"
	"io"
	"io/ioutil"
	"path/filepath"
	"github.com/leandroveronezi/go-recognizer"
)


const fotosDir = "enrrolados"
const dataDir = "modelos"
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
			log.Println("Imagen no valida:",file.Name())
			return
		} else {
			log.Println("Imagen cargada:",file.Name())
		}
	}
	rec.SetSamples()
	fmt.Println("Cargando servidor...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
		fmt.Fprintf(w, `<html oncontextmenu='return false' onkeydown='return false'>
			<head>
				<meta name='viewport' content='width=device-width, initial-scale=1.0'> 
				<title>Enrrolar</title>
				<script async src='opencv.js' onload='openCvReady();'></script>
				<script src='utils.js'></script>
				<style> .labele { color: white; padding: 8px; font-family: Arial; background-color: #ff9800; } .labelt { color: white; padding: 8px; font-family: Arial; background-color: #04aa6d; } @media only screen and (max-width: 992px) { video.camara { height:640px; width:480px; display: block; margin-left: auto; margin-right: auto; } } @media only screen and (min-width: 993px) { video.camara { height:480px; width:640px; display: block; margin-left: auto; margin-right: auto;} } </style>
			</head>
			<body  bgcolor='#000' onload="setTimeout('temporizador()',1000)"> <br>
				<center> 
					<span id='Tiempo' class='labelt'>0</span> 
					<span id='Estado' class='labele'>Iniciando...</span> 
				</center><br><br>
				<center>
					<! canvas id='canvas_output' /><! /canvas> 
					<video id='cam_input' height='480' width='640' class='camara'></video> 
				</center>
				<script>
					/* variables globales para el funcionamiento */
					let tiempo = 0; //variable que sirve como contador de segundos
					let stop = 0; //variable que determina si detenemos el contador
					
					/* aperturamos webcam con opencv */
					function openCvReady() {
						cv['onRuntimeInitialized'] = () => {
							let video = document.getElementById('cam_input');
							/*video.style.display='none';*/ 
							navigator.mediaDevices.getUserMedia({ video: true, audio: false }).then(function (stream) { video.srcObject = stream; video.play(); }).catch(function (err) { console.log('Error: ' + err); });
							let src = new cv.Mat(video.height, video.width, cv.CV_8UC4);
							let gray = new cv.Mat();
							let cap = new cv.VideoCapture(cam_input);
							let faces = new cv.RectVector();
							let faceClassifier = new cv.CascadeClassifier();
							let utils = new Utils('errorMessage');
							let faceCascade = 'haarcascade_frontalface_default.xml';
							utils.createFileFromUrl(faceCascade, faceCascade, () => { faceClassifier.load(faceCascade); });
							const FPS = 40;
							function processVideo() {
								let begin = Date.now();
								cap.read(src);
								cv.cvtColor(src, gray, cv.COLOR_RGBA2GRAY, 0);
								let detectado=0;
								try {
									faceClassifier.detectMultiScale(gray, faces, 1.1, 3, 0);
									for (let i = 0; i < faces.size(); ++i) {
										let face = faces.get(i);
										let point1 = new cv.Point(face.x, face.y);
										let point2 = new cv.Point(face.x + face.width, face.y + face.height);
										cv.rectangle(src, point1, point2, [0, 255, 0, 255]);
										detectado = 1;
									}
								} catch (err) {
									console.log(err);
								}
								if(detectado==1){
									document.getElementById('Estado').innerHTML = 'Rostro detectado';
									document.getElementById('Estado').style.backgroundColor='#04aa6d';
								} else{
									document.getElementById('Estado').innerHTML = 'Rostro no detectado';
									document.getElementById('Estado').style.backgroundColor='#f44336';
								}
								/*cv.imshow('canvas_output', src);*/
								let delay = 1000 / FPS - (Date.now() - begin);
								setTimeout(processVideo, delay);
							} 
							setTimeout(processVideo, 0);
						}
					}
					
					/* temporizador que usa la variable global tiempo para contar los segundos */
					function temporizador() {
						if(stop==0){
							if(document.getElementById('Estado').textContent=='Rostro detectado'){
								tiempo = tiempo + 1;
								document.getElementById('Tiempo').innerHTML = tiempo;
							} else{
								tiempo = 0; //reiniciamos el contador si no detectamos rostro
								document.getElementById('Tiempo').innerHTML = tiempo;
							}
							/*cuando haya pasado 3 segundos de la deteccion de un rostro ejecutar()*/
							if(tiempo == 3){
								stop = 1;
								tiempo = 0;
								document.getElementById('Tiempo').innerHTML = tiempo;
								ejecutar(); 
							}
						}
						setTimeout('temporizador()',1000);
					}
					
					/* generamos el archivo de imagen sin el recuadro */
					function ejecutar(){
						let imageCanvas = document.createElement('canvas');
						let imageCtx = imageCanvas.getContext('2d');
						let v = document.getElementById('cam_input');
						imageCanvas.width = v.videoWidth;
						imageCanvas.height = v.videoHeight;
						imageCtx.drawImage(v, 0, 0, v.videoWidth, v.videoHeight);
						imageCanvas.toBlob(postFile, 'image/jpeg');
					}
					
					/* enviamos el 'file' y el 'identificador' a la url 'enrrolar' por metodo 'POST' */
					function postFile(file) {
						let formdata = new FormData();
						let identificador = prompt('Ingrese un identificador:');
						if(identificador==null){
							stop = 0;
							return
						}
						formdata.append('id', identificador);
						formdata.append('image', file);
						let xhr = new XMLHttpRequest();
						xhr.open('POST', 'enrrolar', true);
						xhr.onload = function () {
							if (this.status === 200){
								alert(this.response); //si se hizo un envio exitoso, sin error; mostramos la respuesta
								stop = 0;
							}
							else{
								alert(xhr); //si se hizo un envio con errores; mostramos el error
								stop = 0;
							}
						};
						xhr.onerror = function () {
							alert('Error de comunicacion con el servidor');
							stop = 0;
						};
						xhr.onabort = function () {
							alert('Peticion de reconocimiento abortada');
							stop = 0;
						};
						xhr.send(formdata);
					}
				</script>
			</body>
		</html>`)
	})
	log.Println("(http://localhost:5000) enrrolador biometrico")
	http.HandleFunc("/identificar.html", func(w http.ResponseWriter, r *http.Request) { 
		fmt.Fprintf(w, `<html oncontextmenu='return false' onkeydown='return false'>
			<head>
				<meta name='viewport' content='width=device-width, initial-scale=1.0'> 
				<title>Identificar</title>
				<script async src='opencv.js' onload='openCvReady();'></script>
				<script src='utils.js'></script>
				<style> .labele { color: white; padding: 8px; font-family: Arial; background-color: #ff9800; } .labelt { color: white; padding: 8px; font-family: Arial; background-color: #04aa6d; } .labeli { color: white; padding: 8px; font-family: Arial; background-color: #f44336; } @media only screen and (max-width: 992px) { video.camara { height:640px; width:480px; display: block; margin-left: auto; margin-right: auto; } } @media only screen and (min-width: 993px) { video.camara { height:480px; width:640px; display: block; margin-left: auto; margin-right: auto;} } </style>
			</head>
			<body  bgcolor='#000' onload="setTimeout('temporizador()',1000)"> <br>
				<center> 
					<span id='Tiempo' class='labelt'>0</span> <span id='Estado' class='labele'>Iniciando...</span> 
					<span id='Identificador' class='labeli'>Desconocido</span> 
				</center><br><br>
				<center>
					<! canvas id='canvas_output' /><! /canvas> 
					<video id='cam_input' height='480' width='640' class='camara'></video> 
				</center>
				<script>
					/* variables globales para el funcionamiento */
					let stop = 0; //variable que determina si detenemos el contador
					let tiempo = 0; //variable que sirve como contador de segundos
					
					/* aperturamos webcam con opencv */
					function openCvReady() {
						cv['onRuntimeInitialized'] = () => {
							let video = document.getElementById('cam_input');
							/*video.style.display='none';*/ 
							navigator.mediaDevices.getUserMedia({ video: true, audio: false }).then(function (stream) { video.srcObject = stream; video.play(); }).catch(function (err) { console.log('Error: ' + err); });
							let src = new cv.Mat(video.height, video.width, cv.CV_8UC4);
							let gray = new cv.Mat();
							let cap = new cv.VideoCapture(cam_input);
							let faces = new cv.RectVector();
							let faceClassifier = new cv.CascadeClassifier();
							let utils = new Utils('errorMessage');
							let faceCascade = 'haarcascade_frontalface_default.xml';
							utils.createFileFromUrl(faceCascade, faceCascade, () => { faceClassifier.load(faceCascade); });
							const FPS = 40;
							function processVideo() {
								let begin = Date.now();
								cap.read(src);
								cv.cvtColor(src, gray, cv.COLOR_RGBA2GRAY, 0);
								let detectado=0;
								try {
									faceClassifier.detectMultiScale(gray, faces, 1.1, 3, 0);
									for (let i = 0; i < faces.size(); ++i) {
										let face = faces.get(i);
										let point1 = new cv.Point(face.x, face.y);
										let point2 = new cv.Point(face.x + face.width, face.y + face.height);
										cv.rectangle(src, point1, point2, [0, 255, 0, 255]);
										detectado = 1;
									}
								} catch (err) {
									console.log(err);
								}
								if(detectado==1){
									document.getElementById('Estado').innerHTML = 'Rostro detectado';
									document.getElementById('Estado').style.backgroundColor='#04aa6d';
								} else{
									document.getElementById('Estado').innerHTML = 'Rostro no detectado';
									document.getElementById('Estado').style.backgroundColor='#f44336';
								}
								/*cv.imshow('canvas_output', src);*/
								let delay = 1000 / FPS - (Date.now() - begin);
								setTimeout(processVideo, delay);
							} 
							setTimeout(processVideo, 0);
						}
					}
					
					/* temporizador que usa la variable global tiempo para contar los segundos */
					function temporizador() {
						if(stop==0){
							if(document.getElementById('Estado').textContent=='Rostro detectado'){
								tiempo = tiempo + 1;
								document.getElementById('Tiempo').innerHTML = tiempo;
							} else{
								tiempo = 0; //reiniciamos contador si no detectamos rostro
								document.getElementById('Tiempo').innerHTML = tiempo;
							}
							/*cuando hayan pasado 3 segundos de la deteccion de un rostro ejecutar()*/
							if(tiempo == 3){
								stop = 1; //detenemos temporizador
								tiempo = 0; //reiniciamos contador
								document.getElementById('Tiempo').innerHTML = tiempo;
								document.getElementById('Identificador').innerHTML = 'Enviando al servidor';
								document.getElementById('Identificador').style.backgroundColor='#ff9800';
								ejecutar(); //enviamos imagen al servidor 
							}
						}
						setTimeout('temporizador()',1000); //volver a ejecutar la funcion, dentro de un segundo
					}
					
					/* generamos el archivo de imagen sin el recuadro */
					function ejecutar(){
						let imageCanvas = document.createElement('canvas');
						let imageCtx = imageCanvas.getContext('2d');
						let v = document.getElementById('cam_input');
						imageCanvas.width = v.videoWidth;
						imageCanvas.height = v.videoHeight;
						imageCtx.drawImage(v, 0, 0, v.videoWidth, v.videoHeight);
						imageCanvas.toBlob(postFile, 'image/jpeg');
					}
					
					/* enviamos el 'file' a la url 'identificar' por metodo 'POST' */
					function postFile(file) {
						let formdata = new FormData();
						formdata.append('image', file);
						let xhr = new XMLHttpRequest();
						xhr.open('POST', 'identificar', true);
						xhr.onload = function () {
							/*si se hizo un envio exitoso, sin error; enviamos la respuerta indentificacion()*/
							if (this.status === 200){
								if (this.response != "NO IDENTIFICADO"){ identificacion(this.response); }
								else{ 
									document.getElementById('Identificador').innerHTML = 'NO IDENTIFICADO';
									document.getElementById('Identificador').style.backgroundColor='#f44336';
									stop = 0; 
								}
							}
						};
						xhr.onerror = function () {
							document.getElementById('Identificador').innerHTML = 'Error de comunicacion';
							alert('Error de comunicacion con el servidor');
							document.getElementById('Identificador').innerHTML = 'NO IDENTIFICADO';
							document.getElementById('Identificador').style.backgroundColor='#f44336';
							stop = 0;
						};
						xhr.onabort = function () {
							document.getElementById('Identificador').innerHTML = 'Peticion abortada';
							alert('Peticion de reconocimiento abortada');
							document.getElementById('Identificador').innerHTML = 'NO IDENTIFICADO';
							document.getElementById('Identificador').style.backgroundColor='#f44336';
							stop = 0;
						};
						xhr.send(formdata);
					}
					
					/* mostramos el Identificador de la persona por 5 segundos */
					async function identificacion(res) {
						document.getElementById('Identificador').innerHTML = res;
						document.getElementById('Identificador').style.backgroundColor='#00008b';
						await sleep(5 * 1000); //hacemos una pausa de 5 segundos
						document.getElementById('Identificador').innerHTML = 'Desconocido';
						document.getElementById('Identificador').style.backgroundColor='#f44336';
						stop = 0;
					}
					
					/* funcion que simula sleep */
					function sleep(ms) {
						return new Promise((resolve) => setTimeout(resolve, ms));
					}
				</script>
			</body>
		</html>`)
	})
	log.Println("(http://localhost:5000/identificar.html) identificador biometrico")
	http.HandleFunc("/opencv.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "opencv.js")})
	//log.Println("(http://localhost:5000/opencv.js) libreria de deteccion de rostros")
	http.HandleFunc("/utils.js", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "utils.js")})
	//log.Println("(http://localhost:5000/utils.js) componentes extras para la libreria opencv")
	http.HandleFunc("/haarcascade_frontalface_default.xml", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "haarcascade_frontalface_default.xml")})
	//log.Println("(http://localhost:5000/haarcascade_frontalface_default.xml) modelo de deteccion de caras")
	http.Handle("/enrrolados/", http.StripPrefix("/enrrolados", (http.FileServer(http.Dir("./enrrolados"))) ))
	log.Println("(http://localhost:5000/enrrolados/) rostros enrrolados")
	http.HandleFunc("/enrrolar", func(w http.ResponseWriter, r *http.Request) { 
		identificador := r.FormValue("id")
		enrrolados, err := os.Stat("./enrrolados/"+identificador)
		if (err == nil) && !enrrolados.IsDir(){
			fmt.Fprintf(w, "Ya existe el identificador: '%s'\n",identificador)
			log.Println("Ya existe el identificador:",identificador)
			return
		}
		if identificador == "" || identificador == "null" {
			http.Error(w, "Error al obtener el identificador", http.StatusBadRequest)
			log.Println("Error al obtener el identificador")
			return
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Error al obtener el imagen", http.StatusBadRequest)
			log.Println("Error al obtener el imagen")
			file.Close()
			return
		}
		defer file.Close()
		dst, err := os.Create("./enrrolados/"+identificador)
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
		fmt.Fprintf(w, "Enrrolado Exitosamente '%s'\n",identificador)
		log.Println("Enrrolado Exitosamente ",identificador)
	})
	//log.Println("(http://localhost:5000/enrrolar) webservice de enrrolamiento")
	http.HandleFunc("/identificar", func(w http.ResponseWriter, r *http.Request) { 
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
	})
	//log.Println("(http://localhost:5000/identificar) webservice de identificacion")
	http.ListenAndServe(":5001", nil)
}
