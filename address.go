package godress

import (
	"regexp"
	"strconv"
	"strings"
)

type Address struct {
	Original        string  `json:"original"`
	HouseNumber     int     `json:"house_number"`
	StreetDirection string  `json:"street_direction"`
	StreetName      string  `json:"street_name"`
	StreetType      string  `json:"street_type"`
	Unit            string  `json:"unit"`
	City            string  `json:"city"`
	State           string  `json:"state"`
	Zip             int     `json:"zip"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
}

const (
	SMALLEST_ZIP_CODE = "01001"
	LARGEST_ZIP_CODE  = "99950"
	NUMBER_REGEX      = `(?m)(\d+)`
)

var (
	states           = []string{"ALABAMA", "AL", "ALASKA", "AK", "ARIZONA", "AZ", "ARKANSAS", "AR", "CALIFORNIA", "CA", "COLORADO", "CO", "CONNECTICUT", "CT", "DELAWARE", "DE", "FLORIDA", "FL", "GEORGIA", "GA", "HAWAII", "HI", "IDAHO", "ID", "ILLINOIS", "IL", "INDIANA", "IN", "IOWA", "IA", "KANSAS", "KS", "KENTUCKY", "KY", "LOUISIANA", "LA", "MAINE", "ME", "MARYLAND", "MD", "MASSACHUSETTS", "MA", "MICHIGAN", "MI", "MINNESOTA", "MN", "MISSISSIPPI", "MS", "MISSOURI", "MO", "MONTANA", "MT", "NEBRASKA", "NE", "NEVADA", "NV", "NEW HAMPSHIRE", "NH", "NEW JERSEY", "NJ", "NEW MEXICO", "NM", "NEW YORK", "NY", "NORTH CAROLINA", "NC", "NORTH DAKOTA", "ND", "OHIO", "OH", "OKLAHOMA", "OK", "OREGON", "OR", "PENNSYLVANIA", "PA", "RHODE ISLAND", "RI", "SOUTH CAROLINA", "SC", "SOUTH DAKOTA", "SD", "TENNESSEE", "TN", "TEXAS", "TX", "UTAH", "UT", "VERMONT", "VT", "VIRGINIA", "VA", "WASHINGTON", "WA", "WEST VIRGINIA", "WV", "WISCONSIN", "WI", "WYOMING", "WY"}
	streetTypes      = []string{"ALY", "ANX", "AVE", "BLVD", "CIR", "CT", "CV", "CRES", "DR", "EXPY", "EXT", "GRV", "HWY", "HL", "KY", "LN", "LOOP", "MALL", "PARK", "PKWY", "PL", "PLZ", "PT", "RD", "ROW", "RUN", "SQ", "ST", "TER", "TRCE", "TRL", "WAY", "ZZ"}
	streetDirections = []string{"N", "NW", "NE", "S", "SW", "SE", "E", "W"}
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
				idx += 1
			}
		} else if IsState(s) && len(s) == 2 {
			a.State = s
		} else if IsZipcode(s) && a.State != "" {
			a.Zip, _ = strconv.Atoi(s)
		} else if (a.StreetDirection != "" || idx == 1) && a.StreetType == "" && a.Unit == "" {
			// is street name
			a.StreetName += (s + " ")
		} else if (a.StreetType != "" || IsPoBox(address) || IsApartment(address)) && a.State == "" && len(s) >= 2 {
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

// Attempts to extract an address from a string with surrounding words.
func Extract(in string) string {
	address_min_characters := 2
	s := strings.Replace(in, "  ", " ", -1)
	x := strings.Split(s, " ")

	for i := 0; i <= len(x); i++ {
		str := func(words []string, offset int) string {

			loops := len(words) - offset
			for i := 0; i < loops; i++ {
				// Enough words for an address
				if len(words[i:loops]) > address_min_characters {
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
	s = strings.Replace(s, " ", "", -1)

	return strings.Contains(strings.ToUpper(s), "POBOX")
}

// Tries to match string with possible U.S. state abbreviations
// and state names.
func IsState(s string) bool {
	for _, value := range states {
		if strings.ToUpper(s) == value {
			return true
		}
	}

	return false
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
	for _, value := range streetTypes {
		if strings.ToUpper(s) == value {
			return true
		}
	}

	return false
}

// Tries to match string with possible street directions
// found in the U.S.
func IsStreetDirection(s string) bool {
	for _, value := range streetDirections {
		if strings.ToUpper(s) == value {
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
