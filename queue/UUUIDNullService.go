package queue

type uuidNullService struct {}
var uuidNullValue = [0]byte{}

// NewUUUID4Service implements a UUIDService by generating UUIDs that are an
// empty array of bytes, for when you don't need UUIDs e.g. where you only
// have at-least-once delivery requirements for a remote target.
func NewUUIDNuullService() UUIDService {
    return uuidNullService{}
}

func (s uuidNullService) Generate() ([]byte, error) {
    return uuidNullValue[:], nil
}
