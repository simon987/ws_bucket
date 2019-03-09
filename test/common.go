package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simon987/ws_bucket/api"
	"io/ioutil"
	"net/http"
)

func Post(path string, x interface{}) *http.Response {

	s := http.Client{}

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://"+api.GetServerAddress()+path, buf)
	handleErr(err)

	//ts := time.Now().Format(time.RFC1123)
	//
	//mac := hmac.New(crypto.SHA256.New, worker.Secret)
	//mac.Write(body)
	//mac.Write([]byte(ts))
	//sig := hex.EncodeToString(mac.Sum(nil))
	//
	//req.Header.Add("X-Worker-Id", strconv.FormatInt(worker.Id, 10))
	//req.Header.Add("X-Signature", sig)
	//req.Header.Add("Timestamp", ts)

	r, err := s.Do(req)
	handleErr(err)

	return r
}

func Get(path string, token string) *http.Response {

	s := http.Client{}

	req, err := http.NewRequest("GET", "http://"+api.GetServerAddress()+path, nil)
	handleErr(err)

	req.Header.Set("X-Upload-Token", token)

	r, err := s.Do(req)
	return r
}

func UnmarshalResponse(r *http.Response, result interface{}) {
	data, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(data))
	err = json.Unmarshal(data, result)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
