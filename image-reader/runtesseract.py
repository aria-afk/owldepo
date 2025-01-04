import subprocess

if __name__ == "__main__":
    subprocess.run(["TESSDATA_PREFIX=", "../trainer/tess_files/tesseract/tessdata/"])
