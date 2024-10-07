# Trail Finder

Trail Finder is a command-line tool designed to parse trail data from a CSV file and filter the trails based on various criteria like restroom availability, picnic areas, ADA facilities, and more. The project includes concurrent CSV parsing, a flexible filtering system, and unit tests to ensure reliable functionality.

## Features

- **Concurrent CSV Parsing**: Efficiently reads large CSV files in chunks and processes the data concurrently.
- **Trail Filtering**: Filters trails based on various amenities such as restrooms, picnic areas, ADA facilities, fishing spots, and more.
- **Command-line Interface (CLI)**: Offers a flexible CLI for users to filter trails by passing command-line flags.
- **HTTP Server Mode**: Optionally run the application as a web server to fetch trail data via HTTP requests.
- **Unit Tests**: Includes comprehensive unit tests for both CSV parsing and filtering.

## Installation

To get started with Trail Finder, clone the repository and build the project using Go or the provided Makefile.

### Using Go

```bash
git clone <repository_url>
cd trailfinder
go build -o trailfinder main.go
```

### Using Makefile

Alternatively, you can use the Makefile to build the project. Run the following commands:

```bash
git clone <repository_url>
cd trailfinder
make
```

The Makefile includes the following targets:

- **build**: Compiles the project and outputs the binary to the `bin` directory.
- **test**: Runs the unit tests.
- **clean**: Cleans up build artifacts.

## Usage

You can run the tool using the command-line interface. Pass the path to the CSV file and specify the filter criteria using flags. You can also run the application in server mode to fetch trail data via HTTP requests.

### Example Command

```bash
./trailfinder --csv trailheads.csv --restrooms Yes --picnic No
```

### Running in Server Mode

To run the application in server mode, use the `--server` flag. You can specify filter criteria as query parameters in the HTTP request.

### Example Server Command

```bash
./trailfinder --server --port 8080
```

### Example HTTP Request

You can then fetch filtered trails using a command like:

```bash
curl "http://localhost:8080/trails?restrooms=Yes&picnic=No"
```

### Available Flags

- `--csv` : Path to the trailheads CSV file (default: `trailheads.csv`)
- `--restrooms` : Filter trails based on restroom availability (`Yes` or `No`)
- `--picnic` : Filter trails based on picnic area availability (`Yes` or `No`)
- `--fishing` : Filter trails based on fishing spot availability
- `--fee` : Filter trails based on whether fees are applicable
- `--bikerack` : Filter trails with bike rack availability
- `--adatoilet` : Filter trails with ADA-compliant toilets
- `--adapicnic` : Filter trails with ADA-compliant picnic areas
- `--adatrail` : Filter trails with ADA-compliant trails
- `--horsetrail` : Filter trails with horse trail availability
- `--recyclebin` : Filter trails with recycle bin availability
- `--dogcompost` : Filter trails with dog compost availability
- `--accessname` : Filter trails by access name
- `--thleash` : Filter trails by TH leash availability
- `--server` : Run the application in server mode to fetch trail data via HTTP requests
- `--port` : Specify the port for the HTTP server (default: `8080`)

## Development

### Project Structure

- **csvparser.go**: Contains the core CSV parsing logic, implemented to handle large datasets concurrently.
- **filter.go**: Implements the filtering logic, enabling users to filter trails by specific amenities.
- **main.go**: The entry point for the CLI application, where users specify filter options and can run the server.
- **csvparser_test.go**: Unit tests for CSV parsing, validating the concurrent parsing mechanism.
- **filter_test.go**: Unit tests for the filtering functionality, ensuring trails are filtered accurately.

### Running Tests

Run the unit tests to verify the correctness of the CSV parser and filtering logic:

```bash
go test ./...
```

## Design Overview

1. **CSV Parsing**: The project uses concurrent CSV parsing to efficiently process large datasets. Trails are read from the CSV file in chunks, and results are streamed through channels.
2. **Trail Filtering**: The filter logic allows for flexible filtering based on user-specified criteria. Multiple filters can be applied in combination. The filtering function can filter trails in parallel, improving performance for large datasets.
3. **Command-line Interface**: The CLI allows users to specify a CSV file and set filtering criteria. The program reads the CSV file, applies filters, and outputs the filtered results.
4. **HTTP Server**: The application can also run as a web server, allowing users to fetch filtered trail data via HTTP requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
