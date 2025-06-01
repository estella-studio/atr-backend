package webdav

import (
	"log"
	"os"

	"github.com/estella-studio/leon-backend/internal/infra/env"
	"github.com/studio-b12/gowebdav"
)

type WebDAVItf interface {
	Upload(path string, fileName string, data *[]byte) error
}

type WebDAV struct {
	client     *gowebdav.Client
	url        string
	user       string
	password   string
	path       string
	permission os.FileMode
}

func NewWebDAV(env *env.Env) WebDAVItf {
	return &WebDAV{
		client:     gowebdav.NewClient(env.WebDAVURL, env.WebDavUser, env.WebDAVPassword),
		url:        env.WebDAVURL,
		user:       env.WebDavUser,
		password:   env.WebDAVPassword,
		path:       env.WebDAVPath,
		permission: os.FileMode(env.WEbDAVPermission),
	}
}

func (w *WebDAV) Upload(path string, fileName string, data *[]byte) error {
	dirPath := w.path + path
	fullPath := dirPath + "/" + fileName

	err := w.client.Mkdir(dirPath, w.permission)
	if err != nil {
		log.Printf("WebDAV: %v\n", err)
	}

	err = w.client.Write(fullPath, *data, w.permission)
	if err != nil {
		log.Printf("WebDAV: %v\n", err)
	}

	return err
}
