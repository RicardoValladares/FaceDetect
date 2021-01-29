from flask import Flask, request, Response
import time

PATH_TO_TEST_IMAGES_DIR = './upload'
app = Flask(__name__)

@app.route('/')
def index():
    return Response(open('./web/page.html').read(), mimetype="text/html")

@app.route('/upload', methods=['POST'])
def image():
    i = request.files['image']
    f = ('%s.jpeg' % time.strftime("%Y%m%d-%H%M%S"))
    i.save('%s/%s' % (PATH_TO_TEST_IMAGES_DIR, f))
    return Response("%s saved" % f)

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0')
