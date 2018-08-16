package voucher

import (
	"testing"
)

var testExpectedMap = map[string]bool{
	"a": true,
	"b": false,
	"c": false,
	"f": true,
}

var testGoodMap = map[string]interface{}{
	"a": true,
	"b": false,
	"c": false,
	"e": 55,
	"f": true,
}

func TestToMapStringBool(t *testing.T) {
	convert := ToMapStringBool(testGoodMap)

	for key, value := range convert {
		if testExpectedMap[key] != value {
			t.Errorf("Value for key %s is not %v (should be %v)", key, value, testExpectedMap[key])
		}
	}
}
