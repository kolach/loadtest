package random

import (
	"testing"
)

func TestAddressIsRandom(t *testing.T) {

	ADDRESSES_TO_GENERATE  := 10
	UNIQUE_ADDRESSES 	   := 10 // should be not less than

	m := make(map[string]string)
	for i := 0; i < ADDRESSES_TO_GENERATE; i++ {
		a, _ := Address("MX")
		m[a] = a
	}

	if len(m) < UNIQUE_ADDRESSES {
		t.Errorf("Unable to generate %d of unique addresses. Generated: %d", UNIQUE_ADDRESSES, len(m))
	}

}