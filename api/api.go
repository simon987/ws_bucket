package api

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/fasthttp/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
)

var WorkDir = getWorkDir()

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var info = Info{
	Name:    "ws_bucket",
	Version: "1.0",
}

var motd = WebsocketMotd{
	Info: info,
	Motd: "Hello, world",
}

type WebApi struct {
	server      fasthttp.Server
	db          *gorm.DB
	MotdMessage *websocket.PreparedMessage
}

func Index(ctx *fasthttp.RequestCtx) {
	Json(info, ctx)
}

func Json(object interface{}, ctx *fasthttp.RequestCtx) {

	resp, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = ctx.Write(resp)
	if err != nil {
		panic(err)
	}
}

func LogRequestMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		logrus.WithFields(logrus.Fields{
			"path":   string(ctx.Path()),
			"header": ctx.Request.Header.String(),
		}).Trace(string(ctx.Method()))

		h(ctx)
	})
}

func New(db *gorm.DB) *WebApi {

	if _, err := os.Stat(WorkDir); err != nil && os.IsNotExist(err) {
		_ = os.Mkdir(WorkDir, 0700)
	}

	api := &WebApi{}

	logrus.SetLevel(getLogLevel())

	router := fasthttprouter.New()
	router.GET("/", LogRequestMiddleware(Index))

	router.POST("/slot", LogRequestMiddleware(api.AllocateUploadSlot))
	router.GET("/slot", LogRequestMiddleware(api.ReadUploadSlot))
	router.GET("/upload", LogRequestMiddleware(api.Upload))

	api.server = fasthttp.Server{
		Handler: router.Handler,
		Name:    "ws_bucket",
	}

	api.db = db
	db.AutoMigrate(&UploadSlot{})

	api.setupMotd()

	return api
}

func (api *WebApi) setupMotd() {
	var data []byte
	data, _ = json.Marshal(motd)
	motdMsg, _ := websocket.NewPreparedMessage(websocket.TextMessage, data)
	api.MotdMessage = motdMsg
}

func (api *WebApi) Run() {
	address := GetServerAddress()

	logrus.WithFields(logrus.Fields{
		"addr": address,
	}).Info("Starting web server")

	err := api.server.ListenAndServe(address)
	if err != nil {
		logrus.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func GetServerAddress() string {
	serverAddress := os.Getenv("WS_BUCKET_ADDR")
	if serverAddress == "" {
		serverAddress = "0.0.0.0:3020"
	}
	return serverAddress
}

func getLogLevel() logrus.Level {
	levelStr := os.Getenv("WS_BUCKET_LOGLEVEL")
	if levelStr == "" {
		return logrus.TraceLevel
	} else {
		level, err := logrus.ParseLevel(levelStr)
		if err != nil {
			panic(err)
		}
		return level
	}
}

func getWorkDir() string {
	workDir := os.Getenv("WS_BUCKET_WORKDIR")
	if workDir == "" {
		path, _ := filepath.Abs("./data")
		return path
	} else {
		path, _ := filepath.Abs(workDir)
		return path
	}
}
