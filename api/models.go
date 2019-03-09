package api

import (
	"fmt"
	"path/filepath"
	"strings"
)

type GenericResponse struct {
	Ok bool `json:"ok"`
}

type AllocateUploadSlotRequest struct {
	Token    string `json:"token"`
	MaxSize  int64  `json:"max_size"`
	FileName string `json:"file_name"`
}

func (req *AllocateUploadSlotRequest) IsValid() bool {
	if len(req.Token) < 3 {
		return false
	}

	if len(req.FileName) <= 0 {
		return false
	}

	path := filepath.Join(WorkDir, req.FileName)
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	if !strings.HasPrefix(pathAbs, WorkDir) {
		fmt.Println(pathAbs)
		fmt.Println(WorkDir)
		return false
	}

	if req.MaxSize < 0 {
		return false
	}

	return true
}

type UploadSlot struct {
	MaxSize  int64  `json:"max_size"`
	Token    string `gorm:"primary_key",json:"token"`
	FileName string `json:"file_name"`
}

type WebsocketMotd struct {
	Info Info   `json:"info"`
	Motd string `json:"motd"`
}
