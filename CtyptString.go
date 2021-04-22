package main

import (
	"errors"
	"net"
)

//CrypterString - crypter string
type CrypterString struct {
	created bool
	count   uint64
	key     []byte
}

//NewCrypterString - create new crypter string
func NewCrypterString() CrypterString {
	return CrypterString{created: true}
}

//SetKey - set crypt/uncrypt key
func (cs *CrypterString) SetKey(key []byte) error {
	if cs.created == false {
		return errors.New("CrypterString not created")
	}

	cs.count = 0
	cs.key = key

	return nil
}

//CryptString - crypt string
func (cs *CrypterString) CryptString(str *string) []byte {
	byteArr := []byte(*str)
	byteCrypt := make([]byte, len(byteArr))

	for i := 0; i < len(byteArr); i++ {
		byteCrypt[i] = cs.CryptChar(byteArr[i])
	}
	return byteCrypt
}

//CryptBytes - Crypt bytes
func (cs *CrypterString) CryptBytes(byteArr []byte) []byte {
	byteCrypt := make([]byte, len(byteArr))
	for i := 0; i < len(byteArr); i++ {
		byteCrypt[i] = cs.CryptChar(byteArr[i])
	}
	return byteCrypt
}

//UncryptString - uncrypt string
func (cs *CrypterString) UncryptString(byteCrypt []byte) string {
	byteArr := make([]byte, len(byteCrypt))

	for i := 0; i < len(byteCrypt); i++ {
		byteArr[i] = cs.CryptChar(byteCrypt[i])
	}
	return string(byteArr)
}

//UncryptStringConn - uncrypt socket string
func (cs *CrypterString) UncryptStringConn(conn net.Conn) (string, error) {
	var myStr string
	bs := make([]byte, 1)
	for {
		_, err := conn.Read(bs)
		if err != nil {
			return "", err
		}
		b := cs.CryptChar(bs[0])
		if b == 0 {
			return myStr, nil
		}
		myStr += string(b)
	}
}

//CryptChar - crypt/uncrypt char
func (cs *CrypterString) CryptChar(b byte) byte {
	ret := b ^ cs.key[cs.count%uint64(len(cs.key))]
	cs.count++
	return ret
}
