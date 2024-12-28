# quick notes on setting up tesseract and how to train

1) clone tesseract repo and follow the installion guide

<https://github.com/tesseract-ocr/tesseract>

<https://tesseract-ocr.github.io/tessdoc/Compiling.html>

2) clone tesstrain

<https://github.com/tesseract-ocr/tesstrain>

3) download eng.traineddata from their models and put it in tesseract

<https://github.com/tesseract-ocr/tessdata_best>

put it in ./tesseract/tessdata

TODO: I should and will rename the parser folder to something else

4) get (or make) the training data from the markdata tool

read the readme in markdata tool or ask me to get data (dependant on scrapper too)

5) copy the training data to tesstrain

mkdir ./tesstrain/data/

cp ./owldepo/parser/markdata/out ./tesstrain/data/

mv ./tesstrain/data/out ./tesstrain/data/maplestory-ground-truth

6) run the thing

`TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=maplestory START_MODEL=eng  TESSDATA=../tesseract/tessdata`

can read more about options.
