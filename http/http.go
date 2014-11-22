package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ReturnJson return json response with provided struct v
func ReturnJson(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.MarshalIndent(v, "", "\t")
	fmt.Fprintf(w, string(js))

	return err
}

func ReturnAcceptJsonOrText(w http.ResponseWriter, r *http.Request, plainText string, jsonResp interface{}) error {
	h := r.Header
	accept := h.Get("Accept")
	if strings.HasPrefix(accept, "application/json") {
		return ReturnJson(w, jsonResp)
	}

	fmt.Fprintf(w, plainText)
	return nil
}

// GetUrl fetch an url request with GET-Method. If a nil httpClient is
// provided, http.DefaultClient will be used.
// Returned the response as string
func GetUrl(client *http.Client, u string) string {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(u)
	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(body)
}

func GetJson(client *http.Client, url string, v interface{}) error {
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)

	return d.Decode(v)
}
