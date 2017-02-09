package webdav

import "strings"

// An entry
type Entry struct {
	Entity
}

// Create a new entry with the given location and client
func NewEntry(client *Client, location string) *Entry {
	return &Entry{Entity{
		Client:   client,
		Location: location,
	}}
}

// Copy this Entry
func (entry *Entry) Copy(location string) error {
	_, err := entry.Client.copy(entry.Location, location, 0)
	return err
}

// Get the data for this Entry
func (entry *Entry) Data() ([]byte, error) {
	res, err := entry.Client.get(entry.Location)

	if err != nil {
		return nil, err
	}

	return entry.Client.read(res)
}

// Ensure the parent exists for this Entry
func (entry *Entry) EnsureParentExists() error {
	parent := entry.Parent()
	exists, err := parent.Exists()

	if err != nil {
		return err
	}

	if !exists {
		return parent.Create()
	}

	return nil
}

// Get the parent Collection resource for this Entry
func (entry *Entry) Parent() *Collection {
	idx := strings.LastIndex(entry.Location, "/")
	return entry.Client.Collection(entry.Location[:idx+1])
}

// Upload the given data to this Entry
func (entry *Entry) Upload(data []byte) error {
	_, err := entry.Client.put(entry.Location, data)
	return err
}
