package webdav

// The Unauthorized error
type UnauthorizedError struct{}

func (err UnauthorizedError) Error() string {
	return "Unauthorized"
}
