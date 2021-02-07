/*
** Made by LUBERT Arnaud, Epitech Student (Promo 2023)
** 06/02/2021 arnaud.lubert@epitech.eu
**
*/
package main

import (
	"encoding/base64"
	"crypto/sha256"
	"encoding/hex"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/md5"
//	"crypto/cipher"
//	"crypto/sha512"
//	"crypto/rand"
//	"crypto/aes"
	"strings"
	"io"
)

func encodeBase64url(bytes []byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func encodeBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func decodeBase64(str string) string {
	data, err := base64.StdEncoding.DecodeString(str)

	if err != nil {
		return ""
	}
	return string(data)
}
/*
// encrypt using AES-256 encryption with 32 bytes key
func encrypt(stringToEncrypt string, key string) (encryptedString string) {
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
        logging(err.Error())
		return ""
	}
	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
        logging(err.Error())
		return ""
	}
	nonce := make([]byte, aesGCM.NonceSize())
    _, err = io.ReadFull(rand.Reader, nonce)

	if err != nil {
        logging(err.Error())
		return ""
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return hex.EncodeToString(ciphertext)
}

// decrypt using AES-256 encryption with 32 bytes key
func decrypt(encryptedString string, key string) (decryptedString string) {
	enc, _ := hex.DecodeString(encryptedString)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
        logging(err.Error())
		return ""
	}
	aesGCM, err := cipher.NewGCM(block)

	if err != nil {
        logging(err.Error())
		return ""
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)

	if err != nil {
        logging(err.Error())
		return ""
	}
	return string(plaintext)
}
*/
// return MD5 hash
func encryptPassword(pass string) string {
	hashBinary := md5.Sum([]byte(strings.ToUpper(pass) + "gyZml6bEhv9LFteFk"))
    return hex.EncodeToString(hashBinary[:])
}

/*
func HmacSHA512(input string, keyStr string) string {
	key := []byte(keyStr)
	dst := make([]byte, hex.DecodedLen(len(key)))
	n, err := hex.Decode(dst, key)

	if err != nil {
		return ""
	}
	mac := hmac.New(sha512.New, []byte(string(dst[:n])))
	mac.Write([]byte(input))
	result := mac.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(result))
}
*/
func HmacSHA256(input string, keyStr string) string {
	sig := hmac.New(sha256.New, []byte(keyStr))
	sig.Write([]byte(input))

	return encodeBase64url(sig.Sum(nil))
}

func SH1(input string) string {
	hash := sha1.New()
	io.WriteString(hash, input)

	return hex.EncodeToString(hash.Sum(nil))
}

// generate a pseudo-UUID
func generateUUID() string {
	buff := make([]byte, 16)
	_, err := rand.Read(buff)

	if err != nil {
	    logging(err.Error())
		return ""
	}
	return hex.EncodeToString(buff)
}
