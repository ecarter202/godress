package godress

import (
	"regexp"
	"strconv"
	"strings"
)

type Address struct {
	Original        string  `arango:"original" json:"original"`
	HouseNumber     int     `arango:"house_number" json:"house_number"`
	StreetDirection string  `arango:"street_direction" json:"street_direction"`
	StreetName      string  `arango:"street_name" json:"street_name"`
	StreetType      string  `arango:"street_type" json:"street_type"`
	Unit            string  `arango:"unit" json:"unit"`
	City            string  `arango:"city" json:"city"`
	County          string  `arango:"county" json:"county"`
	State           string  `arango:"state" json:"state"`
	Zip             int     `arango:"zip" json:"zip"`
	Latitude        float64 `arango:"latitude" json:"latitude"`
	Longitude       float64 `arango:"longitude" json:"longitude"`
}

const (
	SMALLEST_ZIP_CODE = "01001"
	LARGEST_ZIP_CODE  = "99950"
	NUMBER_REGEX      = `(?m)(\d+)`
)

var (
	StreetTypes      = []string{"ALY", "ANX", "AVE", "BLVD", "CIR", "CT", "CV", "CRES", "DR", "EXPY", "EXT", "GRV", "HWY", "HL", "KY", "LN", "LOOP", "MALL", "PARK", "PKWY", "PL", "PLZ", "PT", "RD", "ROW", "RUN", "SQ", "ST", "TER", "TRCE", "TRL", "WAY", "ZZ"}
	StreetDirections = []string{"N", "NW", "NE", "S", "SW", "SE", "E", "W"}
)

// Parses string into separate parts.
func Parse(address string) (*Address, error) {
	a := &Address{}
	a.Original = address
	addressX := strings.FieldsFunc(address, split)

	for idx := 0; idx < len(addressX); idx++ {

		s := addressX[idx]
		s = strings.TrimSpace(s)

		if idx == 0 && !IsPoBox(address) {
			if isInt(s) {
				a.HouseNumber, _ = strconv.Atoi(s)
			}
		} else if IsPoBox(address) && a.HouseNumber == 0 {
			re := regexp.MustCompile(NUMBER_REGEX)
			numbers := re.FindAllString(strings.Replace(address, " ", "", -1), -1)
			if len(numbers) > 0 {
				a.HouseNumber, _ = strconv.Atoi(numbers[0])
			}
			a.StreetName = "PO Box"
			idx += 2
		} else if IsStreetDirection(s) && a.StreetDirection == "" {
			a.StreetDirection = s
		} else if IsStreetType(s) && a.Unit == "" {
			a.StreetType = s
		} else if isApartmentKeyword(s) && a.Unit == "" {
			if idx+1 < len(addressX) {
				a.Unit = addressX[idx+1]
				idx++
			}
		} else if IsState(s) && len(s) == 2 {
			a.State = s
		} else if IsZipcode(s) && a.State != "" {
			a.Zip, _ = strconv.Atoi(s)
		} else if (a.StreetDirection != "" || idx == 1) && a.StreetType == "" && a.Unit == "" && (idx > 0 && !IsStreetDirection(strings.Split(strings.TrimSpace(a.StreetName), " ")[len(strings.Split(strings.TrimSpace(a.StreetName), " "))-1])) && a.City == "" {
			// is street name
			a.StreetName += (s + " ")
		} else if (a.StreetType != "" || IsPoBox(address) || IsApartment(address) || (idx > 0 && IsStreetDirection(strings.Split(strings.TrimSpace(a.StreetName), " ")[len(strings.Split(strings.TrimSpace(a.StreetName), " "))-1]))) && a.State == "" && len(s) >= 2 {
			// is city
			a.City += (s + " ")
		} else if a.StreetType != "" && a.StreetName == "" {
			a.StreetName += s
		} else {
			// fmt.Println("DIDN'T MAKE IT:", s)
		}

	}
	a.StreetName = strings.TrimRight(a.StreetName, " ")
	a.City = strings.TrimRight(a.City, " ")

	return a, nil
}

