package util

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func GenID(prefix string) string {
	id := ksuid.New().String()
	return fmt.Sprintf("%s_%s", prefix, id)
}
