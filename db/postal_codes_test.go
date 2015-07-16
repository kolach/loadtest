package db

import (
	"testing"
	"os"
	"path"
	"encoding/csv"
)

func TestPostalCodesFilename(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"MX", "data/postal_codes/MX.txt"},
		{"BR", "data/postal_codes/BR.txt"},
	}

	for _, c := range cases {
		got := postalCodesFilename(c.in)
		if got != c.want {
			t.Errorf("postalCodesFilename(%q) == %q, want %q", c.in, got, c.want)
		}
	}	
}

func TestLoadPostalCodes(t *testing.T) {

	// Before
	tmpFile := path.Join("data", "postal_codes", "XX.txt")
	f, err := os.Create(tmpFile)	
	if err != nil { t.Errorf("Unable to create temporal file") }	

	defer func() { os.Remove(tmpFile) }()

	tsv := csv.NewWriter(f)
	tsv.Comma = '\t'
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
	tsv.WriteAll(want)
	tsv.Flush()	

	// Test reading file
	got, _ := LoadPostalCodes("XX")
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
	got, err := LoadPostalCodes("YY")
	if got != nil {
		t.Errorf("Expected getting nil data")
	}
	if err == nil {
		t.Errorf("Expected error is %s", err)
	}
}