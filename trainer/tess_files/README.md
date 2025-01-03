# notes for now on how to train

```
TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=maplestory START_MODEL=eng TESSDATA=../tesseract/tessdata

# to run tesseract
tesseract {image path to run on} stdout --tessdata-dir ./tesstrain/data --psm 7 -l maplestory --loglevel ALL
```
