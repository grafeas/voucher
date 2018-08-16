package voucher

// ToMapStringBool takes a map[string]interface{} and converts it to a
// map[string]bool (dropping any values that do not cast to booleans
// cleanly).
func ToMapStringBool(in map[string]interface{}) (out map[string]bool) {
	out = make(map[string]bool, len(in))
	for key, rawValue := range in {
		if value, ok := rawValue.(bool); ok {
			out[key] = value
		}
	}
	return
}
