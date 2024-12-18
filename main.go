package main

import (
	"flag"
	"fmt"
)

func main() {
	rows := flag.Int("rows", 1, "Number of rows to include in the generated CSV file.")
	fields := flag.String("fields", "name,age", "Comma separated list of fields (ex. 'name, age, email') to include in the generated CSV file.")
	filename := flag.String("filename", "output.csv", "Name of the file to write the generated CSV data to.")

	flag.Parse()

	fmt.Printf("Rows: %d\n", *rows)
	fmt.Printf("Fields: %s\n", *fields)
	fmt.Printf("Filename: %s\n", *filename)
}
