package redisclient_test

import (
	"ResiSync/pkg/config"
	redisclient "ResiSync/pkg/database/redis"
	"context"
	"testing"
	"time"
)

func init() {
	config.LoadEnv()
	config.LoadConfig()
}

func TestEncrypt(t *testing.T) {
	c, _ := redisclient.GetRedisClient()
	c.Set(context.TODO(), "some", "13243", time.Hour)
}
