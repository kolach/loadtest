package random

import (
	"math/rand"
	"github.com/kolach/loadtest/postalcodes"
)


func PostalCodesRecord(countryCode string) ([]string, error) {
	db, err := postalcodes.GetDb(countryCode)
	if err != nil { return nil, err }
	line  := rand.Intn(len(db))	
	return db[line], nil
}

func Address(countryCode string) (string, error) {
	db, err := postalcodes.GetDb(countryCode)
	if err != nil { return "", err }
	line  := rand.Intn(len(db))	
	return db[line][postalcodes.PLACE_NAME], nil
}

func AddressGen(countryCode string) (<-chan string, error) {
	log.Debug("Loading address database")
	db, err := postalcodes.GetDb(countryCode)
	log.Debug("Database is loaded!")
	if err != nil { return nil, err }
	c := make(chan string)
	go func() {
		for {
			line := rand.Intn(len(db))
			c <- db[line][postalcodes.PLACE_NAME]
		}
	}()
	return c, nil
}