package random

import (
	"fmt"
	"math/rand"	
)

type gen_func func(string) string
type gen_map map[string]gen_func

var generators gen_map

func PhoneNumber(countryCode, cityCode string) (string, error) {
	gen := generators[countryCode]
	if gen == nil {
		return "", fmt.Errorf("No generator function found for country code %s", countryCode)
	} else {
		return gen(cityCode), nil
	}
}

// mexican phone number random generator
func mx(cityCode string) string {	
	return fmt.Sprintf(
		"+521%s%d%d%d%d%d%d%d%d",
		cityCode,
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
		rand.Intn(9),
	)
}

func init() {
	generators = make(gen_map)
	generators["MX"] = mx
}