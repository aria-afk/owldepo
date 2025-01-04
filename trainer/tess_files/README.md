# notes for now on how to train

```
TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=maplestory START_MODEL=eng TESSDATA=../tesseract/tessdata

# to run tesseract
TESSDATA_PREFIX=../tesseract/tessdata tesseract {image path to run on} stdout --tessdata-dir ./tesstrain/data --psm 7 -l maplestory --loglevel ALL
```

## Really important note

When training is done it will generate a maplestory.traineddata file in ./tesstrain/data/maplestory/

This MUST BE copied into the ./tesseract/tessdata/ dir or it will not work.
