package api

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/fasthttp/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"os"
	"path/filepath"
	"time"
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

type WebApi struct {
	server      fasthttp.Server
	db          *gorm.DB
	MotdMessage *websocket.PreparedMessage
	Cron        *cron.Cron
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
	api.Cron = cron.New()

	logrus.SetLevel(getLogLevel())

	router := fasthttprouter.New()
	router.GET("/", LogRequestMiddleware(Index))

	router.POST("/slot", LogRequestMiddleware(api.AllocateUploadSlot))
	router.GET("/slot", LogRequestMiddleware(api.ReadUploadSlot))
	router.GET("/slot_info", LogRequestMiddleware(api.UploadSlotInfo))
	router.GET("/upload", LogRequestMiddleware(api.Upload))

	api.server = fasthttp.Server{
		Handler: router.Handler,
		Name:    "ws_bucket",
	}

	api.db = db
	db.AutoMigrate(&UploadSlot{})

	return api
}

func (api *WebApi) Run() {
	address := GetServerAddress()
	api.setupCronJobs()

	logrus.WithFields(logrus.Fields{
		"addr": address,
	}).Info("Starting web server")

	err := api.server.ListenAndServe(address)
	if err != nil {
		logrus.Fatalf("Error in ListenAndServe: %s", err)
	}
}
func (api *WebApi) setupCronJobs() {
	duration, _ := time.ParseDuration("5m")
	api.Cron = cron.New()
	schedule := cron.Every(duration)
	api.Cron.Schedule(schedule, cron.FuncJob(api.DisposeStaleUploadSlots))

	api.Cron.Start()

	logrus.WithFields(logrus.Fields{
		"every": duration,
	}).Info("Scheduled job for DisposeStaleUploadSlots")
}

//TODO: Move those to a different file/package
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
