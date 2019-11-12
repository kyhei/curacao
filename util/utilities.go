package util

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ParseJSONRequestBody HTTPリクエストのJSONボディをオブジェクトにパースします。
func ParseJSONRequestBody(r *http.Request, data interface{}) error {

	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		return errors.New("bad Content-Length")
	}

	body := make([]byte, length)
	length, err = r.Body.Read(body)
	if err != nil && err != io.EOF {
		return errors.New("could not read body")
	}

	err = json.Unmarshal(body[:length], &data)
	if err != nil {
		return errors.New("could not parse json srting")
	}

	return nil
}

// SaveUploadedFile Httpリクエストのファイルボディをディスクに保存します。
func SaveUploadedFile(r *http.Request, attrName string, storagePath string) (*string, error) {

	if !strings.HasSuffix(storagePath, "/") {
		return nil, errors.New("storagePath must has '/' suffix")
	}

	file, fileHeader, err := r.FormFile(attrName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	saveFile, err := os.Create(storagePath + fileHeader.Filename)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer saveFile.Close()

	if _, err := io.Copy(saveFile, file); err != nil {
		return nil, err
	}

	savedPath, _ := filepath.Abs(storagePath + fileHeader.Filename)

	return &savedPath, nil
}

// ToJSON return JSON encoding of data
func ToJSON(data interface{}) []byte {
	resp, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	return resp
}
