package source

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// WebRepository is a struct that implements the Repository interface for
// handling configuration data fetched from a remote HTTP endpoint (web URL).
type WebRepository struct {
	sync.RWMutex                        // RWMutex to synchronize access to data during refresh
	Name         string                 // Name of the configuration source
	data         map[string]interface{} // Map to store the configuration data
	URL          *url.URL               // URL representing the remote HTTP endpoint (web URL)
	rawData      []byte                 // Raw data of the YAML configuration file
}

// GetName returns the name of the configuration source.
func (w *WebRepository) GetName() string {
	return w.Name
}

// GetData returns the configuration data as a map of configuration names to their respective models.
func (w *WebRepository) GetData(configName string) (config interface{}, isPresent bool) {
	w.RLock()
	defer w.RUnlock()
	config, isPresent = w.data[configName]
	return config, isPresent
}

// GetRawData returns the raw data of the YAML configuration file.
func (w *WebRepository) GetRawData() []byte {
	w.RLock()
	defer w.RUnlock()
	return w.rawData
}

// Refresh fetches the YAML file from the remote HTTP endpoint (web URL),
// unmarshal it into the data map.
func (w *WebRepository) Refresh() error {
	w.Lock()
	defer w.Unlock()

	// Create an HTTP request to fetch the YAML file from the remote web URL.
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, w.URL.String(), nil)
	if err != nil {
		logrus.Debug("error creating request")
		return err
	}

	// Perform the HTTP request to get the YAML file content.
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logrus.Debug("error doing request")
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.WithError(err).Debug("error closing response body")
		}
	}(resp.Body)

	// Read the file content from the response body.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Debug("error reading file")
		return err
	}

	// Unmarshal the YAML data into the data map.
	err = yaml.Unmarshal(data, &w.data)
	if err != nil {
		logrus.Debug("error unmarshalling file")
		return err
	}

	// Store the raw data of the YAML file.
	w.rawData = data

	return nil
}
