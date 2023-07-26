package elements

import "math"

func GetDiscount(price float32, previousPrice float32) int {
	if previousPrice != 0 {
		discountFloat := ((previousPrice - price) / previousPrice) * 100
		return int(math.Round(float64(discountFloat)))
	}
	return 0
}

