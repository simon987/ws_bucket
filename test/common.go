package test

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/simon987/ws_bucket/api"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Post(path string, x interface{}) *http.Response {

	secret := os.Getenv("WS_BUCKET_SECRET")
	if secret == "" {
		secret = "default_secret"
	}

	s := http.Client{}

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://"+api.GetServerAddress()+path, buf)
	handleErr(err)

	ts := time.Now().Format(time.RFC1123)
	mac := hmac.New(crypto.SHA256.New, []byte(secret))
	mac.Write(body)
	mac.Write([]byte(ts))
	sig := hex.EncodeToString(mac.Sum(nil))
	req.Header.Add("X-Signature", sig)
	req.Header.Add("Timestamp", ts)

	r, err := s.Do(req)
	handleErr(err)

	return r
}

func Get(path string, token string) *http.Response {

	secret := os.Getenv("WS_BUCKET_SECRET")
	if secret == "" {
		secret = "default_secret"
	}

	s := http.Client{}

	req, err := http.NewRequest("GET", "http://"+api.GetServerAddress()+path, nil)
	handleErr(err)

	ts := time.Now().Format(time.RFC1123)
	mac := hmac.New(crypto.SHA256.New, []byte(secret))
	mac.Write([]byte(path))
	mac.Write([]byte(ts))
	sig := hex.EncodeToString(mac.Sum(nil))
	req.Header.Add("X-Signature", sig)
	req.Header.Add("Timestamp", ts)

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
