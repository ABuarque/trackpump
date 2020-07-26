package filestorage

import (
	"fmt"
	"io"
	"trackpump/storage"
	"trackpump/usecase/service"
)

type fileStorage struct {
	client *storage.PCloudClient
}

// NewPcloudStorage returns a instance of pCloud client
func NewPcloudStorage(client *storage.PCloudClient) service.Storage {
	return &fileStorage{
		client: client,
	}
}

func (fs *fileStorage) Put(fileName string, data io.Reader) (string, error) {
	url, err := fs.client.Put(fileName, data)
	if err != nil {
		return "", fmt.Errorf("failed to save file on pCloud, err %q", err)
	}
	return url, nil
}
