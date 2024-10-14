#!/bin/bash

# Create the output directory on the desktop if it doesn't exist
mkdir -p ~/Desktop/go_txt

# Find all .go files in the current directory and all subdirectories
find . -type f -name "*.go" | while read -r file; do
    # Get the filename without the path and .go extension
    filename=$(basename "${file%.go}")
    # Get the relative path of the file (excluding the filename)
    relpath=$(dirname "$file")
    # Create a unique name for the output file
    # This prevents overwriting files with the same name from different directories
    unique_name="${relpath//\//_}_${filename}"
    # Remove the leading ._ from the unique name if present
    unique_name="${unique_name#._}"
    # Concatenate the content of the .go file to a new .txt file in the go_txt directory on the desktop
    cat "$file" > "$HOME/Desktop/go_txt/${unique_name}.txt"
    echo "Processed $file -> $HOME/Desktop/go_txt/${unique_name}.txt"
done

echo "All .go files have been processed. Output files are in the 'go_txt' directory on your desktop."
