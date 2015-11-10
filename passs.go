package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"github.com/howeyc/gopass"
)

const (
	KeySize = 32
	DirPerm = 0700
)

var (
	HomePath = os.Getenv("HOME")
	AppDir   = HomePath + "/.passs"
	AppFile  = AppDir + "/passs"
)

var (
	AesKey  []byte
	Records RecordsSlice
)

var (
	filePrepared bool
)

func DirectoryExists() bool {
	stat, err := os.Stat(AppDir)
	if err != nil {
		return false
	}
	if !stat.IsDir() {
		return false
	}
	return true
}

func FileExists() bool {
	if !DirectoryExists() {
		return false
	}

	_, err := os.Stat(AppFile)
	if err != nil {
		return false
	}
	return true
}

func PadKey(key []byte) (nkey []byte) {
	nkey = make([]byte, KeySize)

	for i, _ := range nkey {
		if i < len(key) {
			nkey[i] = key[i]
		} else {
			nkey[i] = 0
		}
	}

	return
}

func PadMessage(msg []byte) (nmsg []byte) {
	length := len(msg)
	if length%aes.BlockSize != 0 {
		length = (length/aes.BlockSize + 1) * aes.BlockSize
	}
	nmsg = make([]byte, length)

	for i, _ := range nmsg {
		if i < len(msg) {
			nmsg[i] = msg[i]
		} else {
			nmsg[i] = 0
		}
	}

	return
}

func RequestKeyFirstTime() {
	fmt.Print("Make up the passs storage key: ")
	key := gopass.GetPasswd()
	fmt.Print("Repeat the key: ")
	key2 := gopass.GetPasswd()

	if len(key) != len(key2) {
		fmt.Println("Entered keys are different!")
		RequestKeyFirstTime()
		return
	}

	for i, _ := range key {
		if key[i] != key2[i] {
			fmt.Println("Entered keys are different!")
			RequestKeyFirstTime()
			return
		}
	}

	AesKey = PadKey(key)
}

func Encrypt(key, msg []byte) []byte {
	msg = PadMessage(msg)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	crypted := make([]byte, aes.BlockSize+len(msg))
	iv := crypted[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(crypted[aes.BlockSize:], msg)

	return crypted
}

func Decrypt(key, msg []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := msg[:aes.BlockSize]
	msg = msg[aes.BlockSize:]
	decrypted := make([]byte, len(msg))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, msg)

	return decrypted
}

func CreateFile(key []byte) {
	if !DirectoryExists() {
		err := os.Mkdir(AppDir, DirPerm)
		if err != nil {
			panic(err)
		}
	}

	fd, err := os.Create(AppFile)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	fd.Write(Encrypt(key, key))
}

func OpenDecryptFile(key []byte) []byte {
	fd, err := os.Open(AppFile)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	text, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(err)
	}

	decrypted := Decrypt(key, text)
	decrypted = bytes.Trim(decrypted, "\x00")
	return decrypted
}

func EncryptAndSave() {
	fd, err := os.Create(AppFile)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	sort.Sort(Records)
	jsonData, err := json.Marshal(&Records)
	if err != nil {
		panic(err)
	}

	text := append(AesKey, jsonData...)
	text = Encrypt(AesKey, text)

	fd.Write(text)
}

func CheckKey(key []byte) bool {
	text := OpenDecryptFile(key)
	textkey := text[:len(key)]

	for i, _ := range key {
		if key[i] != textkey[i] {
			return false
		}
	}

	return true
}

func RequestKey() {
	fmt.Print("Passs storage key: ")
	key := gopass.GetPasswd()
	if len(key) < 1 {
		os.Exit(0)
	}

	key = PadKey(key)

	if !CheckKey(key) {
		fmt.Println("Incorrent key!")
		RequestKey()
		return
	}

	AesKey = key
}

func PrepareFile() {
	if filePrepared {
		return
	}

	if !FileExists() {
		RequestKeyFirstTime()
		fmt.Printf("Creating a new passs storage file at %v\n", AppFile)
		CreateFile(AesKey)
	} else {
		RequestKey()
	}

	text := OpenDecryptFile(AesKey)
	if len(text) <= len(AesKey) {
		Records = nil
	} else {
		err := json.Unmarshal(text[len(AesKey):], &Records)
		if err != nil {
			panic(err)
		}
	}

	filePrepared = true
}
