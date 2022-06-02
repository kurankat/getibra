package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

var testLocation = LocalityData{
	localityName: "Currie",
	country:      "Australia",
	state:        "State of Tasmania",
	lat:          "-39.93125",
	long:         "143.85099",
	llMethod:     "Decimal degrees",
	datum:        "GDA94",
}

func Test_Main(t *testing.T) {
	for key, value := range testLocation {
		fmt.Printf("Header: %s,\tValue: %s", key, value)
	}

	equals(t, true, true)
}

// Locality Name	Variants	State	Country	Lat/Long Method	Latitude 1	Longitude 1	Datum
// Currie		State of Tasmania	Australia	Decimal degrees	-39.93125	143.85099	GDA94

// https://spatial.ala.org.au/ws/intersect/1048/-23.1/149.1
