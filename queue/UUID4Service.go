package queue

import (
    "crypto/rand"
)

type uuid4Service struct {}

// NewUUUID4Service implements a UUIDService by generating Version-4 UUIDs
// from random numbers.
func NewUUID4Service() UUIDService {
    return uuid4Service{}
}

func (s uuid4Service) Generate() ([]byte, error) {
    b := make([]byte, 16) // 128 bits
    _, err := rand.Read(b)
    if err != nil { return nil, err }
    
    // set the four most significant bits of the 7th byte to 0100 (4)
    b[6] = (b[6] & 0b00001111) | 0b01000000
    
    // set the two most significant bits of the 9th byte to 10
    b[8] = (b[8] & 0b00111111) | 0b10000000
    
    return b, nil
}
