package elements

import (
	"regexp"
	"strconv"
	"strings"
)

func GetReviews(eText string) int {
	regex, _ := regexp.Compile(`[\d\.]+`)
	match := regex.FindString(eText)
	removedDots := strings.ReplaceAll(match, ".", "")
	reviews, err := strconv.Atoi(removedDots)
	if err != nil {
		return 0
	}
	return reviews
}