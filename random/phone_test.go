package random

import "testing"

func TestPhoneNumberMXIsRandom(t *testing.T) {

	NUMBERS_TO_GENERATE := 1000
	UNIQUE_NUMBERS 	    := 1000 // should be not less than

	phoneBook := make(map[string]string)
	for i := 0; i < NUMBERS_TO_GENERATE; i++ {
		phone, _ := PhoneNumber("MX", "55")
		phoneBook[phone] = phone
	}

	if len(phoneBook) < UNIQUE_NUMBERS {
		t.Errorf("Unable to generate %d of unique numbers. Generated: %d", UNIQUE_NUMBERS, len(phoneBook))
	}

}

func TestPhoneNumberError(t *testing.T) {
	want := "No generator function found for country code XX"
	_, err := PhoneNumber("XX", "55")
	if err == nil {
		t.Errorf("Error expected becasue countrycode is not supported")
	} else {
		if err.Error() != want {
			t.Errorf("Error wanted %q, but got %q", want, err)
		}
	}
}
