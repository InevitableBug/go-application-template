package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type EncryptedString string

/***********************************************************************************************************************
	The configuration structure.
 **********************************************************************************************************************/
type Configuration struct {
	Database DatabaseConfiguration `json:"database"`
	Log      LogConfiguration      `json:"log"`
	HTTP     HTTPConfiguration     `json:"http"`
}

type DatabaseConfiguration struct {
	Host               string          `json:"host"`
	Port               int             `json:"port"`
	User               string          `json:"user"`
	Password           EncryptedString `json:"password"`
	Schema             string          `json:"schema"`
	MaxOpenConnections int             `json:"max_open_connections"`
	MaxIdleConnections int             `json:"max_idle_connections"`
}

type LogConfiguration struct {
	LogFile  string `json:"log_file"`
	LogLevel string `json:"log_level"`
}

type HTTPConfiguration struct {
	Port         int    `json:"port"`
	DocumentRoot string `json:"document_root"`
}

/**********************************************************************************************************************/

// Instance is a singleton of the Configuration struct
var Instance *Configuration

const configFile = "config.json"

var key []byte

const keyStr = "b917c7265ef7d769df1be19212cd681803719e48ebd68f08f68b93cf13c7a2f5"

const encryptedPrefix = "{enc}"
const plainPrefix = "{pln}"

func init() {
	Instance = &Configuration{}

	key, _ = hex.DecodeString(keyStr)

	_, err := os.Stat(configFile)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, Instance)
	if err != nil {
		panic(err)
	}

}

func (es *EncryptedString) UnmarshalJSON(b []byte) error {
	var s string
	var plaintext, ciphertext string

	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("1")
		return err
	}

	if strings.HasPrefix(s, encryptedPrefix) {
		// if the property has the encrypted prefix decrypt the value and replace it in the struct

		ciphertext = s[len(encryptedPrefix):]
		plaintext, err = decrypt(ciphertext)

	} else if strings.HasPrefix(s, plainPrefix) {
		// if the property has the plain prefix print the encrypted string to stdout

		plaintext = s[len(plainPrefix):]
		ciphertext, err = encrypt(plaintext)

		fmt.Printf("encrypted property: %s\n", ciphertext)

	} else {
		plaintext = s
	}

	*es = EncryptedString(plaintext)

	return err
}

func encrypt(plaintext string) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	return hex.EncodeToString(ciphertext) + "." + hex.EncodeToString(nonce), nil
}

func decrypt(encoded string) (string, error) {

	s := strings.Split(encoded, ".")
	if len(s) != 2 {
		return "", errors.New("invalid ciphertext")
	}

	ciphertext, err := hex.DecodeString(s[0])
	if err != nil {
		return "", err
	}

	nonce, err := hex.DecodeString(s[1])
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
