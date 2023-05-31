import base64
import os
import time
import face_recognition
from flask import Flask, request, Response

PATH_IMAGES_DIR = './enrrolados'
enrrolados_faces = []
res = os.listdir(PATH_IMAGES_DIR)
app = Flask(__name__)


@app.route('/')
def enrrolarhtml():
    return Response(open('enrrolar.html').read(), mimetype="text/html")


@app.route('/identificar.html')
def identificarhtml():
    return Response(open('identificar.html').read(), mimetype="text/html")


@app.route('/utils.js')
def utilsjs():
    return Response(open('utils.js').read(), mimetype="text/javascript")


@app.route('/opencv.js')
def opencvjs():
    return Response(open('opencv.js').read(), mimetype="text/javascript")


@app.route('/haarcascade_frontalface_default.xml')
def haarcascade_frontalface_defaultxml():
    return Response(open('haarcascade_frontalface_default.xml').read(), mimetype="application/xml")


@app.route('/enrrolar', methods=['POST'])
def enrrolar():
    identificador = request.form.get("id")
    file = request.files['image']
    try:
        if res.index(identificador) != None:
            return Response("Ya existe el identificador: " + identificador)
        else:
            raise ValueError('Error al encontrar indice del dato')
            raise Exception('Error al encontrar indice del dato')
    except:
        try:
            if file.filename != '':
                file.save('%s/%s' % (PATH_IMAGES_DIR, identificador))
                face_image = face_recognition.load_image_file(PATH_IMAGES_DIR + '/' + identificador)
                unk_face_encoding = face_recognition.face_encodings(face_image)[0]
                enrrolados_faces.append(unk_face_encoding)
                res.append(identificador)
                return Response("Enrrolado Exitosamente " + identificador)
            else:
                return Response("Error al obtener el imagen")
        except:
            return Response("Error al enrrolar")


@app.route('/identificar', methods=['POST'])
def identificar():
    try:
        file = request.files['image']
        nuevonombre = ('%s' % time.strftime("%Y%m%d%H%M%S"))
        file.save('%s/%s' % ("./temp", nuevonombre))
        face_image = face_recognition.load_image_file("./temp" + '/' + nuevonombre)
        unk_face_encoding = face_recognition.face_encodings(face_image)[0]
        results = face_recognition.compare_faces(enrrolados_faces, unk_face_encoding, tolerance=0.4)
        return Response(res[results.index(True)])
    except:
        return Response("NO IDENTIFICADO")


if __name__ == '__main__':
    for x in res:
        face_image = face_recognition.load_image_file(PATH_IMAGES_DIR + '/' + x)
        try:
            unk_face_encoding = face_recognition.face_encodings(face_image)[0]
            enrrolados_faces.append(unk_face_encoding)
            print("Imagen cargada: " + x)
        except:
            print("Imagen no valida: " + x)
            quit()
    print("http://localhost:5001")
    app.run(debug=True, host='0.0.0.0', port=5001)
