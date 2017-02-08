package webdav

import "time"

// The base entity type
type Entity struct {
	// The client for this Entity
	Client *Client

	// The location for this Entity
	Location string

	// The created time for this Entity
	CreatedTime time.Time

	// The last modified time for this Entity
	LastModifiedTime time.Time
}

// Copy this Entity
func (entity *Entity) Copy(location string) error {
	return nil
}

// Get the created time for this Entity
func (entity *Entity) Created() time.Time {
	return entity.CreatedTime
}

// Delete this Entity
func (entity *Entity) Delete() error {
	_, err := entity.Client.delete(entity.Location)
	return err
}

// Returns true if this Entity exists, or false otherwise
func (entity *Entity) Exists() (bool, error) {
	res, err := entity.Client.propfind(entity.Location, 1)

	if err != nil {
		return false, err
	}

	return !entity.Client.isResponseError(res), nil
}

// Get the last modified time for this Entity
func (entity *Entity) LastModified() time.Time {
	return entity.LastModifiedTime
}

// Move this Entity
func (entity *Entity) Move(location string) error {
	return nil
}
