#!/bin/bash
TESSDATA_PREFIX=./tesseract/tessdata/
cd ./tesstrain/
make training MODEL_NAME=maplestory START_MODEL=eng TESSDATA=../tesseract/tessdata
