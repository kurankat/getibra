package main

import "fmt"

var inputHeaders = []string{
	"Locality Name",
	"Variants",
	"State",
	"Country",
	"Lat/Long Method",
	"Latitude 1",
	"Longitude 1",
	"Datum",
}
var outputHeaders []string

type LocalityData struct {
	localityName string
	variants     string
	country      string
	state        string
	bioregion    string
	lat          string
	long         string
	llMethod     string
	datum        string
}

type Request struct {
	lat     string
	lon     string
	layerID string
}

type Response struct {
	field       string
	description string
	layername   string
	pid         string
	value       string
}

func main() {
	fmt.Println("Main")
}
