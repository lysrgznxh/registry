package files

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

func DownLoad(ctx context.Context, url string, path string) error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr, Timeout: time.Hour}

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	// Get the data
	//resp, err := http.Get(url)
	defer resp.Body.Close()
	// Create output file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	// copy stream
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
