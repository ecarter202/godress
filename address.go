package godress

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	smallestZipCode    = "01001"
	largestZipCode     = "99950"
	numberRegexPattern = `(?m)(\d+)`
)

var (
	numberRegex = regexp.MustCompile(numberRegexPattern)
)

// Address represents a street address' parts.
type Address struct {
	Hash            string  `arango:"hash" json:"hash"`
	Original        string  `arango:"original" json:"original"`
	HouseNumber     string  `arango:"house_number" json:"house_number"`
	StreetDirection string  `arango:"street_direction" json:"street_direction"`
	StreetName      string  `arango:"street_name" json:"street_name"`
	StreetType      string  `arango:"street_type" json:"street_type"`
	Unit            string  `arango:"unit" json:"unit"`
	City            string  `arango:"city" json:"city"`
	County          string  `arango:"county" json:"county"`
	State           string  `arango:"state" json:"state"`
	PostalCode      string  `arango:"postal_code" json:"postal_code"`
	Country         string  `arango:"country" json:"country"`
	Latitude        float64 `arango:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude       float64 `arango:"longitude,omitempty" json:"longitude,omitempty"`
}

// MustParse parses the address, ignoring any errors.
func MustParse(address string) (addr *Address) {
	addr, _ = Parse(address)

	return addr
}

// Parse parses a string into an address struct.
func Parse(address string) (a *Address, err error) {
	stripped := regexp.MustCompile(`/\s+/ /`).ReplaceAllString(address, "")
	stripped = regexp.MustCompile(`\.`).ReplaceAllString(stripped, "")
	stripped = strings.ToUpper(stripped)

	a = &Address{Original: stripped}
	a.Hash = fmt.Sprintf("%x", md5.Sum([]byte(stripped)))

	x := strings.FieldsFunc(stripped, split)

	if IsPoBox(address) {
		a, err = parsePoBox(a, x)

		return
	}

	var (
		currentValue string
		cityWords    []string
	)

	for i := 0; i < len(x); i++ {
		currentValue = x[i]
		if i == 0 {
			if isInt(currentValue) {
				a.HouseNumber = currentValue
			}
		} else if IsStreetDirection(currentValue) && a.StreetDirection == "" {
			a.StreetDirection = currentValue
		} else if len(cityWords) == 0 && IsStreetType(currentValue) && a.Unit == "" {
			a.StreetType = currentValue
		} else if isApartmentKeyword(currentValue) && a.Unit == "" {
			if i+1 < len(x) {
				a.Unit = x[i+1]
				i++
			}
		} else if IsState(currentValue) {
			a.State = StateAbbreviation(currentValue)
		} else if IsZipcode(currentValue) && a.State != "" && a.PostalCode == "" {
			a.PostalCode = strings.Split(currentValue, "-")[0]
		} else if (a.StreetDirection != "" || i == 1) && a.StreetType == "" && a.Unit == "" && (i > 0 && !IsStreetDirection(strings.Split(strings.TrimSpace(a.StreetName), " ")[len(strings.Split(strings.TrimSpace(a.StreetName), " "))-1])) && a.City == "" {
			a.StreetName += (currentValue + " ")
		} else if a.State == "" && len(currentValue) >= 2 {
			cityWords = append(cityWords, currentValue)
		} else if a.StreetType != "" && a.StreetName == "" {
			a.StreetName += currentValue
		}
	}

	a.StreetName = strings.TrimRight(a.StreetName, " ")
	a.City = strings.Join(cityWords, " ")

	return
}

func parsePoBox(a *Address, x []string) (*Address, error) {
	var (
		currentValue string
		cityWords    []string
	)

	a.StreetName = "PO BOX"

	for i := 0; i < len(x); i++ {
		currentValue = x[i]
		if currentValue == "PO" || currentValue == "BOX" {
			continue
		}

		if a.HouseNumber == "" && isInt(currentValue) {
			a.HouseNumber = currentValue
		} else if IsState(currentValue) && len(currentValue) == 2 {
			a.State = currentValue
		} else if IsZipcode(currentValue) && a.State != "" && a.PostalCode == "" {
			a.PostalCode = strings.Split(currentValue, "-")[0]
		} else if len(currentValue) >= 2 {
			cityWords = append(cityWords, currentValue)
		}
	}
	a.City = strings.Join(cityWords, " ")

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
					if addr.StreetType != "" && addr.HouseNumber != "" {
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
		address = a.StreetName + " " + a.HouseNumber
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
	if a.PostalCode != "" {
		address += " " + a.PostalCode
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
		address = a.StreetName + " " + a.HouseNumber
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
