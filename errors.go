package webdav

// The NotFound error
type NotFoundError struct {
	// The url for this error
	Url string
}

func (err NotFoundError) Error() string {
	return "NotFound"
}

// The Unauthorized error
type UnauthorizedError struct {
	// The url for this error
	Url string
}

func (err UnauthorizedError) Error() string {
	return "Unauthorized"
}

// The Unknown error
type UnknownError struct {
	// The url for this error
	Url string
}

func (err UnknownError) Error() string {
	return "Unknown"
}
