package goseaweed

type Seaweed interface {
	PutObject(objectName string, content []byte) error
	GetObject(objectName string) ([]byte, error)
	RemoveObject(objectName string) error
}
