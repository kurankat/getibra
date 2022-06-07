package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var outputHeaders = []string{
	"Locality Name",
	"Variants",
	"State",
	"Country",
	"Lat/Long Method",
	"Latitude 1",
	"Longitude 1",
	"Datum",
	"bioregion",
}

type Request struct {
	lat     string
	lon     string
	layerID string
}

type Response struct {
	Field       string `json:"field"`
	Description string `json:"description"`
	Layername   string `json:"layername"`
	Pid         string `json:"pid"`
	Value       string `json:"value"`
}

var alaClient = &http.Client{}

func main() {
	timeStamp := time.Now().Format("20060102T150405")
	importFile, err := os.Open("Tasmania_test.csv")
	dealWith(err)
	defer importFile.Close()

	exportName := fmt.Sprintf("ibraLocalities-%s.csv", timeStamp)
	exportFile, err := os.OpenFile(exportName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	dealWith(err)
	defer exportFile.Close()

	localityReader := csv.NewReader(importFile)
	dealWith(err)
	localityReader.Comma = ','

	localityWriter := csv.NewWriter(exportFile)
	localityWriter.Comma = ','
	localityWriter.Write(outputHeaders)

	// Read and discard header row
	headerRow, err := localityReader.Read()
	dealWith(err)

	if len(headerRow) != len(outputHeaders)-1 {
		panic("CSV file a sthe wrong number of fields")
	}

	for {
		record, err := localityReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		bioregion := getBioregion(record[5], record[6])
		record = append(record, bioregion)

		fmt.Println(record)
		localityWriter.Write(record)
	}
	localityWriter.Flush()
}

func getBioregion(lat, long string) (bioregion string) {
	region := &Response{}

	requestURL := fmt.Sprintf("https://spatial.ala.org.au/ws/intersect/1048/%s/%s", lat, long)
	_ = getJson(requestURL, region)

	return region.Value
}

func getJson(url string, target interface{}) error {
	r, err := alaClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	rawResp, _ := ioutil.ReadAll(r.Body)
	jsonResp := bytes.ReplaceAll(bytes.ReplaceAll(rawResp, []byte("["), []byte("")),
		[]byte("]"), []byte(""))

	jsonReader := bytes.NewReader(jsonResp)
	return json.NewDecoder(jsonReader).Decode(target)
}

func dealWith(err error) {
	if err != nil {
		panic(err)
	}
}
