package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

const KeySize = 2048

// Generates private and public key
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

// Convert private key to bytes (pem)
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(priv)
}

// Convert public key to bytes (pem)
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(pub)
}

func FileBytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	privPem, _ := pem.Decode(priv)
	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New(fmt.Sprintf("RSA private key is of the wrong type :%s", privPem.Type))
	}

	return BytesToPrivateKey(privPem.Bytes)
}

// Convert bytes (pem) to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(priv)
}

func FileBytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	pubPem, _ := pem.Decode(pub)
	if pubPem.Type != "RSA PUBLIC KEY" {
		return nil, errors.New(fmt.Sprintf("RSA public key is of the wrong type :%s", pubPem.Type))
	}

	return BytesToPublicKey(pubPem.Bytes)
}

// Convert bytes (pem) to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	return x509.ParsePKCS1PublicKey(pub)
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}

// Encrypt message using public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()

	// Chunk the message into smaller parts
	var chunkSize = pub.N.BitLen()/8 - 2*hash.Size() - 2
	var result []byte
	chunks := chunkBy[byte](msg, chunkSize)
	for _, chunk := range chunks {
		ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, chunk, nil)
		if err != nil {
			return []byte{}, err
		}
		result = append(result, ciphertext...)
	}

	return result, nil
}

// Decrypt message using private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	dec_msg := []byte("")

	for _, chnk := range chunkBy[byte](ciphertext, priv.N.BitLen()/8) {
		plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, chnk, nil)
		if err != nil {
			return []byte{}, err
		}
		dec_msg = append(dec_msg, plaintext...)
	}

	return dec_msg, nil
}

// Example
/*
func main() {
	msg := "Hello, world!"                             // Message to encryption
	PrKey, PubKey, _ := GenerateKeyPair(KeySize)       // Private and Public key generated
	prkstr := PrivateKeyToBytes(PrKey)                 // Private and Public key generated
	fmt.Println("PrivateKey: ", string(prkstr))        // Print serialized private key
	pbkstr, _ := PublicKeyToBytes(PubKey)              // Serialize public key
	fmt.Println("PublicKey: ", string(pbkstr))         // Print serialized public key
	dePub, _ := BytesToPublicKey(pbkstr)               // Deserialize public key
	enc, _ := EncryptWithPublicKey([]byte(msg), dePub) // Encrypt message with deserialized public key
	fmt.Println("Encrypted: ", string(enc))            // Show encrypted message
	dePrv, _ := BytesToPrivateKey(prkstr)              // Show encrypted message
	dec, _ := DecryptWithPrivateKey(enc, dePrv)        // Decrypt message with deserialized private key
	fmt.Println("Decrypted: ", string(dec))            // Show decrypted key
}
*/
