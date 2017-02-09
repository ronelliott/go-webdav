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

// Parse the given response, returning the entity and it's collection status
func (entity *Entity) parse(res *Response) (*Entity, bool) {
	var isDir bool
	out := &Entity{
		Client:   entity.Client,
		Location: entity.Client.location(res),
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
		switch err.(type) {
		case NotFoundError:
			return false, nil
		}

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
	_, err := entity.Client.move(entity.Location, location)
	return err
}
