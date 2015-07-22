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

func PhoneNumberGen(countryCode, cityCode string) (<-chan string, error) {
	f := generators[countryCode]
	if f == nil {
		return nil, fmt.Errorf("No generator function found for country code %s", countryCode)
	}
	c := make(chan string)
	go func() {
		for {
			c <- f(cityCode)
		}
	}()
	return c, nil
}

func init() {
	generators = make(gen_map)
	generators["MX"] = mx
}