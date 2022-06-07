package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

var testLocation = LocalityData{
	localityName: "Currie",
	country:      "Australia",
	state:        "State of Tasmania",
	lat:          "-39.93125",
	long:         "143.85099",
	llMethod:     "Decimal degrees",
	datum:        "GDA94",
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestGetBioregion(t *testing.T) {
	equals(t, "King", getBioregion(testLocation.lat, testLocation.long))
}
