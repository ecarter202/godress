package godress

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	streetTypesByFull = map[string]string{"alley": "aly", "annex": "anx", "arcade": "arc", "avenue": "ave", "bayou": "yu", "beach": "bch", "bend": "bnd", "bluff": "blf", "bottom": "btm", "boulevard": "blvd", "branch": "br", "bridge": "brg", "brook": "brk", "burg": "bg", "bypass": "byp", "camp": "cp", "canyon": "cyn", "cape": "cpe", "causeway": "cswy", "center": "ctr", "circle": "cir", "cliffs": "clfs", "club": "clb", "corner": "cor", "corners": "cors", "course": "crse", "court": "ct", "courts": "cts", "cove": "cv", "creek": "crk", "crescent": "cres", "crossing": "xing", "dale": "dl", "dam": "dm", "divide": "dv", "drive": "dr", "estates": "est", "expressway": "expy", "extension": "ext", "fall": "fall", "falls": "fls", "ferry": "fry", "field": "fld", "fields": "flds", "flats": "flt", "ford": "for", "forest": "frst", "forge": "fgr", "fork": "fork", "forks": "frks", "fort": "ft", "freeway": "fwy", "gardens": "gdns", "gateway": "gtwy", "glen": "gln", "green": "gn", "grove": "grv", "harbor": "hbr", "haven": "hvn", "heights": "hts", "highway": "hwy", "hill": "hl", "hills": "hls", "hollow": "holw", "inlet": "inlt", "island": "is", "islands": "iss", "isle": "isle", "junction": "jct", "key": "cy", "knolls": "knls", "lake": "lk", "lakes": "lks", "landing": "lndg", "lane": "ln", "light": "lgt", "loaf": "lf", "locks": "lcks", "lodge": "ldg", "loop": "loop", "mall": "mall", "manor": "mnr", "meadows": "mdws", "mill": "ml", "mills": "mls", "mission": "msn", "mount": "mt", "mountain": "mtn", "neck": "nck", "orchard": "orch", "oval": "oval", "park": "park", "parkway": "pky", "pass": "pass", "path": "path", "pike": "pike", "pines": "pnes", "place": "pl", "plain": "pln", "plains": "plns", "plaza": "plz", "point": "pt", "port": "prt", "prairie": "pr", "radial": "radl", "ranch": "rnch", "rapids": "rpds", "rest": "rst", "ridge": "rdg", "river": "riv", "road": "rd", "row": "row", "run": "run", "shoal": "shl", "shoals": "shls", "shore": "shr", "shores": "shrs", "spring": "spg", "springs": "spgs", "spur": "spur", "square": "sq", "station": "sta", "stravenues": "stra", "stream": "strm", "street": "st", "summit": "smt", "terrace": "ter", "trace": "trce", "track": "trak", "trail": "trl", "trailer": "trlr", "tunnel": "tunl", "turnpike": "tpke", "union": "un", "valley": "vly", "viaduct": "via", "view": "vw", "village": "vlg", "ville": "vl", "vista": "vis", "walk": "walk", "way": "way", "wells": "wls"}
	streetTypesByAbbr = map[string]string{"aly": "alley", "anx": "annex", "arc": "arcade", "ave": "avenue", "yu": "bayou", "bch": "beach", "bnd": "bend", "blf": "bluff", "btm": "bottom", "blvd": "boulevard", "br": "branch", "brg": "bridge", "brk": "brook", "bg": "burg", "byp": "bypass", "cp": "camp", "cyn": "canyon", "cpe": "cape", "cswy": "causeway", "ctr": "center", "cir": "circle", "clfs": "cliffs", "clb": "club", "cor": "corner", "cors": "corners", "crse": "course", "ct": "court", "cts": "courts", "cv": "cove", "crk": "creek", "cres": "crescent", "xing": "crossing", "dl": "dale", "dm": "dam", "dv": "divide", "dr": "drive", "est": "estates", "expy": "expressway", "ext": "extension", "fall": "fall", "fls": "falls", "fry": "ferry", "fld": "field", "flds": "fields", "flt": "flats", "for": "ford", "frst": "forest", "fgr": "forge", "fork": "fork", "frks": "forks", "ft": "fort", "fwy": "freeway", "gdns": "gardens", "gtwy": "gateway", "gln": "glen", "gn": "green", "grv": "grove", "hbr": "harbor", "hvn": "haven", "hts": "heights", "hwy": "highway", "hl": "hill", "hls": "hills", "holw": "hollow", "inlt": "inlet", "is": "island", "iss": "islands", "isle": "isle", "jct": "junction", "cy": "key", "knls": "knolls", "lk": "lake", "lks": "lakes", "lndg": "landing", "ln": "lane", "lgt": "light", "lf": "loaf", "lcks": "locks", "ldg": "lodge", "loop": "loop", "mall": "mall", "mnr": "manor", "mdws": "meadows", "ml": "mill", "mls": "mills", "msn": "mission", "mt": "mount", "mtn": "mountain", "nck": "neck", "orch": "orchard", "oval": "oval", "park": "park", "pky": "parkway", "pass": "pass", "path": "path", "pike": "pike", "pnes": "pines", "pl": "place", "pln": "plain", "plns": "plains", "plz": "plaza", "pt": "point", "prt": "port", "pr": "prairie", "radl": "radial", "rnch": "ranch", "rpds": "rapids", "rst": "rest", "rdg": "ridge", "riv": "river", "rd": "road", "row": "row", "run": "run", "shl": "shoal", "shls": "shoals", "shr": "shore", "shrs": "shores", "spg": "spring", "spgs": "springs", "spur": "spur", "sq": "square", "sta": "station", "stra": "stravenues", "strm": "stream", "st": "street", "smt": "summit", "ter": "terrace", "trce": "trace", "trak": "track", "trl": "trail", "trlr": "trailer", "tunl": "tunnel", "tpke": "turnpike", "un": "union", "vly": "valley", "via": "viaduct", "vw": "view", "vlg": "village", "vl": "ville", "vis": "vista", "walk": "walk", "way": "way", "wls": "wells"}
	streetDirections  = []string{"N", "NW", "NE", "S", "SW", "SE", "E", "W"}
)

