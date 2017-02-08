package webdav

import "fmt"
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

// Get the data for this Entry
func (entry *Entry) Data() ([]byte, error) {
	res, err := entry.Client.get(entry.Location)

	if err != nil {
		return nil, err
	}

	return entry.Client.read(res)
}

// Get the parent Collection resource for this Entry
func (entry *Entry) Parent() *Collection {
	idx := strings.LastIndex(entry.Location, "/")
	return entry.Client.Collection(entry.Location[:idx+1])
}

// Upload the given data to this Entry
func (entry *Entry) Upload(data []byte) error {
	res, err := entry.Client.put(entry.Location, data)

	if err != nil {
		return err
	}

	data, err = entry.Client.read(res)

	if err != nil {
		return err
	}

	fmt.Println(string(data))
	return nil
}
