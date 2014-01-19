package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ReturnJson return json response with provided struct v
func ReturnJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	js, _ := json.MarshalIndent(v, "", "\t")
	fmt.Fprintf(w, string(js))
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