// Street represents a street, as in a part of a street address.
type Street struct {
	HouseNumber     int    `json:"house_number"`
	StreetDirection string `json:"street_direction"`
	StreetName      string `json:"street_name"`
	StreetType      string `json:"street_type"`
	Unit            string `json:"unit"`
}

// ParseStreet atempts to parse a string into the parts of a street.
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

// SetStreet will set an addresses street values from a parsed street.
func (a *Address) SetStreet(street *Street) {
	a.HouseNumber = street.HouseNumber
	a.StreetName = street.StreetName
	a.StreetDirection = street.StreetDirection
	a.StreetType = street.StreetType
	a.Unit = street.Unit
}

// String will return a parsed street as a string.
func (s *Street) String() string {
	if s.StreetName == "PO Box" {
		return fmt.Sprintf("PO Box %v", s.HouseNumber)
	}

	if s.Unit != "" {
		s.Unit = fmt.Sprintf("Unit %s", s.Unit)
	}

	return strings.TrimSpace(fmt.Sprintf("%v %s %s %s %v", s.HouseNumber, s.StreetDirection, s.StreetName, s.StreetType, s.Unit))
}

// IsStreetType attempts to match string with possible street types
// found in the U.S. (military excluded, I believe)
func IsStreetType(s string) bool {
	s = strings.ToLower(s)

	_, ok := streetTypesByAbbr[s]

	return ok
}

// Tries to match string with possible street directions
// found in the U.S.
func IsStreetDirection(s string) bool {
	for _, value := range streetDirections {
		if strings.ToUpper(strings.TrimSpace(s)) == value {
			return true
		}
	}

	return false
}

// StreetTypeAbbr takes the full name of a street type i.e. Circle
// and returns the abbreviation for it i.e. Cir
// If no match is found, the supplied string is returned.
func StreetTypeAbbr(full string) (abbr string) {
	var ok bool
	abbr, ok = streetTypesByFull[strings.ToLower(full)]
	if !ok {
		return full
	}

	return strings.Title(abbr)
}
