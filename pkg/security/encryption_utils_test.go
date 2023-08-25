package security_test

import (
	"ResiSync/pkg/config"
	"ResiSync/pkg/security"
	"testing"
)

func init() {
	config.LoadConfig()
}

func TestEncrypt(t *testing.T) {
	t.Error(security.EncryptPassword("bCA6KgNNyz2N8SdkGyicWxg4OHLi93"))
}
