package webdav

import "strings"

// A collection callback
type CollectionCallback func(*Collection) error

// A entry callback
type EntryCallback func(*Entry) error

// The collection type
type Collection struct {
	Entity
}

// Create a new collection
func NewCollection(client *Client, location string) *Collection {
	return &Collection{Entity{
		Client:   client,
		Location: location,
	}}
}

// Parse the given response, returning the entity and it's collection status
func (col *Collection) parse(res *Response) (*Entity, bool) {
	var isDir bool
	out := &Entity{
		Client:   col.Entity.Client,
		Location: res.Href,
	}

	for _, prop := range res.PropStat.Props {
		// Check if this resource is a collection
		if prop.ResourceType.IsCollection() {
			isDir = true
		}

		out.CreatedTime = prop.CreationDate
		out.LastModifiedTime = prop.LastModified.Time
	}

	return out, isDir
}

// Return the Collections for this Collection
func (col *Collection) Collections(recursive bool) ([]Collection, error) {
	out := []Collection{}

	err := col.Walk(nil,
		func(collection *Collection) error {
			out = append(out, *collection)
			return nil
		}, recursive)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Return the contents for this Collection
func (col *Collection) Contents(recursive bool) ([]interface{}, error) {
	out := []interface{}{}

	err := col.Walk(
		func(entry *Entry) error {
			out = append(out, *entry)
			return nil
		},
		func(collection *Collection) error {
			out = append(out, *collection)
			return nil
		}, recursive)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Create this Collection
func (col *Collection) Create() error {
	_, err := col.Client.mkcol(col.Location)
	return err
}

// Return the Entries for this Collection
func (col *Collection) Entries(recursive bool) ([]Entry, error) {
	out := []Entry{}

	err := col.Walk(
		func(entry *Entry) error {
			out = append(out, *entry)
			return nil
		}, nil, recursive)

	if err != nil {
		return nil, err
	}

	return out, nil
}

// Walk this collection
func (col *Collection) Walk(
	entryCallback EntryCallback,
	collectionCallback CollectionCallback,
	recursive bool) error {

	res, err := col.Client.propfind(col.Location, 1)

	if err != nil {
		return err
	}

	// TODO: check for response errors here

	data, err := col.Client.read(res)

	if err != nil {
		return err
	}

	ms, err := NewMultiStatusFromData(data)

	if err != nil {
		return err
	}

	for _, response := range ms.Responses {
		// Skip those without a successful status
		if !strings.HasSuffix(response.PropStat.Status, "200 OK") {
			continue
		}

		// Skip this location
		if response.Href == col.Location {
			continue
		}

		entity, isDir := col.parse(&response)

		// Create the correct type and append it to the slice
		if isDir {
			collection := &Collection{*entity}

			if collectionCallback != nil {
				err = collectionCallback(collection)

				if err != nil {
					return err
				}
			}

			// Recurse if needed
			if recursive {
				err = collection.Walk(entryCallback, collectionCallback, recursive)
			}
		} else if entryCallback != nil {
			err = entryCallback(&Entry{*entity})
		}

		if err != nil {
			return err
		}
	}

	return nil
}
