# markdata

tool to help mark training data faster (hopefully)

## Running

`go run markdata.go`

each image will use xdg-open to show the cropped preview and print out the vim command to open a new text writer.

for gaps between words that are not spaces use 5 spaces (ie; username->storename->quantity...)

once you have written the text file close the xdg-open window and move onto the next.

to "skip" an image simply close the window without writting a vim file and it will be pruned from the dataset.
