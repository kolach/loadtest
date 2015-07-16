package db

import (
	"fmt"
	"os"
	"path"
	"encoding/csv"
)

const (
	COUNTRY_CODE int = iota
	POSTAL_CODE  
	PLACE_NAME   
	ADMIN_NAME1  
	ADMIN_CODE1  
	ADMIN_NAME2  
	ADMIN_CODE2  
	ADMIN_NAME3  
	ADMIN_CODE3  
	LATITUDE     
	LONGITUDE    
	ACCURACY     
	FIELDS_PER_RECORD
)

// postalCodesFilename returns filename for given country code
// according to internal nameing conventions
func postalCodesFilename(countryCode string) string {
	return path.Join("data", "postal_codes", fmt.Sprintf("%s.txt", countryCode))
}

// LoadPostalCodes returns a postal codes database dump for given country code
// The dumps are taken from the geonames website http://download.geonames.org/export/zip/
func LoadPostalCodes(countryCode string) ([][]string, error) {
	f, err := os.Open(postalCodesFilename(countryCode))
	if err != nil {
		return nil, fmt.Errorf("LoadPostalCodes: node postal code found for %s", countryCode)
	}

	defer f.Close()

	tsv := csv.NewReader(f)
	tsv.Comma 			 = '\t'
	tsv.FieldsPerRecord  = FIELDS_PER_RECORD
	tsv.LazyQuotes 		 = true
	tsv.TrailingComma 	 = true  // retain rather than remove empty slots
	tsv.TrimLeadingSpace = false // retain rather than remove empty slots

	return tsv.ReadAll()
}

