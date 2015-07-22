package postalcodes

import (
	"fmt"
	"bytes"
	"encoding/csv"
	"github.com/op/go-logging"
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

var log = logging.MustGetLogger("random")

type PostalCodesDb [][]string
type PostalCodesMap map[string]PostalCodesDb
var postalCodesMap = make(PostalCodesMap)

// GetDb returns a postal codes database dump for given country code
// The dumps are taken from the geonames website http://download.geonames.org/export/zip/
func GetDb(countryCode string) ([][]string, error) {

	db := postalCodesMap[countryCode]
	if db == nil {

		log.Debug("Loading postalcodes database for country: %s", countryCode)

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


		db, err = tsv.ReadAll()
		if err != nil {
			return nil, err
		}

		postalCodesMap[countryCode] = db

		log.Debug("Database is loaded for countrycode: %s", countryCode)
	}

	return db, nil
}

