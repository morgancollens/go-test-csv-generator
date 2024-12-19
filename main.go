package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

var validFields = map[string]bool{
	"name":       true,
	"age":        true,
	"email":      true,
	"firstName":  true,
	"lastName":   true,
	"middleName": true,
	"city":       true,
	"jobTitle":   true,
}

var generators = map[string]func(BaseFields) string{
	"name":       func(fields BaseFields) string { return fields.Name },
	"age":        func(fields BaseFields) string { return strconv.Itoa(gofakeit.Number(18, 99)) },
	"email":      func(fields BaseFields) string { return fields.Email },
	"firstName":  func(fields BaseFields) string { return fields.FirstName },
	"lastName":   func(fields BaseFields) string { return fields.LastName },
	"middleName": func(fields BaseFields) string { return gofakeit.MiddleName() },
	"city":       func(fields BaseFields) string { return gofakeit.City() },
	"jobTitle":   func(fields BaseFields) string { return gofakeit.JobTitle() },
}

type BaseFields struct {
	Name      string
	FirstName string
	LastName  string
	Email     string
}

func validateFlags(rows int, fields string, filename string) error {
	if rows <= 0 {
		return fmt.Errorf("invalid number of rows: %d", rows)
	}

	if fields == "" {
		return fmt.Errorf("fields cannot be empty")
	}

	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	return nil
}

func validateSelectedFields(fields string) []string {
	var invalidFields []string
	fieldSlice := strings.Split(fields, ",")
	for _, userField := range fieldSlice {
		if !validFields[userField] {
			invalidFields = append(invalidFields, userField)
		}
	}

	return invalidFields
}

// To maintain consistency between certain fields, base fields are generated for each row
// regardless of whether they are included in the fields list.
func generateBaseFields() BaseFields {
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	emailDomain := gofakeit.DomainName()
	name := fmt.Sprintf("%s %s", firstName, lastName)
	email := fmt.Sprintf("%s.%s@%s", strings.ToLower(firstName), strings.ToLower(lastName), emailDomain)

	return BaseFields{
		Name:      name,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
}

func generateCsvData(rows int, fields string, outputDir string, filename string) error {
	fieldSlice := strings.Split(fields, ",")

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	filePath := filepath.Join(outputDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(fieldSlice); err != nil {
		return fmt.Errorf("failed to write header row: %v", err)
	}

	for i := 0; i < rows; i++ {
		row := []string{}
		baseFields := generateBaseFields()
		for _, field := range fieldSlice {
			row = append(row, generators[field](baseFields))
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %v", err)
		}
	}

	return nil
}

func main() {
	rows := flag.Int("rows", 1, "Number of rows to include in the generated CSV file.")
	fields := flag.String("fields", "name,age", "Comma separated list of fields (ex. 'name,age,email') to include in the generated CSV file.")
	filename := flag.String("filename", "output.csv", "Name of the file to write the generated CSV data to.")
	seed := flag.Int("seed", 0, "Seed for random number generation.")

	flag.Parse()

	startTime := time.Now()

	if err := validateFlags(*rows, *fields, *filename); err != nil {
		log.Fatalf("Invalid flags: %v", err)
	}

	invalidFields := validateSelectedFields(*fields)
	if len(invalidFields) > 0 {
		log.Fatalf("Unable to generate CSV data. Invalid fields selected: %s", strings.Join(invalidFields, ", "))
	}

	fmt.Printf("Rows: %d\n", *rows)
	fmt.Printf("Fields: %s\n", *fields)
	fmt.Printf("Filename: %s\n", *filename)
	fmt.Printf("Generating CSV file...\n")

	gofakeit.Seed(*seed)

	outputDir := "output"
	generateCsvData(*rows, *fields, outputDir, *filename)

	elapsed := time.Since(startTime)

	fmt.Printf("CSV file generated at %s/%s. (Elapsed time: %f seconds)\n", outputDir, *filename, elapsed.Seconds())
}
