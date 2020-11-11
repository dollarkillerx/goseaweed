package goseaweed

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type seaweedFS struct {
	serverURL string
	timeout   time.Duration
}

func NewSeaweedFs(serverURL string, timeout time.Duration) Seaweed {
	if timeout <= 100 {
		timeout = time.Second * 10
	}
	return &seaweedFS{serverURL: serverURL, timeout: timeout}
}

func (s *seaweedFS) PutObject(objectName string, content []byte) error {
	url := fmt.Sprintf("%s/%s", s.serverURL, objectName)
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.WriteField("file", bytes.NewBuffer(content).String())
	if err != nil {
		return errors.WithStack(err)
	}
	err = writer.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	client := &http.Client{Timeout: s.timeout}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return errors.WithStack(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	if res.StatusCode != http.StatusCreated {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.WithStack(err)
		}
		return fmt.Errorf("create failed, response:%s", string(body))
	}
	return nil
}

func (s *seaweedFS) GetObject(objectName string) ([]byte, error) {
	body, err := http.Get(fmt.Sprintf("%s/%s", s.serverURL, objectName))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer body.Body.Close()
	return ioutil.ReadAll(body.Body)
}

func (s *seaweedFS) RemoveObject(objectName string) error {
	url := fmt.Sprintf("%s/%s", s.serverURL, objectName)
	method := "DELETE"

	client := &http.Client{Timeout: s.timeout}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return errors.WithStack(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.WithStack(err)
	}
	if res.StatusCode != http.StatusNoContent {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.WithStack(err)
		}
		return fmt.Errorf("delete failed, response:%s", string(body))
	}
	return nil
}
