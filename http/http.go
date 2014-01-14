package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpReturnJson(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")

	js, _ := json.MarshalIndent(v, "", "\t")
	fmt.Fprintf(w, string(js))
}

func HttpGetUrl(client *http.Client, u string) string {
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
