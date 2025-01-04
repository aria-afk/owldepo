import os 
import subprocess

if __name__ == "__main__":
    # Download git repos
    if os.path.isdir("./tesseract/") is not True:
        subprocess.run(["git", "clone", "https://github.com/tesseract-ocr/tesseract.git"])
    if os.path.isdir("./tesstrain/") is not True:
        subprocess.run(["git", "clone", "https://github.com/tesseract-ocr/tesstrain.git"])
    # Copy eng train data to tesseract/tessdata
    if os.path.isfile("./tesseract/tessdata/eng.traineddata") is not True:
        subprocess.run(["mv", "./eng.traineddata", " ", "./tesseract/tessdata/"])
    # get training data 
    # TODO: This should eventually come from s3 instead of local 
    if os.path.isdir("./tesstrain/data/maplestory-ground-truth/") is not True:
        if os.path.isdir("./tesstrain/data/") is not True:
            subprocess.run(["mkdir", "./tesstrain/data"])
        subprocess.run(["cp", "-r", "../data-marker/maplestory-ground-truth/", "./tesstrain/data/"])
