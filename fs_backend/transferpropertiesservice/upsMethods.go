package transferpropertiesservice

import (
	"time"

	"fs_backend/apierrors"
	"fs_backend/models"
)

func (ups *TransferPropertiesService) Get(uploadId string) (models.FileTransferProperties, error) {
	if ups.cache.Has(uploadId) {
		res := ups.cache.Get(uploadId)
		if res.IsExpired() {
			ups.cache.Delete(uploadId)
			return models.FileTransferProperties{}, apierrors.UploadIdNotFound{UploadId: uploadId}
		} else {
			return res.Value(), nil
		}
	} else {
		return models.FileTransferProperties{}, apierrors.UploadIdNotFound{UploadId: uploadId}
	}
}

func (ups TransferPropertiesService) Set(properties models.FileTransferProperties) {
	ups.cache.Set(properties.LinkId, properties, time.Hour*24)
}

func (ups TransferPropertiesService) Delete(uploadId string) {
	ups.cache.Delete(uploadId)
}
