package postalcodes

import (
	"testing"
)


func TestGetDb(t *testing.T) {


	want := [][]string{
		{
			"MX", // COUNTRY_CODE
			"95756", // POSTAL_CODE
			"San Martin", // PLACE_NAME
			"Veracruz de Ignacio de la Llave", // ADMIN_NAME1
			"VER", // ADMIN_CODE1
			"San Andrés Tuxtla", // ADMIN_NAME2
			"141", // ADMIN_CODE2
			"San Andrés Tuxtla", // ADMIN_NAME3
			"15", // ADMIN_CODE3
			"19.6115", // LATITUDE 
			"-96.8888", // LONGITUDE
			"4", // ACCURACY
		}, {
			"MX", // COUNTRY_CODE
			"95760", // POSTAL_CODE
			"Las Rocas", // PLACE_NAME
			"Veracruz de Ignacio de la Llave", // ADMIN_NAME1
			"VER", // ADMIN_CODE1
			"San Andrés Tuxtla", // ADMIN_NAME2
			"141", // ADMIN_CODE2
			"San Andrés Tuxtla", // ADMIN_NAME3
			"15", // ADMIN_CODE3
			"19.6033", // LATITUDE 
			"-96.4011", // LONGITUDE
			"4", // ACCURACY
		},
	}

	// Test reading file
	got, _ := GetDb("XX")
	if len(got) != len(want) {
		t.Errorf("Error reading file. Expected len %d, but read %d", len(want), len(got))
	}
	for n, gotRec := range got {
		for i := 0; i < FIELDS_PER_RECORD; i++ {
			if gotRec[i] != want[n][i] {
				t.Errorf("Error reading file at line %d. Expected %s == %s", n, gotRec[i], want[n][i])
			}
		}
	}	
}

func TestPostalCodesNotFound(t *testing.T) {
	// Test reading file
	got, err := GetDb("YY")
	if got != nil {
		t.Errorf("Expected getting nil data")
	}
	if err == nil {
		t.Errorf("Expected error is %s", err)
	}
}