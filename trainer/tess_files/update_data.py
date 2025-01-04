import os
import subprocess

if __name__ == "__main__":
    # update new training data set to tesstrain 
    # TODO: again this should be s3 eventually
    if os.path.isdir("./tesstrain/data/maplestory-ground-truth/") is True:
        subprocess.run(["cp", "-r", "../data-marker/maplestory-ground-truth/", "./tesstrain/data/"])
