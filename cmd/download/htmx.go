package download

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ross96D/gohtmx-stack/cmd/download/storage"
)

const url = "https://unpkg.com/htmx.org@latest"

func LatestHtmx() (htmx []byte, err error) {
	var ver string
	ver, err = LatestHtmxVersion()
	if err != nil {
		return
	}
	if storage.GetHtmxVersion() != ver {
		storage.SetHtmxVersion(ver)
		return downloadHtmx(ver)
	}
	if htmx, err = storage.GetHtmxFile(); err != nil {
		return downloadHtmx(ver)
	}
	return
}

func LatestHtmxVersion() (version string, err error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Head(url)
	if err != nil {
		err = fmt.Errorf("htmx-version %w", err)
		return
	}
	url, err := resp.Location()
	if err != nil {
		err = fmt.Errorf("htmx-version %w", err)
		return
	}
	splitted := strings.Split(url.Path, "@")
	if len(splitted) != 2 {
		err = errors.New("location was not formatted as expected. Result was " + url.Path + ". Expect somehting like /htmx.org@1.9.8 that has only one '@' and after the '@' has the version")
		return
	}
	version = splitted[1]
	return
}

func downloadHtmx(version string) ([]byte, error) {
	resp, err := http.Get(strings.ReplaceAll(url, "latest", version))
	if err != nil {
		err = fmt.Errorf("download-htmx %w", err)
		return nil, err
	}
	defer resp.Body.Close()
	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("download-htmx %w", err)
		return nil, err
	}
	go storage.SaveHtmxFile(b)
	return b, nil
}
