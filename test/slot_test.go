package test

import (
	"bytes"
	"github.com/fasthttp/websocket"
	"github.com/simon987/ws_bucket/api"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestAllocateUploadInvalidMaxSize(t *testing.T) {

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "valid",
		Token:    "valid",
		MaxSize:  -1,
	}).Ok != false {
		t.Error()
	}
}

func TestAllocateUploadSlotInvalidToken(t *testing.T) {

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "valid",
		Token:    "",
		MaxSize:  100,
	}).Ok != false {
		t.Error()
	}
}

func TestAllocateUploadSlotUnsafePath(t *testing.T) {

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "../test.png",
		Token:    "valid",
		MaxSize:  100,
	}).Ok != false {
		t.Error()
	}

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "test/../../test.png",
		Token:    "valid",
		MaxSize:  100,
	}).Ok != false {
		t.Error()
	}
}

func TestDuplicateUploadSlot(t *testing.T) {

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "test.png",
		Token:    "testdupe",
		MaxSize:  100,
	}).Ok != true {
		t.Error()
	}

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "test.png",
		Token:    "testdupe",
		MaxSize:  100,
	}).Ok != false {
		t.Error()
	}
}

func TestTokenInQueryString(t *testing.T) {

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "test.png",
		Token:    "testquery",
		MaxSize:  100,
	}).Ok != true {
		t.Error()
	}

	conn := ws("testquery")
	conn.WriteMessage(websocket.BinaryMessage, []byte("test"))
	conn.Close()

	time.Sleep(time.Millisecond * 20)
	r, err := http.Get("http://" + api.GetServerAddress() + "/slot?token=testquery")
	handleErr(err)

	data, err := ioutil.ReadAll(r.Body)
	handleErr(err)

	if bytes.Compare(data, []byte("test")) != 0 {
		t.Error()
	}

}

func TestUploadSlotInfo(t *testing.T) {

	if allocateUploadSlot(api.AllocateUploadSlotRequest{
		FileName: "testuploadslotinfo.png",
		Token:    "testuploadslotinfo",
		MaxSize:  123,
	}).Ok != true {
		t.Error()
	}

	resp := getSlotInfo("testuploadslotinfo")

	if resp.FileName != "testuploadslotinfo.png" {
		t.Error()
	}
	if resp.Token != "testuploadslotinfo" {
		t.Error()
	}
	if resp.MaxSize != 123 {
		t.Error()
	}
}

func allocateUploadSlot(request api.AllocateUploadSlotRequest) (ar *api.GenericResponse) {
	resp := Post("/slot", request)
	UnmarshalResponse(resp, &ar)
	return
}

func getSlotInfo(token string) (ar *api.GetUploadSlotResponse) {
	resp := Get("/slot_info", token)
	UnmarshalResponse(resp, &ar)
	return
}
