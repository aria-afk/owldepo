# data-marker

Tool used to simplify the creation process of training data.

Expected training data will come from the output of running the scrapper.

## TODO

- Make a data set of unique strings and not duplicates
 ie; dont put in lots of Power elixir ones
 so we can check existing vim txt files for the entered string or sth idk

- Make this windows compatible for Rushi

- Would be nice to have vim buffer open automatically

## Usage

To begin the data-marker run `go run data-marker.go` (in this dir)

The default amount of images to parse is `10` to change this run with the `-iterations` flag. Example: `go run data-marker.go -iterations 100`

An image window will open using `xdg-open` and the script will stdout a vim command
to open a writer for the image file.

**README: skipping a file**

to skip a file simply close the image preview and DO NOT create a vim text file.
this will cause the preview to be automatically pruned and not entered into the
training data set.

Copy the vim command in a terminal and run it, then write the text exactly as seen
in the open image and `:wq` the file.
