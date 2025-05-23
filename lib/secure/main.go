/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2024-12-25 16:08:39
 * @LastEditTime: 2024-12-25 16:16:56
 */
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
)

func encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainBytes := []byte(plaintext)
	// 对于CBC模式，需要使用PKCS#7填充plaintext到blocksize的整数倍
	plainBytes = pad(plainBytes, aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(plainBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainBytes)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, ct string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(data) < aes.BlockSize {
		return "", err
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)
	data = unpad(data, aes.BlockSize)
	return string(data), nil
}

// pad 使用PKCS#7标准填充数据
func pad(buf []byte, blockSize int) []byte {
	padding := blockSize - (len(buf) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(buf, padtext...)
}

// unpad 移除PKCS#7标准填充的数据
func unpad(buf []byte, blockSize int) []byte {
	if len(buf)%blockSize != 0 {
		return nil
	}
	padding := int(buf[len(buf)-1])
	return buf[:len(buf)-padding]
}

func main() {
	key := []byte("1234567890123456") // 16字节长度
	plaintext := "fsfdfdfdfdfdfdfdfdfdfdfdfd"

	// 加密
	ciphertext, err := encrypt(key, plaintext)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Encrypted:", ciphertext)

	// 解密
	decrypted, err := decrypt(key, ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Decrypted:", decrypted)
}
