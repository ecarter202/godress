package godress

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Street struct {
	HouseNumber     int    `json:"house_number"`
	StreetDirection string `json:"street_direction"`
	StreetName      string `json:"street_name"`
	StreetType      string `json:"street_type"`
	Unit            string `json:"unit"`
}

func ParseStreet(street string) *Street {
	s := &Street{}

	if IsPoBox(street) {
		re := regexp.MustCompile("[0-9]+")
		matches := re.FindStringSubmatch(street)
		if len(matches) >= 1 {
			s.HouseNumber, _ = strconv.Atoi(strings.TrimSpace(matches[0]))
		}
		s.StreetName = "PO Box"
		return s
	} else {
		streetX := strings.Split(strings.TrimSpace(street), " ")
		s.HouseNumber, _ = strconv.Atoi(streetX[0])
		for idx, value := range streetX {
			if idx == 0 && s.HouseNumber != 0 {
				continue
			} else if IsStreetDirection(value) {
				s.StreetDirection = value
			} else if IsStreetType(value) {
				s.StreetType = value
			} else if isApartmentKeyword(value) && s.Unit == "" {
				if idx+1 < len(streetX) {
					s.Unit = streetX[idx+1]
					idx += 1
				}
			} else {
				s.StreetName += fmt.Sprintf("%s ", value)
			}
		}
	}

	s.StreetName = strings.TrimSpace(s.StreetName)
	return s
}

func (s *Street) String() string {
	if s.StreetName == "PO Box" {
		return fmt.Sprintf("PO Box %v", s.HouseNumber)
	} else {
		if s.Unit != "" {
			s.Unit = fmt.Sprintf("Unit %s", s.Unit)
		}
		return strings.TrimSpace(fmt.Sprintf("%v %s %s %s %v", s.HouseNumber, s.StreetDirection, s.StreetName, s.StreetType, s.Unit))
	}
}
