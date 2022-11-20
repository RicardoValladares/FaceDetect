# Rostro Biométrico
Este proyecto es un servidor web hecho para la detección y reconocimiento de rostros; implementando un conjunto de herramientas golang, javascript y C++.

<p align="center">
  <img alt="Light" src="https://raw.githubusercontent.com/RicardoValladares/FaceDetect/facedetection/desktop.png" width="50%">
  &nbsp; &nbsp; &nbsp; &nbsp;
  <img alt="Dark" src="https://raw.githubusercontent.com/RicardoValladares/FaceDetect/facedetection/celphone.png" width="15%">
</p>


Librerías usadas:
- Go-Face (https://github.com/Kagami/go-face)
- Go-Recognizer (https://github.com/leandroveronezi/go-recognizer)
- OpenCV (https://docs.opencv.org/3.4/d5/d10/tutorial_js_root.html)
- Dlib (http://dlib.net/)

Modelos entrenados para el reconocimiento, detección y predicción de rostros:
- @davisking (https://github.com/davisking/dlib-models)

## Notas de: @davisking (https://github.com/davisking/dlib-models)

### Reconocimiento - dlib_face_recognition_resnet_model_v1.dat
La red se entrenó desde cero con un conjunto de datos de aproximadamente 3 millones de rostros. Este conjunto de datos se deriva de dos conjuntos de datos. El conjunto de datos de face scrub (http://vintage.winklerbros.net/facescrub.html) y el conjunto de datos de VGG (http://www.robots.ox.ac.uk/~vgg/data/vgg_face/). Se hizo este modelo entrenando una CNN de reconocimiento facial y luego usando métodos de agrupación de gráficos y mucha revisión manual para limpiar el conjunto de datos. Al final, aproximadamente la mitad de las imágenes son de VGG y face scrub. Además, el número total de identidades individuales en el conjunto de datos es 7485. Se evitaron superposiciones de identidades en LFW.

### Detección - mmod_human_face_detector.dat
Este modelo está entrenado con el conjunto de datos: http://dlib.net/files/data/dlib_face_detection_dataset-2016-09-30.tar.gz. Se creó el conjunto de datos encontrando imágenes de rostros en muchos conjuntos de datos de imágenes disponibles públicamente (excluyendo el conjunto de datos FDDB). En particular, hay imágenes de ImageNet, AFLW, Pascal VOC, el conjunto de datos VGG, WIDER y face scrub. Todas las anotaciones en el conjunto de datos fueron creadas usando la herramienta imglab de dlib.

### Predicción - shape_predictor_5_face_landmarks.dat
Este es un modelo de referencia de 5 puntos que identifica las esquinas de los ojos y la parte inferior de la nariz. Se entrenó con el conjunto de datos que se encuentra en http://dlib.net/files/data/dlib_faces_5points.tar, que consta de 7198 rostros. Se creó este conjunto de datos descargando imágenes de Internet y anotándolas con la herramienta imglab de dlib. Este modelo está diseñado para funcionar bien con el detector de rostros HOG de dlib y el detector de rostros CNN (mmod_human_face_detector.dat).
