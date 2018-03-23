package main

import (
	"encoding/base64"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
)

func main() {
	appId := "wx4f4bc4dec97d474b"

    sessionKey := "tiihtNczf5v6AKRyjwEUhQ=="
    encryptedData := "CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZMQmRzooG2xrDcvSnxIMXFufNstNGTyaGS9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+3hVbJSRgv+4lGOETKUQz6OYStslQ142dNCuabNPGBzlooOmB231qMM85d2/fV6ChevvXvQP8Hkue1poOFtnEtpyxVLW1zAo6/1Xx1COxFvrc2d7UL/lmHInNlxuacJXwu0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn/Hz7saL8xz+W//FRAUid1OksQaQx4CMs8LOddcQhULW4ucetDf96JcR3g0gfRK4PC7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns/8wR2SiRS7MNACwTyrGvt9ts8p12PKFdlqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYVoKlaRv85IfVunYzO0IKXsyl7JCUjCpoG20f0a04COwfneQAGGwd5oa+T8yO5hzuyDb/XcxxmK01EpqOyuxINew=="
    iv := "r7BXXKkLb8qrSNn05n0qiA=="


	base64Sessionkey := make([]byte, base64.StdEncoding.DecodedLen(len(sessionKey)))
	base64.StdEncoding.Decode(base64Sessionkey, []byte(sessionKey))
	fmt.Println("%v", base64Sessionkey)

	base64EncryptData := make([]byte, base64.StdEncoding.DecodedLen(len(encryptedData)))
	base64.StdEncoding.Decode(base64EncryptData, []byte(encryptedData))

	fmt.Println("%v", base64EncryptData)


	base64Iv := make([]byte, base64.StdEncoding.DecodedLen(len(iv)))
	base64.URLEncoding.Decode(base64Iv, []byte(iv))

	fmt.Println("%v", base64Iv)


	//fmt.Printf("base text = %v\n", string(base64Sessionkey))
	str := decrypt(base64Sessionkey, base64EncryptData, base64Iv)

	fmt.Printf("appid = %v\n", appId)
	fmt.Printf("result = %v \n", str)

}


func decrypt(key []byte, securemess, iv []byte) (decodedmess string) {
	//cipherText, err := base64.URLEncoding.DecodeString(securemess)
	//if err != nil {
	//	return
	//}

	block, err := aes.NewCipher(key[:16])
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}

	if len(securemess) < aes.BlockSize {
		fmt.Println("Ciphertext block size is too short!")
		return
	}

	dest := make([]byte, len(securemess))
	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	//iv := cipherText[:aes.BlockSize]
	//cipherText = cipherText[aes.BlockSize:]

	fmt.Println(block.BlockSize())
	stream := cipher.NewCBCDecrypter(block, iv[:16])
	stream.CryptBlocks(dest, securemess[:len(securemess)-2])
	//stream := cipher.NewCFBDecrypter(block, iv[:16])
	// XORKeyStream can work in-place if the two arguments are the same.
	//stream.XORKeyStream(dest, securemess[:len(securemess)-2])

	decodedmess = string(dest)
	return
}