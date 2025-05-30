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
	Env    *env.Env
	client *gowebdav.Client
}

func NewWebDAV(env *env.Env) WebDAVItf {
	return &WebDAV{
		Env:    env,
		client: gowebdav.NewClient(env.WebDAVURL, env.WebDavUser, env.WebDAVPassword),
	}
}

func (w *WebDAV) Upload(path string, fileName string, data *[]byte) error {
	dirPath := w.Env.WebDAVPath + path
	fullPath := dirPath + "/" + fileName

	err := w.client.Mkdir(dirPath, os.FileMode(w.Env.WEbDAVPermission))
	if err != nil {
		log.Println(err)
	}

	err = w.client.Write(fullPath, *data, os.FileMode(w.Env.WEbDAVPermission))
	if err != nil {
		log.Println(err)
	}

	return err
}
