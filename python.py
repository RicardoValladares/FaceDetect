from flask import Flask, request, Response
import time

app = Flask(__name__)

@app.route('/')
def principal():
    return Response(open('./web/index.html').read(), mimetype="text/html")

@app.route('/index.html')
def index():
    return Response(open('./web/index.html').read(), mimetype="text/html")

@app.route('/opencv.js')
def opencvjs():
    return Response(open('./web/opencv.js').read(), mimetype="text/javascript")

@app.route('/utils.js')
def utiljs():
    return Response(open('./web/utils.js').read(), mimetype="text/javascript")

@app.route('/haarcascade_frontalface_default.xml')
def frontalfacexml():
    return Response(open('./web/haarcascade_frontalface_default.xml').read(), mimetype="application/xml")

@app.route('/haarcascade_eye.xml')
def eyexml():
    return Response(open('./web/haarcascade_eye.xml').read(), mimetype="application/xml")


if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')
