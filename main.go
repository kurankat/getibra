package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kurankat/csvdict"
)

var outputHeaders = []string{
	"Locality_Name",
	"Variants",
	"bioregion",
	"State",
	"Country",
	"Lat/Long_Method",
	"Latitude_1",
	"Longitude_1",
	"Datum",
}

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
	Field       string `json:"field"`
	Description string `json:"description"`
	Layername   string `json:"layername"`
	Pid         string `json:"pid"`
	Value       string `json:"value"`
}

var alaClient = &http.Client{}

func main() {
	importFile, err := os.Open("Tasmania_test.csv")
	if err != nil {
		panic(err)
	}
	defer importFile.Close()

	exportFile, err := os.OpenFile("ibraLocalities.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer exportFile.Close()

	localityReader, err := csvdict.NewDictReader(importFile)
	if err != nil {
		panic(err)
	}

	localityWriter := csvdict.NewDictWriter(exportFile, outputHeaders)
	localityWriter.WriteHeaders()

	for {
		// Read line into memory
		record, err := localityReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		lat := record["Latitude_1"]
		long := record["Longitude_1"]

		record["bioregion"] = getBioregion(lat, long)

		localityRecord := newLocality(record)
		// for key, value := range record {
		// 	fmt.Printf("\"%s\": \"%s\"\n", key, value)
		// }

		fmt.Println(record)
		fmt.Println("Locality Name (map): ", record["Locality_Name"])
		fmt.Println("Locality Name (struct): ", localityRecord.localityName)

		// fmt.Println(outputHeaders)

		localityWriter.Write(record)
	}
	localityWriter.Flush()
}

func newLocality(record map[string]string) *LocalityData {
	locData := &LocalityData{
		localityName: record["Locality_Name"],
		variants:     record["Variants"],
		country:      record["Country"],
		state:        record["State"],
		lat:          record["Latitude_1"],
		long:         record["Longitude_1"],
		llMethod:     record["Lat/Long_Method"],
		datum:        record["Datum"],
	}

	locData.bioregion = getBioregion(locData.lat, locData.long)
	return locData
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
