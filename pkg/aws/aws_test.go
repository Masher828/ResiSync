package aws_services_test

import (
	"ResiSync/pkg/config"
	pkg_constants "ResiSync/pkg/constants"
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func init() {
	config.LoadEnv()
	config.LoadConfig()
}

func TestEncrypt(t *testing.T) {
	fmt.Println(viper.GetString(pkg_constants.AWSS3Bucket))
}
