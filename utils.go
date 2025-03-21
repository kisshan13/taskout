package taskout

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

// GenerateID creates a lightweight unique ID based on the current timestamp and a random component.
func generateId() (string, error) {
	now := time.Now().UnixNano()
	randBytes := make([]byte, 4)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate id for task")
	}

	randPart := binary.BigEndian.Uint32(randBytes)

	return fmt.Sprintf("%x-%x", now, randPart), nil
}
