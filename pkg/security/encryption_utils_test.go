package security_test

import (
	"ResiSync/pkg/config"
	"ResiSync/pkg/security"
	"testing"
)

func init() {
	config.LoadEnv()
	config.LoadConfig()
}

func TestEncrypt(t *testing.T) {
	t.Error(security.EncryptPassword(""))
}

func TestHashing(t *testing.T) {
	// var requestContext = models.ResiSyncRequestContext{Log: logger.GetBasicLogger()}
	// pass, hash, _ := security.Hashpassword(requestContext, 16, "1234567890")

	t.Error(security.ComparePassword("4ce064409fb2c01feeedcaca0a5fcd6e1f1c9c44d9722d1402004a5fe0ed212a6d72c9642e77b51831610d8018549cce27f6eac0d5143d0858f6313e8a1a5774", "c7baca383b11c46acff3b95179bbf342", "1234567890"))
}
