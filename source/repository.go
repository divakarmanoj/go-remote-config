package source

// Repository is an interface that defines the contract for a configuration data repository.
// Any type implementing this interface must provide methods to retrieve the configuration data
// and to refresh the data when required.
type Repository interface {
	GetName() string
	// GetData returns the configuration data as a map of configuration names to their respective models.
	GetData(string) (interface{}, bool)

	// GetRawData returns the raw data of the configuration file.
	GetRawData() []byte

	// Refresh updates the configuration data by fetching the latest data from the data source,
	// such as a file, database, or remote service. The method should handle any necessary
	// synchronization or locking to ensure safe access to the data during the refresh process.
	// It should return an error if there was a problem while fetching or updating the data.
	// The caller of this method should handle the error appropriately.
	Refresh() error
}
