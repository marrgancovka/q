package pkg

import (
	"hack/internal/config"
	"time"
)

func GetCurrentTime() string {
	return time.Now().Format(config.LogsTimeFormat)
}
