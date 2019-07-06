package monobank

import "strconv"

// ToBanknote - converts coins(i) into banknotes with decimal.
//
// `minorUnits` - symbols after dot.
//
// TODO: write benchmarks for ToBanknote() and next:
//
//     rate := math.Pow(10, float64(minorUnits))
//     return float64(i) / rate
func ToBanknote(i int64, minorUnits int) string {
	s := strconv.FormatInt(i, 10)

	// indent:
	for len(s) <= minorUnits {
		s = "0" + s
	}

	return s[:len(s)-minorUnits] + "." + s[len(s)-minorUnits:]
}
