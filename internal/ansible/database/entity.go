package database

type Entity interface {
	GetId() string
	Type() string
	GetName() string
	SetName(name string)

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}
