package test

import (
	"bytes"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/simon987/ws_bucket/api"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestWebSocketUploadSmallFile(t *testing.T) {

	id := uuid.New()

	allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "testfile",
		Token:    id.String(),
		MaxSize:  math.MaxInt64,
	})

	c := ws(id.String())

	err := c.WriteMessage(websocket.BinaryMessage, []byte("testuploadsmallfile"))
	handleErr(err)

	err = c.Close()
	handleErr(err)

	resp := readUploadSlot(id.String())

	if bytes.Compare(resp, []byte("testuploadsmallfile")) != 0 {
		t.Error()
	}
}

func TestWebSocketUploadOverwritesFile(t *testing.T) {

	id := uuid.New()

	allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "testuploadoverwrites",
		Token:    id.String(),
		MaxSize:  math.MaxInt64,
	})

	c := ws(id.String())

	err := c.WriteMessage(websocket.BinaryMessage, []byte("testuploadsmallfile"))
	handleErr(err)

	err = c.Close()
	handleErr(err)

	time.Sleep(time.Millisecond * 50)
	resp := readUploadSlot(id.String())

	if bytes.Compare(resp, []byte("testuploadsmallfile")) != 0 {
		t.Error()
	}

	c1 := ws(id.String())

	err = c1.WriteMessage(websocket.BinaryMessage, []byte("newvalue"))
	handleErr(err)

	err = c1.Close()
	handleErr(err)

	time.Sleep(time.Millisecond * 50)

	resp = readUploadSlot(id.String())

	if bytes.Compare(resp, []byte("newvalue")) != 0 {
		t.Error()
	}
}

func TestWebSocketUploadLargeFile(t *testing.T) {

	id := uuid.New()

	allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "testlargefile",
		Token:    id.String(),
		MaxSize:  math.MaxInt64,
	})

	c := ws(id.String())

	chunk := make([]byte, 100000)
	_ = copy(chunk, "test")
	_ = c.WriteMessage(websocket.BinaryMessage, chunk)

	err := c.Close()
	handleErr(err)

	resp := readUploadSlot(id.String())

	if bytes.Compare(resp, chunk) != 0 {
		t.Error()
	}
}

func TestWebSocketUploadMaxSize(t *testing.T) {

	id := uuid.New()

	allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "testmaxsize",
		Token:    id.String(),
		MaxSize:  10,
	})

	c := ws(id.String())

	chunk := make([]byte, 100000)
	_ = copy(chunk, "test")
	_ = c.WriteMessage(websocket.BinaryMessage, chunk)

	err := c.Close()
	handleErr(err)

	resp := readUploadSlot(id.String())

	if len(resp) != 10 {
		t.Error()
	}
}

func readUploadSlot(token string) []byte {

	time.Sleep(time.Millisecond * 20)
	r := Get("/slot", token)

	data, err := ioutil.ReadAll(r.Body)
	handleErr(err)

	return data
}

func ws(slot string) *websocket.Conn {

	u := url.URL{Scheme: "ws", Host: "localhost:3021", Path: "/upload"}
	fmt.Printf("Connecting to %s", u.String())

	header := http.Header{}
	header.Add("X-Upload-Token", slot)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	handleErr(err)

	c.EnableWriteCompression(true)

	return c
}
