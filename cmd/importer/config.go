package main

import "github.com/namsral/flag"

type Config struct {
	ApiURL     string
	InputFile  string
	OutputFile string
}

func NewConfig() Config {
	var apiURL, inputFile, outputFile string

	flag.StringVar(&inputFile, "input", "", "File containing file paths to import")
	flag.StringVar(&outputFile, "output", "", "Write unsuccessfully imported PDFs to")
	flag.StringVar(&apiURL, "api_url", "http://localhost", "The URL of the API")
	flag.Parse()

	return Config{
		ApiURL:     apiURL,
		InputFile:  inputFile,
		OutputFile: outputFile,
	}
}