package middlewares

import (
	"fmt"
	"time"
)

func GenerateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
