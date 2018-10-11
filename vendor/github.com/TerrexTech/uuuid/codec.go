package uuuid

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// FromBytes returns a UUID generated from the raw byte slice input.
// It will return an error if the slice isn't 16 bytes long.
func FromBytes(input []byte) (UUID, error) {
	u, err := uuid.FromBytes(input)
	if err != nil {
		err = errors.Wrap(err, "Error converting Bytes to UUID")
		return UUID{}, err
	}
	return UUID{u}, nil
}

// FromBytesOrNil returns a UUID generated from the raw byte slice input.
// Same behavior as FromBytes(), but returns uuid.Nil instead of an error.
func FromBytesOrNil(input []byte) UUID {
	return UUID{uuid.FromBytesOrNil(input)}
}

// FromString returns a UUID parsed from the input string.
// Input is expected in a form accepted by UnmarshalText.
func FromString(input string) (UUID, error) {
	u, err := uuid.FromString(input)
	if err != nil {
		err = errors.Wrap(err, "Error converting String to UUID")
		return UUID{}, err
	}
	return UUID{u}, nil
}

// FromStringOrNil returns a UUID parsed from the input string.
// Same behavior as FromString(), but returns uuid.Nil instead of an error.
func FromStringOrNil(input string) UUID {
	return UUID{uuid.FromStringOrNil(input)}
}
