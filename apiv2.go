package gtranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/text/language"

	"github.com/robertkrimen/otto"
)

var ttk otto.Value

func init() {
	ttk, _ = otto.ToValue("0")
}

const (
	defaultNumberOfRetries = 2
)

func translate(httpClient *http.Client, text, from, to string, withVerification bool, tries int, delay time.Duration) (string, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if tries <= 0 {
		tries = defaultNumberOfRetries
	}

	if withVerification {
		if _, err := language.Parse(from); err != nil && from != "auto" {
			log.Println("[WARNING], '" + from + "' is a invalid language, switching to 'auto'")
			from = "auto"
		}
		if _, err := language.Parse(to); err != nil {
			log.Println("[WARNING], '" + to + "' is a invalid language, switching to 'en'")
			to = "en"
		}
	}

	t, _ := otto.ToValue(text)

	urll := fmt.Sprintf("https://translate.%s/translate_a/single", GoogleHost)

	token := get(t, ttk)

	data := map[string]string{
		"client": "gtx",
		"sl":     from,
		"tl":     to,
		"hl":     to,
		// "dt":     []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"},
		"ie":   "UTF-8",
		"oe":   "UTF-8",
		"otf":  "1",
		"ssel": "0",
		"tsel": "0",
		"kc":   "7",
		"q":    text,
	}

	u, err := url.Parse(urll)
	if err != nil {
		return "", nil
	}

	parameters := url.Values{}

	for k, v := range data {
		parameters.Add(k, v)
	}
	for _, v := range []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"} {
		parameters.Add("dt", v)
	}

	parameters.Add("tk", token)
	u.RawQuery = parameters.Encode()

	var r *http.Response

	for tries > 0 {
		r, err = httpClient.Get(u.String())
		if err != nil {
			if errors.Is(err, http.ErrHandlerTimeout) {
				return "", errBadNetwork
			}
		} else if r.StatusCode == http.StatusOK {
			break
		}

		if r.Body != nil {
			_ = r.Body.Close()
		}
		tries--
		time.Sleep(delay)
	}
	if err != nil {
		return "", err
	}
	defer func() { _ = r.Body.Close() }()
	if r.StatusCode != http.StatusOK {
		return "", errStatusCodeNotOK
	}

	var resp []interface{}
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return "", err
	}
	responseText := ""
	for _, obj := range resp[0].([]interface{}) {
		if len(obj.([]interface{})) == 0 {
			break
		}

		t, ok := obj.([]interface{})[0].(string)
		if ok {
			responseText += t
		}
	}

	return responseText, nil
}
