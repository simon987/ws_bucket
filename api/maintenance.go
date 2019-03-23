package api

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

func (api *WebApi) DisposeStaleUploadSlots() {

	var toDispose []UploadSlot
	api.db.Where("? >= to_dispose_date", time.Now().Unix()).Find(&toDispose)

	for _, slot := range toDispose {
		path := filepath.Join(WorkDir, slot.FileName)

		err := os.Remove(path)
		api.db.Where("token = ?", slot.Token).Delete(UploadSlot{})

		logrus.WithFields(logrus.Fields{
			"fileName": slot.FileName,
			"err":      err,
		}).Trace("Deleted file")
	}

	logrus.WithFields(logrus.Fields{
		"staleUploadSlots": len(toDispose),
	}).Info("Disposed stale upload slots")
}
