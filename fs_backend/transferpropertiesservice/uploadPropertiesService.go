package transferpropertiesservice

import (
	"time"

	"fs_backend/models"

	"github.com/jellydator/ttlcache/v3"
)

type TransferPropertiesService struct {
	isRunning bool
	cache     *ttlcache.Cache[string, models.FileTransferProperties]
}

func (ups *TransferPropertiesService) Start() {
	ups.cache = ttlcache.New[string, models.FileTransferProperties](
		ttlcache.WithTTL[string, models.FileTransferProperties](24 * time.Hour),
	)
	go ups.cache.Start()
	ups.isRunning = true
}

func (ups *TransferPropertiesService) Stop() {
	ups.cache.Stop()
	ups.cache.DeleteAll()
	ups.isRunning = false
}
