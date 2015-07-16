package postalcodes

import (
	"fmt"
	"bytes"
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


// LoadPostalCodes returns a postal codes database dump for given country code
// The dumps are taken from the geonames website http://download.geonames.org/export/zip/
func LoadPostalCodes(countryCode string) ([][]string, error) {

	data, err := Asset(fmt.Sprintf("data/%s.txt", countryCode))
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(data)

	tsv := csv.NewReader(reader)
	tsv.Comma 			 = '\t'
	tsv.FieldsPerRecord  = FIELDS_PER_RECORD
	tsv.LazyQuotes 		 = true
	tsv.TrailingComma 	 = true  // retain rather than remove empty slots
	tsv.TrimLeadingSpace = false // retain rather than remove empty slots

	return tsv.ReadAll()
}

