package util

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	url2 "net/url"
	"os"
)

type TinyPNGResponse struct {
	Input  Input
	Output OutPut
}

type Input struct {
	Size int64
	Type string
}

type OutPut struct {
	Size   int64
	Type   string
	Width  int64
	Height int64
	Ratio  float64
	Url    string
}

//CompressImageByTinyPNGAPI compress image by the tinypng api.
func CompressImageByTinyPNGAPI(ctx context.Context, fPath string, token string) error {
	const compressURL = "https://api.tinify.com/shrink"
	originFile, err := os.Open(fPath)
	originFile.Sync()
	if err != nil {
		return nil
	}
	defer originFile.Close()
	compressReq, err := http.NewRequestWithContext(ctx, "POST", compressURL, originFile)
	if err != nil {
		return err
	}
	compressReq.URL.User = url2.UserPassword("api", token)
	compressRes, err := http.DefaultClient.Do(compressReq)
	if err != nil {
		return err
	}
	defer compressRes.Body.Close()
	var tr TinyPNGResponse
	err = json.NewDecoder(compressRes.Body).Decode(&tr)
	if err != nil {
		return err
	}
	downloadReq, err := http.NewRequestWithContext(ctx, "GET", tr.Output.Url, nil)
	if err != nil {
		return err
	}
	downloadReq.URL.User = url2.UserPassword("api", token)
	downloadRes, err := http.DefaultClient.Do(downloadReq)
	if err != nil {
		return err
	}
	defer downloadRes.Body.Close()
	originFileAgain, err := os.OpenFile(fPath, os.O_WRONLY|os.O_TRUNC, 0666)
	defer originFileAgain.Close()
	_, err = io.Copy(originFileAgain, downloadRes.Body)
	return err
}
