# notes for now on how to train

# If you dont have the project set up yet

`python3 init.py`

**And for now normalize the dataset file names**

(fsr makefile doesnt like some of the characters in the base file names)

`cd ./filenormalizer/ && go run filenormalizer.go`

```
TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=maplestory START_MODEL=eng TESSDATA=../tesseract/tessdata

# to run tesseract
TESSDATA_PREFIX=../tesseract/tessdata tesseract {image path to run on} stdout --tessdata-dir ./tesstrain/data --psm 7 -l maplestory --loglevel ALL
```

## Really important note

When training is done it will generate a maplestory.traineddata file in ./tesstrain/data/maplestory/

This MUST BE copied into the ./tesseract/tessdata/ dir or it will not work.
