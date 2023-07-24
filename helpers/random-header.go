package helpers

import (
	"math/rand"
	"strings"
	"github.com/EDDYCJY/fake-useragent"
)

func isInArray(arr []string, element string) bool {
	for _, item := range arr {
		if item == element {
			return true
		}
	}
	return false
}

func getRandomNumber(min, max int) int {
	// Generate a random number within the range (min <= number < max)
	return rand.Intn(max-min) + min
}

func UpgradeInsecureRequests() string {
	return "1"

}

func accept() string {
	mediaTypes := []string{
		"application/json",
		"application/xml",
		"text/plain",
		"image/jpeg",
		"image/png",
		"audio/mpeg",
		"video/mp4",
		"application/pdf",
		"application/octet-stream",
		"text/css",
		"application/javascript",
	}

	var accept []string
	accept = append(accept, "text/html")

	for i := 0; i < 5; i++ {
		randomIndex := rand.Intn(len(mediaTypes))
		if !isInArray(accept, mediaTypes[randomIndex]) {
			accept = append(accept, mediaTypes[randomIndex])
		}
	}

	return strings.Join(accept, ",")
}

func acceptLanguage() string {
	ptBRVariations := []string{
		"pt-BR",
		"pt-BR,pt;q=0.9",
		"pt-BR;q=0.8,en-US;q=0.5",
		"pt-BR,en-US;q=0.7,en;q=0.3",
		"pt-BR;q=0.6,en-US;q=0.4,fr;q=0.2",
	}

	var languages []string
	var max = getRandomNumber(2, 5)
	for i := 0; i < max; i++ {
		randomIndex := rand.Intn(len(ptBRVariations))
		if !isInArray(languages, ptBRVariations[randomIndex]) {
			languages = append(languages, ptBRVariations[randomIndex])
		}
	}

	return strings.Join(languages, ",")
}

// Build creates an map representing a random request header
//
// Returns:
// A map containing the header keys and values
func RandomHeader() map[string]string {
	fakeHeader := make(map[string]string)
	fakeHeader["Upgrade-Insecure-Requests"] = UpgradeInsecureRequests()
	fakeHeader["Accept"] = accept()
	fakeHeader["Accept-Language"] = acceptLanguage()
	fakeHeader["Accept-Encoding"] = "*"
	fakeHeader["Connection"] = "keep-alive"
	fakeHeader["Sec-Fetch-Dest"] = "document"
	fakeHeader["Sec-Fetch-Mode"] = "navigate"
	fakeHeader["Sec-Fetch-Site"] = "none"
	fakeHeader["Sec-Fetch-User"] = "?1"
	fakeHeader["User-Agent"] = browser.Computer()
	fakeHeader["Referer"] = "https =//google.com/"

	return fakeHeader
}