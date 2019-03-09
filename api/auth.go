package api

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"github.com/valyala/fasthttp"
	"math"
	"os"
	"time"
)

var Secret = []byte(getApiSecret())

func getApiSecret() string {

	secret := os.Getenv("WS_BUCKET_SECRET")
	if secret == "" {
		return "default_secret"
	} else {
		return secret
	}
}

func validateRequest(ctx *fasthttp.RequestCtx) error {

	signature := ctx.Request.Header.Peek("X-Signature")
	timeStampStr := string(ctx.Request.Header.Peek("Timestamp"))

	if timeStampStr == "" {
		return errors.New("date is not specified")
	}

	timestamp, err := time.Parse(time.RFC1123, timeStampStr)
	if err != nil {
		return err
	}

	if math.Abs(float64(timestamp.Unix()-time.Now().Unix())) > 60 {
		return errors.New("invalid Timestamp")
	}

	var body []byte
	if ctx.Request.Header.IsGet() {
		body = ctx.Request.RequestURI()
	} else {
		body = ctx.Request.Body()
	}

	mac := hmac.New(crypto.SHA256.New, Secret)
	mac.Write(body)
	mac.Write([]byte(timeStampStr))

	expectedMac := make([]byte, 64)
	hex.Encode(expectedMac, mac.Sum(nil))
	matches := bytes.Compare(expectedMac, signature) == 0

	if !matches {
		return errors.New("signature does not match")
	}
	return nil
}
