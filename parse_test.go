package godress

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/fatih/color"
)

func TestParse(t *testing.T) {
	address1 := &Address{
		Original:        "123 N CENTER ST LEHI, UT 84043",
		Hash:            "aa7118672d057b7ff9ec7fedb1abc850",
		HouseNumber:     "123",
		StreetDirection: "N",
		StreetName:      "CENTER",
		StreetType:      "ST",
		City:            "LEHI",
		State:           "UT",
		PostalCode:      "84043",
	}

	address2 := &Address{
		Original:        "137 N 800 E SPANISH FORK, UT 84660",
		Hash:            "f42142554958f74799bdb1e992a0a60e",
		HouseNumber:     "137",
		StreetDirection: "N",
		StreetName:      "800 E",
		City:            "SPANISH FORK",
		State:           "UT",
		PostalCode:      "84660",
	}

	address3 := &Address{
		Original:        "2505 NE 135TH ST, SEATTLE, WA 98125",
		Hash:            "7fa506f8000b944bf69dcfab008c6604",
		HouseNumber:     "2505",
		StreetDirection: "NE",
		StreetName:      "135TH",
		StreetType:      "ST",
		City:            "SEATTLE",
		State:           "WA",
		PostalCode:      "98125",
	}

	address4 := &Address{
		Original:        "PO BOX 523029 WEST CHESTER, PA 18630",
		Hash:            "886b9d089a09b109643a477eb84e9ac7",
		HouseNumber:     "523029",
		StreetDirection: "",
		StreetName:      "PO BOX",
		StreetType:      "",
		City:            "WEST CHESTER",
		State:           "PA",
		PostalCode:      "18630",
	}

	tests := map[string]*Address{
		"123 N Center St Lehi, UT 84043":       address1,
		"123 N Center St. Lehi, UT 84043":      address1,
		"137 N 800 E Spanish Fork, UT 84660":   address2,
		"2505 NE 135th St, Seattle, WA 98125":  address3,
		"PO BOX 523029 West Chester, PA 18630": address4,
	}

	for s, a := range tests {
		if pa, err := Parse(s); err != nil {
			t.Errorf("error testing %s: %v", s, err)
		} else if *pa != *a {
			prettyPrint(t, a, pa)
		}
	}
}

func prettyPrint(t *testing.T, testInput, testOutput *Address) {
	var (
		expectedData, gotData []byte
		err                   error
	)

	if expectedData, err = json.MarshalIndent(testInput, "", "  "); err != nil {
		log.Fatalf(err.Error())
	}

	if gotData, err = json.MarshalIndent(testOutput, "", "  "); err != nil {
		log.Fatalf(err.Error())
	}

	t.Errorf("\nexpected: %s\n got: %s\n\n", color.GreenString(string(expectedData)), color.RedString(string(gotData)))
}
