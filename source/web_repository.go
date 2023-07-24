package source

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

type WebRepository struct {
	lastUpdateSeconds int64
	data              string
	url               *url.URL
}

func (w *WebRepository) GetData(ctx context.Context) (string, error) {
	if ((time.Now().Unix() - w.lastUpdateSeconds) < 10) && w.data != "" {
		logrus.Debug("returning cached file")
		return w.data, nil
	}
	logrus.Debug("fetching file")

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, w.url.String(), nil)
	if err != nil {
		logrus.Debug("error creating request")
		return "", err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Debug("error doing request")
		return "", err
	}
	defer resp.Body.Close()

	logrus.Debug("reading file")
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Debug("error reading file")
		return "", err
	}

	logrus.Debug("setting data")
	w.data = string(data)
	logrus.Debug("setting lastUpdateSeconds")
	w.lastUpdateSeconds = time.Now().Unix()
	return w.data, nil
}

func (w *WebRepository) GetType() string {
	return "http"
}

func (w *WebRepository) GetPath() string {
	return w.url.String()
}

func (w *WebRepository) GetUrl() *url.URL {
	return w.url
}

func NewWebRepository(webUrl string) (Repository, error) {
	parsedUrl, err := url.Parse(webUrl)
	if err != nil {
		return nil, err
	}
	return &WebRepository{url: parsedUrl}, nil
}