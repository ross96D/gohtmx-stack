package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

type htmxStorage struct {
	Version string `json:"version"`
	Error   bool   `json:"error"`
}

type storage struct {
	Htmx htmxStorage `json:"htmx"`
}

var data storage
var StoragePath string

const HtmxFileName = "htmx.js"

func update() {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("Warning: saving data file failed %v\n", err)
		}
	}()
	b, err := json.Marshal(data)
	f, err := os.Create(path.Join(StoragePath, "data.json"))
	_, err = f.Write(b)
}

func SetHtmxVersion(location string) {
	data.Htmx.Version = location
	update()
}

func GetHtmxVersion() string {
	return data.Htmx.Version
}

func SaveHtmxFile(b []byte) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("save-htmx-file %w", err)
		}
	}()
	f, err := os.Create(path.Join(StoragePath, HtmxFileName))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		data.Htmx.Error = true
		return err
	}
	return nil
}

func GetHtmxFile() (htmx []byte, err error) {
	if data.Htmx.Error {
		err = errors.New("htmx file is corrupted")
		return
	}
	htmxFile := path.Join(StoragePath, HtmxFileName)
	var f *os.File
	f, err = os.Open(htmxFile)
	if err != nil {
		err = fmt.Errorf("get-htmx-file %w", err)
		return
	}
	defer f.Close()
	if htmx, err = io.ReadAll(f); err != nil {
		err = fmt.Errorf("get-htmx-file %w", err)
	}
	return
}

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		//! TODO maybe this does not have to panic, just remove al sort of cache behaivour
		panic(err)
	}
	StoragePath = path.Join(home, ".ghtstack")

	dataLocation := path.Join(StoragePath, "data.json")

	if _, err := os.Stat(dataLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dataLocation), 0755); err != nil {
			panic(err)
		}
		if _, err := os.Create(dataLocation); err != nil {
			panic(err)
		}
	}
	f, err := os.ReadFile(dataLocation)
	if err != nil {
		//! TODO maybe this does not have to panic, just remove al sort of cache behaivour
		panic(err)
	}
	json.Unmarshal(f, &data)
}
