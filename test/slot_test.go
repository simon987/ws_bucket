package test

import (
	"github.com/simon987/ws_bucket/api"
	"testing"
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

func allocateUploadSlot(request api.AllocateUploadSlotRequest) (ar *api.GenericResponse) {
	resp := Post("/slot", request)
	UnmarshalResponse(resp, &ar)
	return
}
