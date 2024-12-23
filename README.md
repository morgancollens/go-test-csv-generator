# Go Test CSV Generator

## Overview

Go Test CSV Generator is a simple tool designed to help quickly generate CSV files locally with mock data spanning millions of rows. This can be useful for testing, data analysis, and performance benchmarking.

## Features

- Generate CSV files with customizable mock data
- Support for millions of rows
- Easy to use and configure

## Installation

To install the Go Test CSV Generator, clone the repository and build the project:

```bash
git clone https://github.com/morgancollens/go-test-csv-generator.git
cd go-test-csv-generator
go build
```

## Usage

To generate a CSV file with mock data, run the following command:

```bash
./go-test-csv-generator -rows=1000000 -fields=name,age,email -filename=test_data.csv -seed=18283
```

### Command Line Options

- `-rows`: Number of rows to generate (default: 1)
- `-fields`: List of fields (or columns) to output data for (default: name,age)
- `-filename`: Output file name (default: output.csv)
- `-seed`: A number that can be used to generate consistent output instead of randomized output (default: 0)

### Supported fields

Currently the tool supports generation of the following fields:
- `name`
- `age`
- `email`
- `firstName`
- `lastName`
- `middleName`
- `city`
- `jobTitle`

## How to run tests

```bash
    cd go-test-csv-generator
    go test -v
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the ISC License. See the [LICENSE](LICENSE) file for details.
