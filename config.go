package main

import (
	"fmt"
	"strings"
        "encoding/json"
	"os"
	"io/ioutil"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)



var bt = []byte{46, 57, 24, 85, 35, 24, 74, 35, 87, 88, 98, 66, 32, 14, 05, 35}

func Encode(b []byte) string {
 return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) []byte {
 data, err := base64.StdEncoding.DecodeString(s)
 if err != nil {
  panic(err)
 }
 return data
}

func Encrypt(text, key string) (string, error) {
 block, err := aes.NewCipher([]byte(key))
 if err != nil {
  return "", err
 }
 plainText := []byte(text)
 cfb := cipher.NewCFBEncrypter(block, bt)
 cipherText := make([]byte, len(plainText))
 cfb.XORKeyStream(cipherText, plainText)
 return Encode(cipherText), nil
}

func Decrypt(text, key string) (string, error) {
 block, err := aes.NewCipher([]byte(key))
 if err != nil {
  return "", err
 }
 cipherText := Decode(text)
 cfb := cipher.NewCFBDecrypter(block, bt)
 plainText := make([]byte, len(cipherText))
 cfb.XORKeyStream(plainText, cipherText)
 return string(plainText), nil
}

// config
type Configuration struct {
    User    []string
    Password   []string
    Server   []string
    Port   []string
    Db   []string
    Init   []string
}

var sqlcred = getconfig()

func getconfig() string {

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	errcode := decoder.Decode(&configuration)
	if errcode != nil {
		fmt.Println("error:", errcode)
	}


	if strings.Join(configuration.Init," ") == "1" {
		fmt.Println(md5hash(`admin`+keysalt+``))
		fmt.Println("Initial Passwort wird geändert...")
	
		changepw, err := ioutil.ReadFile("conf.json")
		if err != nil {
			fmt.Println("error:", errcode)
		}
		encText, err := Encrypt(strings.Join(configuration.Password," "), key)
		newpw := bytes.Replace(changepw, []byte(`"Password": ["`+strings.Join(configuration.Password," ")+`"],`), []byte(`"Password": ["`+encText+`"],`), -1)
		if err = ioutil.WriteFile("conf.json", newpw, 0660); err != nil {
			fmt.Println("error:", errcode)
		}  
	
	
		changeinit, err := ioutil.ReadFile("conf.json")
		if err != nil {
			fmt.Println("error:", errcode)
		}
		newinit := bytes.Replace(changeinit, []byte(`"Init": ["1"]`), []byte(`"Init": ["0"]`), -1)
		if err = ioutil.WriteFile("conf.json", newinit, 0660); err != nil {
			fmt.Println("error:", errcode)
		} 

		return ``+strings.Join(configuration.User," ")+`:`+strings.Join(configuration.Password," ")+`@tcp(`+strings.Join(configuration.Server," ")+`:`+strings.Join(configuration.Port," ")+`)/`+strings.Join(configuration.Db," ")+``

	} else
	{
		decText, err := Decrypt(strings.Join(configuration.Password," "), key)
		if err != nil {
			fmt.Println("error:", errcode)
		}

		return ``+strings.Join(configuration.User," ")+`:`+decText+`@tcp(`+strings.Join(configuration.Server," ")+`:`+strings.Join(configuration.Port," ")+`)/`+strings.Join(configuration.Db," ")+``
	}

}
