package commands

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

func NewUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func TestUUID(t *testing.T) {
	t.Log(uuid.New().String())
	t.Log(uuid.New().String())
}
