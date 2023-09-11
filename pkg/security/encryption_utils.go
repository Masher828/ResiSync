package security

import (
	pkg_constants "ResiSync/pkg/constants"
	pkg_models "ResiSync/pkg/models"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	return b, err
}

func Aes256GCMEncode(plainText []byte, encryptionKey []byte) ([]byte, []byte, error) {

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	nonce, err := GenerateRandomBytes(12)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, []byte{}, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, []byte{}, err
	}

	cipherText := aesGCM.Seal(nil, nonce, plainText, nil)

	return cipherText, nonce, err
}

func Aes256GCMDecode(cipherText, encryptionKey, nonce []byte) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

func EncryptPassword(plainPassword string) (string, string, error) {
	return EncryptString(plainPassword, viper.GetString(pkg_constants.AwsEncryptionKey))
}

func EncryptString(plainText, base64EncryptionKey string) (string, string, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(base64EncryptionKey)

	if err != nil {
		log.Println(err)
		return "", "", err
	}

	encryptAes, nonce, err := Aes256GCMEncode([]byte(plainText), encryptionKey)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return base64.StdEncoding.EncodeToString(encryptAes), base64.StdEncoding.EncodeToString(nonce), nil
}

func DecryptPassword(encryptedPassword, passwordNonce string) (string, error) {
	return DecryptString(encryptedPassword, viper.GetString(pkg_constants.AwsEncryptionKey), passwordNonce)
}

func DecryptString(base64CipherText, base64EncryptionKey, base64Nonce string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(base64CipherText)
	if err != nil {
		return "", err
	}

	encryptionKey, err := base64.StdEncoding.DecodeString(base64EncryptionKey)
	if err != nil {
		return "", err
	}

	secretNonce, err := base64.StdEncoding.DecodeString(base64Nonce)
	if err != nil {
		return "", err
	}

	decryptedAes, err := Aes256GCMDecode(cipherText, encryptionKey, secretNonce)
	if err != nil {
		return "", err
	}

	return string(decryptedAes[:]), nil
}

func HashPasswordWithSalt(password string, salt []byte) string {

	var passwordBytes = append([]byte(password), salt...)

	var sha512Hasher = sha512.New()

	sha512Hasher.Write(passwordBytes)

	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	return hex.EncodeToString(hashedPasswordBytes)
}

func Hashpassword(requestContext pkg_models.ResiSyncRequestContext, hashLength int, password string) (string, string, error) {

	log := requestContext.Log

	salt, err := GenerateRandomBytes(hashLength)
	if err != nil {
		log.Error("error while generating salt", zap.Error(err))
		return "", "", err
	}

	return HashPasswordWithSalt(password, salt), hex.EncodeToString(salt), nil
}

func ComparePassword(hashedPassword, salt, password string) bool {

	saltByte, err := hex.DecodeString(salt)
	if err != nil {
		return false
	}
	return hashedPassword == HashPasswordWithSalt(password, saltByte)
}
