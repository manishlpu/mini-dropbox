package models

import "time"

const (
	STATUS_INACTIVE = iota
	STATUS_ACTIVE
)

var (
	FileStatus_name = map[int8]string{
		STATUS_INACTIVE: "inactive",
		STATUS_ACTIVE:   "active",
	}

	FileStatus_value = map[string]int8{
		"active":   STATUS_ACTIVE,
		"inactive": STATUS_INACTIVE,
	}
)

type FileStatus int8

func (x FileStatus) String() string {
	if val, ok := FileStatus_name[int8(x)]; ok {
		return val
	}
	return "inactive"
}

type Metadata struct {
	ID          int64      `db:"id" json:"id,omitempty"`
	Filename    string     `db:"filename" json:"filename"`
	SizeInBytes int64      `db:"size_in_bytes" json:"size_in_bytes"`
	S3ObjectKey string     `db:"s3_object_key" json:"s3_object_key"`
	Description string     `db:"description" json:"description,omitempty"`
	MimeType    string     `db:"mime_type" json:"mime_type,omitempty"`
	Status      FileStatus `db:"status" json:"-"`
	// PrevKey     string     `db:"prev_key" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
