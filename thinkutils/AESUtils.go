package thinkutils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type aesutils struct {
}

func (this aesutils) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func (this aesutils) pkcs7Unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

func (this aesutils) EncryptString(szTxt, szKey string) (string, error) {
	key := []byte(szKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7填充
	plainBytes := []byte(szTxt)
	plainBytes = this.pkcs7Pad(plainBytes, aes.BlockSize)

	// 生成随机IV
	ciphertext := make([]byte, aes.BlockSize+len(plainBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// CBC模式加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainBytes)

	// 返回Base64编码
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (this aesutils) DecryptString(szTxt, szKey string) (string, error) {
	key := []byte(szKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 解码Base64
	cipherBytes, err := base64.StdEncoding.DecodeString(szTxt)
	if err != nil {
		return "", err
	}

	// 检查长度
	if len(cipherBytes) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// 提取IV
	iv := cipherBytes[:aes.BlockSize]
	cipherBytes = cipherBytes[aes.BlockSize:]

	// CBC模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherBytes, cipherBytes)

	// 去除填充
	plainBytes := this.pkcs7Unpad(cipherBytes)
	return string(plainBytes), nil
}
