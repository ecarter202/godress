package godress

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Address represents a street address' parts.
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
	Latitude        float64 `arango:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude       float64 `arango:"longitude,omitempty" json:"longitude,omitempty"`
}

const (
	smallestZipCode = "01001"
	largestZipCode  = "99950"
	numberRegex     = `(?m)(\d+)`
)

// MustParse parses the address, ignoring any errors.
func MustParse(address string) (addr *Address) {
	addr, _ = Parse(address)

	return addr
}

// Parse parses a string into an address struct.
func Parse(address string) (*Address, error) {
	a := &Address{}
	a.Original = address
	newAddress := strings.NewReplacer("  ", " ", ".", "").Replace(regexp.MustCompile(`(?s)\(.*\)`).ReplaceAllString(address, ""))
	addressX := strings.FieldsFunc(newAddress, split)

	for idx := 0; idx < len(addressX); idx++ {

		s := addressX[idx]
		s = strings.TrimSpace(s)

		if idx == 0 && !IsPoBox(address) {
			if isInt(s) {
				a.HouseNumber, _ = strconv.Atoi(s)
			}
		} else if IsPoBox(address) && a.HouseNumber == 0 {
			re := regexp.MustCompile(numberRegex)
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
		} else if IsZipcode(s) && a.State != "" && a.Zip == 0 {
			a.Zip, _ = strconv.Atoi(strings.Split(s, "-")[0])
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

// IsApartment checks an address string for indication of apt number.
func IsApartment(s string) bool {
	sX := strings.FieldsFunc(s, split)

	for _, s = range sX {
		s = strings.ToLower(strings.TrimSpace(s))
		if _, ok := unitTerms[s]; ok {
			return true
		}
	}

	return false
}

// IsPoBox checks an address string for indication of being a POBox.
func IsPoBox(s string) bool {
	s = strings.NewReplacer(" ", "", ".", "").Replace(s)

	return strings.Contains(strings.ToUpper(s), "POBOX")
}

// IsZipcode determines if a string is a valid zip code
// or not. (trailing +4 is removed)
func IsZipcode(s string) bool {
	s = strings.Split(s, "-")[0]

	if isInt(s) && len(s) == 5 && s >= smallestZipCode && s <= largestZipCode {
		return true
	}

	return false
}

// String formats an address, returning it as a string.
func (a *Address) String() string {
	var address string
	if a.StreetName == "PO Box" {
		address = a.StreetName + " " + strconv.Itoa(a.HouseNumber)
	} else {
		if strings.ToUpper(a.State) == "WA" || strings.Index(a.Original, fmt.Sprintf(" %s", a.StreetDirection)) > strings.Index(a.Original, fmt.Sprintf(" %s", a.StreetType)) && a.StreetDirection != "" {
			address = fmt.Sprintf("%v %s %s %s", a.HouseNumber, a.StreetName, a.StreetType, a.StreetDirection)
		} else {
			address = fmt.Sprintf("%v %s %s %s", a.HouseNumber, a.StreetDirection, a.StreetName, a.StreetType)
		}
		if a.Unit != "" {
			address += " #" + a.Unit
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

// Street formats an address to string (street only).
func (a *Address) Street() string {
	var address string
	if a == nil {
		return ""
	} else if !strings.EqualFold(a.Original, "") {
		return strings.Split(a.Original, a.City)[0]
	} else if a.StreetName == "PO Box" {
		address = a.StreetName + " " + strconv.Itoa(a.HouseNumber)
	} else {
		if strings.ToUpper(a.State) == "WA" || strings.Index(a.Original, fmt.Sprintf(" %s", a.StreetDirection)) > strings.Index(a.Original, fmt.Sprintf(" %s", a.StreetType)) && a.StreetDirection != "" {
			address = fmt.Sprintf("%v %s %s %s", a.HouseNumber, a.StreetName, a.StreetType, a.StreetDirection)
		} else {
			address = fmt.Sprintf("%v %s %s %s", a.HouseNumber, a.StreetDirection, a.StreetName, a.StreetType)
		}
		if a.Unit != "" {
			address += " #" + a.Unit
		}
	}
	address = strings.TrimSpace(address)

	return strings.Replace(address, "  ", " ", -1)
}

func isApartmentKeyword(s string) bool {
	s = strings.ToUpper(strings.TrimSpace(s))
	if s == "APT" || s == "#" || s == "UNIT" {
		return true
	}

	return false
}

func isInt(s string) bool {
	_, err := strconv.Atoi(strings.TrimSpace(s))

	return err == nil
}

func split(r rune) bool {
	return r == ',' || r == ' '
}
