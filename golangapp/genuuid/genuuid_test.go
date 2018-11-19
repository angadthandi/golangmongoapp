package genuuid

import (
	"fmt"
	"testing"
)

func TestGenUUID(t *testing.T) {
	uuid := GenUUID()

	fmt.Printf("uuid: %v\n", uuid)

	if uuid == "" {
		t.Error("TestGenUUID Failed to generate uuid")
	}
}
