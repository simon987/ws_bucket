package api

import (
	"encoding/json"
	"github.com/fasthttp/websocket"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const WsBufferSize = 4096

var Mutexes sync.Map
var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:    WsBufferSize,
	WriteBufferSize:   WsBufferSize,
	EnableCompression: true,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

func (api *WebApi) AllocateUploadSlot(ctx *fasthttp.RequestCtx) {

	err := validateRequest(ctx)
	if err != nil {
		ctx.Response.Header.SetStatusCode(401)
		Json(GenericResponse{
			Ok:      false,
			Message: err.Error(),
		}, ctx)
		return
	}

	req := &AllocateUploadSlotRequest{}
	err = json.Unmarshal(ctx.Request.Body(), req)
	if err != nil {
		ctx.Response.Header.SetStatusCode(400)
		Json(GenericResponse{
			Ok: false,
		}, ctx)
		return
	}

	if !req.IsValid() {
		ctx.Response.Header.SetStatusCode(400)
		Json(GenericResponse{
			Ok: false,
		}, ctx)
		return
	}

	err = api.allocateUploadSlot(req)

	if err == nil {
		Json(GenericResponse{
			Ok: true,
		}, ctx)
	} else {
		Json(GenericResponse{
			Ok:      false,
			Message: err.Error(),
		}, ctx)
	}
}

func (api *WebApi) Upload(ctx *fasthttp.RequestCtx) {

	token := string(ctx.Request.Header.Peek("X-Upload-Token"))
	slot := UploadSlot{}
	err := api.db.Where("token=?", token).First(&slot).Error
	if err != nil {
		ctx.Response.Header.SetStatusCode(400)
		logrus.WithError(err).WithFields(logrus.Fields{
			"token": token,
		}).Warning("Upload slot not found")
		return
	}

	logrus.WithFields(logrus.Fields{
		"slot": slot,
	}).Info("Upgrading connection")

	err = upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer ws.Close()

		mt, reader, err := ws.NextReader()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"token": token,
			}).Warning("Client disconnected before sending the first byte")
			return
		}
		if mt != websocket.BinaryMessage {
			return
		}

		mu, _ := Mutexes.LoadOrStore(slot.Token, &sync.RWMutex{})
		mu.(*sync.RWMutex).Lock()
		path := filepath.Join(WorkDir, slot.FileName)
		fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			logrus.WithError(err).Error("Error while opening file for writing")
		}

		buf := make([]byte, WsBufferSize)
		totalRead := int64(0)
		for totalRead < slot.MaxSize {
			read, err := reader.Read(buf)

			var toWrite int
			if totalRead+int64(read) > slot.MaxSize {
				toWrite = int(slot.MaxSize - totalRead)
			} else {
				toWrite = read
			}

			_, _ = fp.Write(buf[:toWrite])
			if err == io.EOF {
				break
			}
			totalRead += int64(read)
		}

		logrus.WithFields(logrus.Fields{
			"totalRead": totalRead,
		}).Info("Finished reading")
		err = fp.Close()
		if err != nil {
			logrus.WithError(err).Error("Error while closing file")
		}
		mu.(*sync.RWMutex).Unlock()
		mu.(*sync.RWMutex).RLock()

		executeUploadHook(slot)

		mu.(*sync.RWMutex).RUnlock()
	})
	if err != nil {
		logrus.WithError(err).Error("Error while upgrading connexion")
	}
}

func executeUploadHook(slot UploadSlot) {

	path := filepath.Join(WorkDir, slot.FileName)

	commandStr := strings.Replace(slot.UploadHook, "$1", "\""+path+"\"", -1)
	cmd := exec.Command("bash", "-c", commandStr)
	output, err := cmd.CombinedOutput()

	logrus.WithFields(logrus.Fields{
		"output":     string(output),
		"err":        err,
		"commandStr": commandStr,
	}).Info("Execute upload hook")
}

func (api *WebApi) ReadUploadSlot(ctx *fasthttp.RequestCtx) {

	tokenStr := string(ctx.Request.Header.Peek("X-Upload-Token"))
	if tokenStr == "" {
		tokenStr = string(ctx.Request.URI().QueryArgs().Peek("token"))
	}

	slot := UploadSlot{}
	err := api.db.Where("token=?", tokenStr).First(&slot).Error

	if err != nil {
		ctx.Response.Header.SetStatusCode(404)
		logrus.WithError(err).WithFields(logrus.Fields{
			"token": tokenStr,
		}).Warning("Upload slot not found")
		return
	}

	logrus.WithFields(logrus.Fields{
		"slot": slot,
	}).Info("Reading")

	path := filepath.Join(WorkDir, slot.FileName)

	mu, _ := Mutexes.LoadOrStore(slot.Token, &sync.RWMutex{})
	mu.(*sync.RWMutex).RLock()
	fp, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		logrus.WithError(err).Error("Error while opening file for reading")
		mu.(*sync.RWMutex).RUnlock()
		return
	}

	buf := make([]byte, WsBufferSize)
	response := ctx.Response.BodyWriter()
	for {
		read, err := fp.Read(buf)
		_, _ = response.Write(buf[:read])
		if err == io.EOF {
			break
		}
	}
	err = fp.Close()
	if err != nil {
		logrus.WithError(err).Error("Error while closing file for reading")
	}
	mu.(*sync.RWMutex).RUnlock()
}

func (api *WebApi) allocateUploadSlot(req *AllocateUploadSlotRequest) error {

	slot := &UploadSlot{
		MaxSize:       req.MaxSize,
		FileName:      req.FileName,
		Token:         req.Token,
		ToDisposeDate: req.ToDisposeDate,
		UploadHook:    req.UploadHook,
	}

	logrus.WithFields(logrus.Fields{
		"slot": slot,
	}).Info("Allocated new upload slot")

	return api.db.Create(slot).Error
}