// Extract attempts to extract an address from a string with surrounding words.
func Extract(in string) string {
	addressMinCharacters := 2
	s := strings.Replace(in, "  ", " ", -1)
	x := strings.Split(s, " ")

	for i := 0; i <= len(x); i++ {
		str := func(words []string, offset int) string {

			loops := len(words) - offset
			for i := 0; i < loops; i++ {
				// Enough words for an address
				if len(words[i:loops]) > addressMinCharacters {
					str := strings.Join(words[i:loops], " ")
					addr, _ := Parse(str)
					if addr.StreetType != "" && addr.HouseNumber != 0 {
						return str
					}
				}
			}

			return ""
		}(x, i)

		if str != "" {
			return str
		}
	}

	return ""
}

// Checks address string for indication of
// apt number.
func IsApartment(s string) bool {
	sX := strings.FieldsFunc(s, split)

	for _, s = range sX {
		s = strings.ToUpper(strings.TrimSpace(s))
		if s == "APT" || s == "#" || s == "UNIT" {
			return true
		}
	}

	return false
}

// Checks address string for indication of
// being a po box.
func IsPoBox(s string) bool {
	s = strings.NewReplacer(" ", "", ".", "").Replace(s)

	return strings.Contains(strings.ToUpper(s), "POBOX")
}

// Determines if string is a valid zip code
// or not. (trailing +4 is removed)
func IsZipcode(s string) bool {
	s = strings.Split(s, "-")[0]

	if isInt(s) && len(s) == 5 && s >= SMALLEST_ZIP_CODE && s <= LARGEST_ZIP_CODE {
		return true
	}

	return false
}

// Tries to match string with possible street types
// found in the U.S. (military excluded, I believe)
func IsStreetType(s string) bool {
	for _, value := range StreetTypes {
		if strings.ToUpper(s) == value {
			return true
		}
	}

	return false
}

// Tries to match string with possible street directions
// found in the U.S.
func IsStreetDirection(s string) bool {
	for _, value := range StreetDirections {
		if strings.ToUpper(strings.TrimSpace(s)) == value {
			return true
		}
	}

	return false
}

// Formats an address to string.
func (a *Address) String() string {
	var address string
	if a.StreetName == "PO Box" {
		address = a.StreetName + " " + strconv.Itoa(a.HouseNumber)
	} else {
		address = strconv.Itoa(a.HouseNumber) + " " + a.StreetDirection + " " + a.StreetName + " " + a.StreetType
		if a.Unit != "" {
			address += " # " + a.Unit
		}
	}
	address = strings.TrimSpace(address)

	if a.City != "" {
		address += " " + a.City
	}
	if a.State != "" {
		address += ", " + a.State
	}
	if a.Zip != 0 {
		address += " " + strconv.Itoa(a.Zip)
	}

	return strings.Replace(address, "  ", " ", -1)
}

// Formats an address to string (street only).
func (a *Address) Street() string {
	var address string
	if a.StreetName == "PO Box" {
		address = a.StreetName + " " + strconv.Itoa(a.HouseNumber)
	} else {
		address = strconv.Itoa(a.HouseNumber) + " " + a.StreetDirection + " " + a.StreetName + " " + a.StreetType
		if a.Unit != "" {
			address += " # " + a.Unit
		}
	}
	address = strings.TrimSpace(address)

	return strings.Replace(address, "  ", " ", -1)
}

// Checks string for apartment keyword.
func isApartmentKeyword(s string) bool {

	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "APT" || s == "#" || s == "UNIT" {
		return true
	}

	return false
}

// Determines if string is an integer.
func isInt(s string) bool {
	_, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return false
	}

	return true
}

// split func for strings.FieldsFunc splitting
// up an address into separate parts.
func split(r rune) bool {
	return r == ',' || r == ' '
}
