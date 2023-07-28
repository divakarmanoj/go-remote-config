package go_remote_config

import (
	"context"
	"errors"
	"github.com/divakarmanoj/go-remote-config/source"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"time"
)

type Client struct {
	Repository      source.Repository
	RefreshInterval time.Duration
	cancel          context.CancelFunc
}

// NewClient creates a new Client with the provided context, repository,
// and refresh interval. It starts a background goroutine to periodically
// refresh the configuration data from the repository based on the given
// refresh interval. The function returns the created Client.
func NewClient(ctx context.Context, repository source.Repository, refreshInterval time.Duration) *Client {
	// Create a new context and its corresponding cancel function
	// for the Client. This allows us to control the lifetime of the
	// background refresh goroutine.
	ctx, cancel := context.WithCancel(ctx)

	// Create the Client instance with the provided repository and refresh interval.
	client := &Client{
		Repository:      repository,
		RefreshInterval: refreshInterval,
		cancel:          cancel, // Store the cancel function in the Client struct for later use.
	}

	// Refresh the configuration data for the first time to ensure the
	// Client is initialized with the latest data before it is used.
	err := client.Repository.Refresh()
	if err != nil {
		logrus.WithError(err).Error("error refreshing repository")
	}

	// Start the background refresh goroutine by calling the refresh function
	// with the newly created context and the client as arguments.
	go refresh(ctx, client)

	// Return the created Client instance, which is now ready to use.
	return client
}

// refresh is a goroutine that periodically refreshes the configuration data
// from the repository based on the provided refresh interval. It stops
// refreshing when the given context is canceled.
func refresh(ctx context.Context, client *Client) {
	ticker := time.NewTicker(client.RefreshInterval) // Create a new ticker with the given refresh interval
	for {
		select {
		case <-ticker.C:
			// The ticker has ticked, indicating it's time to refresh the data
			err := client.Repository.Refresh() // Call the Refresh method of the repository to update the configuration data
			if err != nil {
				logrus.WithError(err).Error("error refreshing repository")
			}
		case <-ctx.Done():
			// The context is canceled, indicating the refresh routine should stop
			return
		}
	}
}

// Close stops the background refresh goroutine of the Client by canceling
// its associated context. This function allows graceful termination of the
// background routine and prevents potential goroutine leaks. It should be
// called when the Client is no longer needed to release resources properly.
func (c *Client) Close() {
	// Call the Cancel function associated with the Client's context.
	// This cancels the context, causing the background refresh goroutine
	// (started by NewClient) to return and terminate gracefully.
	c.cancel()
}

// GetConfig retrieves the configuration with the given name from the repository
// and stores it in the provided data pointer. It returns an error if the
// configuration is not found, the data argument is not a non-nil pointer, or
// the type of the data is not compatible with the type in the repository.
func (c *Client) GetConfig(name string, data interface{}) error {
	// Get the configuration data from the repository
	config, ok := c.Repository.GetData(name)
	if !ok {
		return errors.New("config not found")
	}
	//
	marshal, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	// Unmarshal the configuration data into the provided data pointer
	err = yaml.Unmarshal(marshal, data)
	if err != nil {
		return err
	}

	return nil
}