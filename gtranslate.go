package gtranslate

import (
	"golang.org/x/text/language"
	"net/http"
	"time"
)

var GoogleHost = "google.com"

// TranslationParams is a util struct to pass as parameter to indicate how to translate
type TranslationParams struct {
	From       string
	To         string
	Tries      int
	Delay      time.Duration
	GoogleHost string
	HTTPClient *http.Client
}

// Translate translate a text using native tags offer by go language
func Translate(text string, from language.Tag, to language.Tag, googleHost ...string) (string, error) {
	if len(googleHost) != 0 && googleHost[0] != "" {
		GoogleHost = googleHost[0]
	}
	translated, err := translate(http.DefaultClient, text, from.String(), to.String(), false, 2, 0)
	if err != nil {
		return "", err
	}

	return translated, nil
}

// TranslateWithParams translate a text with simple params as string
func TranslateWithParams(text string, params TranslationParams) (string, error) {
	if params.GoogleHost == "" {
		GoogleHost = "google.com"
	} else {
		GoogleHost = params.GoogleHost
	}
	if params.HTTPClient == nil {
		params.HTTPClient = http.DefaultClient
	}
	translated, err := translate(params.HTTPClient, text, params.From, params.To, true, params.Tries, params.Delay)
	if err != nil {
		return "", err
	}
	return translated, nil
}
