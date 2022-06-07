package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

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
var outputHeaders = []string{
	"Locality Name",
	"Variants",
	"bioregion",
	"State",
	"Country",
	"Lat/Long Method",
	"Latitude 1",
	"Longitude 1",
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
	getBioregion("-39.93125", "143.85099")
}

func getJson(url string, target interface{}) error {
	r, err := alaClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	jsonResp, _ := ioutil.ReadAll(r.Body)
	stringResp := string(jsonResp)
	trimmed1 := strings.ReplaceAll(stringResp, "[", "")
	trimmed2 := strings.ReplaceAll(trimmed1, "]", "")

	forDecoder := strings.NewReader(trimmed2)
	return json.NewDecoder(forDecoder).Decode(target)
}

func getBioregion(lat, long string) (bioregion string) {
	region := &Response{}

	requestURL := fmt.Sprintf("https://spatial.ala.org.au/ws/intersect/1048/%s/%s", lat, long)
	getJson(requestURL, region)

	return region.Value
}
