package godress

import "strings"

// Term represents a company name term.
type Term struct {
	Abbreviation string
	Label        string
}

var (
	unitTerms = map[string]*Term{
		"apt":   &Term{"Apt", "Apartment"},
		"ste":   &Term{"Ste", "Suite"},
		"suite": &Term{"Suite", "Suite"},
		"unit":  &Term{"unit", "Unit"},
		"#":     &Term{"#", "Number"},
	}
)

// Scrub will remove any unit term from the addess.
func ScrubUnit(address string) string {
	addressX := strings.Split(strings.TrimSpace(address), " ")
	for i, w := range addressX {
		if _, contains := unitTerms[strings.ToLower(w)]; contains {
			return ScrubUnit(strings.Join(append(addressX[:i], addressX[i+1:]...), " "))
		}
	}

	return strings.Join(addressX, " ")
}
