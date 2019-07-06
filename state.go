package godress

import "strings"

var (
	States = map[string]string{
		"alabama":        "al",
		"alaska":         "ak",
		"arizona":        "az",
		"arkansas":       "ar",
		"california":     "ca",
		"colorado":       "co",
		"connecticut":    "ct",
		"delaware":       "de",
		"florida":        "fl",
		"georgia":        "ga",
		"hawaii":         "hi",
		"idaho":          "id",
		"illinois":       "il",
		"indiana":        "in",
		"iowa":           "ia",
		"kansas":         "ks",
		"kentucky":       "ky",
		"louisiana":      "la",
		"maine":          "me",
		"maryland":       "md",
		"massachusetts":  "ma",
		"michigan":       "mi",
		"minnesota":      "mn",
		"mississippi":    "ms",
		"missouri":       "mo",
		"montana":        "mt",
		"nebraska":       "ne",
		"nevada":         "nv",
		"new hampshire":  "nh",
		"new jersey":     "nj",
		"new mexico":     "nm",
		"new york":       "ny",
		"north carolina": "nc",
		"north dakota":   "nd",
		"ohio":           "oh",
		"oklahoma":       "ok",
		"oregon":         "or",
		"pennsylvania":   "pa",
		"rhode island":   "ri",
		"south carolina": "sc",
		"south dakota":   "sd",
		"tennessee":      "tn",
		"texas":          "tx",
		"utah":           "ut",
		"vermont":        "vt",
		"virginia":       "va",
		"washington":     "wa",
		"west virginia":  "wv",
		"wisconsin":      "wi",
		"wyoming":        "wy",
	}
)

// IsState determines if the provided string is a valid state.
func IsState(s string) bool {
	if States[s] == "" {
		for _, v := range States {
			if strings.ToLower(s) == v {
				return true
			}
		}
		return false
	}

	return true
}

// StateAbbreviation gets the 2 letter abbreviation for a state.
func StateAbbreviation(s string) string {
	if len(s) > 2 {
		return strings.ToUpper(States[strings.ToLower(strings.TrimSpace(s))])
	}

	return strings.ToUpper(s)
}
