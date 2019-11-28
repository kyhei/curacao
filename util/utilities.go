package util

import (
	"bytes"
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

// ParseJSONRequestBody Parse HTTP JSON request to interface{}
func ParseJSONRequestBody(r *http.Request, data interface{}) error {

	if _, err := strconv.Atoi(r.Header.Get("Content-Length")); err != nil {
		return errors.New("bad Content-Length")
	}

	body := r.Body
	defer body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, body); err != nil {
		return errors.New(err.Error())
	}

	err := json.Unmarshal(buf.Bytes(), data)

	if err != nil {
		return errors.New("could not parse json string")
	}

	return nil
}

// SaveUploadedFile Save uploaded file to local disk
func SaveUploadedFile(r *http.Request, attrName string, storagePath string) (*string, error) {

	if !strings.HasSuffix(storagePath, "/") {
		return nil, errors.New("storagePath must has '/' suffix")
	}

	file, fileHeader, err := r.FormFile(attrName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.Mkdir(storagePath, 0777)
	}

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
