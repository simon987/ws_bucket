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
)

func TestWebsocketReturnsMotd(t *testing.T) {

	id := uuid.New()
	allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "testmotd",
		MaxSize:  0,
		Token:    id.String(),
	})

	c := ws(id.String())
	motd := &api.WebsocketMotd{}
	err := c.ReadJSON(&motd)
	handleErr(err)

	if len(motd.Motd) <= 0 {
		t.Error()
	}
	if len(motd.Info.Version) <= 0 {
		t.Error()
	}
}

func TestWebSocketUploadSmallFile(t *testing.T) {

	id := uuid.New()

	allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "testfile",
		Token:    id.String(),
		MaxSize:  math.MaxInt64,
	})

	c := ws(id.String())
	_, _, err := c.ReadMessage()
	handleErr(err)

	err = c.WriteMessage(websocket.BinaryMessage, []byte("testuploadsmallfile"))
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
	_, _, err := c.ReadMessage()
	handleErr(err)

	err = c.WriteMessage(websocket.BinaryMessage, []byte("testuploadsmallfile"))
	handleErr(err)

	err = c.Close()
	handleErr(err)

	c1 := ws(id.String())
	_, _, err = c1.ReadMessage()
	handleErr(err)

	err = c1.WriteMessage(websocket.BinaryMessage, []byte("newvalue"))
	handleErr(err)

	err = c1.Close()
	handleErr(err)

	resp := readUploadSlot(id.String())

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
	_, _, err := c.ReadMessage()
	handleErr(err)

	chunk := make([]byte, 100000)
	_ = copy(chunk, "test")
	_ = c.WriteMessage(websocket.BinaryMessage, chunk)

	err = c.Close()
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
	_, _, err := c.ReadMessage()
	handleErr(err)

	chunk := make([]byte, 100000)
	_ = copy(chunk, "test")
	_ = c.WriteMessage(websocket.BinaryMessage, chunk)

	err = c.Close()
	handleErr(err)

	resp := readUploadSlot(id.String())

	if len(resp) != 10 {
		t.Error()
	}
}

func readUploadSlot(token string) []byte {

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

	return c
}
