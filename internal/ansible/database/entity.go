package database

// Entity interface represents an Entity that could be a group member in the database
type Entity interface {
	GetID() string
	Type() string
	GetName() string
	SetName(name string)

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}
