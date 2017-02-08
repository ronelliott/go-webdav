package webdav

import (
	"encoding/xml"
	"strings"
	"time"
)

// Allow parsing the last modified time
type LastModifiedTime struct {
	time.Time
}

// Unmarshal the xml
func (lmt *LastModifiedTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	value := ""
	err := d.DecodeElement(&value, &start)

	if err != nil {
		return err
	}

	parsed, err := time.Parse(time.RFC1123, value)

	if err != nil {
		return err
	}

	*lmt = LastModifiedTime{parsed}
	return nil
}

// The lockentry type
type LockEntry struct {
	Scope LockScope `xml:"lockscope"`
	Type  LockType  `xml:"locktype"`
}

// The lockscope type
type LockScope struct {
	Exclusive *string `xml:"exclusive"`
	Shared    *string `xml:"shared"`
}

// Return true if this LockScope is exclusive
func (scope *LockScope) IsExclusive() bool {
	return scope.Exclusive != nil
}

// Return true if this LockScope is shared
func (scope *LockScope) IsShared() bool {
	return scope.Shared != nil
}

// The locktype type
type LockType struct {
	Read  *string `xml:"read"`
	Write *string `xml:"write"`
}

// Return true if this LockType is read
func (scope *LockType) IsRead() bool {
	return scope.Read != nil
}

// Return true if this LockType is write
func (scope *LockType) IsWrite() bool {
	return scope.Write != nil
}

// The multistatus type
type MultiStatus struct {
	Responses []Response `xml:"response"`
}

// Create a new MultiStatus instance from the given data
func NewMultiStatusFromData(data []byte) (*MultiStatus, error) {
	ms := &MultiStatus{}
	err := xml.Unmarshal(data, ms)

	if err != nil {
		return nil, err
	}

	return ms, nil
}

// The propstat type
type PropStat struct {
	Props  []Prop `xml:"prop"`
	Status string `xml:"status"`
}

// The prop type
type Prop struct {
	ContentLength uint64           `xml:"getcontentlength"`
	CreationDate  time.Time        `xml:"creationdate"`
	Etag          string           `xml:"getetag"`
	Executable    bool             `xml:"executable"`
	LastModified  LastModifiedTime `xml:"getlastmodified"`
	ResourceType  ResourceType     `xml:"resourcetype"`
	SupportedLock SupportedLock    `xml:"supportedlock"`
}

// The resourcetype type
type ResourceType struct {
	Collection *string `xml:"collection"`
}

// Return true if this ResourceType is a collection
func (rt *ResourceType) IsCollection() bool {
	return rt.Collection != nil
}

// The response type
type Response struct {
	Href     string   `xml:"href"`
	PropStat PropStat `xml:"propstat"`
}

// Return the location for this Response
func (res *Response) Location(base string) string {
	return strings.Replace(res.Href, base, "", 1)
}

// The supportedlock type
type SupportedLock struct {
	LockEntries []LockEntry `xml:"lockentry"`
}
