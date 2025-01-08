# notes for now on how to train

# If you dont have the project set up yet

`python3 init.py`

**And for now normalize the dataset file names**

(fsr makefile doesnt like some of the characters in the base file names)

`cd ./filenormalizer/ && go run filenormalizer.go`

Running the training and model:

Note: I have commited the latest model and will update so you dont have to train.

```
# from tess_files dir but can change file pathing if u want to run somewhere else
# This trains the model
TESSDATA_PREFIX=../tesseract/tessdata make training MODEL_NAME=maplestory START_MODEL=eng TESSDATA=../tesseract/tessdata

# Then to run tesseract and see eval as stdout
TESSDATA_PREFIX=../tesseract/tessdata tesseract {image path to run on} stdout --tessdata-dir ./tesstrain/data --psm 13 -l maplestory --loglevel ALL
```

## Really important note

Update: I will commit the newest model so you dont need to run the trainer yourself!

When training is done it will generate a maplestory.traineddata file in ./tesstrain/data/maplestory/

This MUST BE copied into the ./tesseract/tessdata/ dir or it will not work.
