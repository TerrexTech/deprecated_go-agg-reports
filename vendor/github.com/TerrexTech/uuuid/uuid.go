package uuuid

import (
	"github.com/gocql/gocql"
	"github.com/gofrs/uuid"
)

// UUID represents a v4 UUID.
type UUID struct {
	uuid.UUID
}

// NewV4 returns a new v4 UUID
func NewV4() (UUID, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return UUID{}, err
	}
	return UUID{
		uid,
	}, nil
}

// MarshalCQL converts the uuid into GoCql-compatible []byte.
func (u UUID) MarshalCQL(info gocql.TypeInfo) ([]byte, error) {
	return u.MarshalBinary()
}

// UnmarshalCQL converts GoCql UUID-type to local UUID.
func (u *UUID) UnmarshalCQL(info gocql.TypeInfo, data []byte) error {
	return u.UnmarshalBinary(data)
}
