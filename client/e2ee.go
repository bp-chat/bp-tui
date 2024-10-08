package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

const ivSize = 12

type InitialVector = [ivSize]byte

type KeySet struct {
	PrivateKey ed25519.PrivateKey //Private component of identity key
	PublicKey  ed25519.PublicKey  // public component of identity key
	PreKey     *ecdh.PrivateKey   // This pre signed key is a temporary identity key
	Ek         *ecdh.PrivateKey   // this ephemeral key should be created for each exchange
	Signature  []byte             // signature, used to authenticate the preKey.
}

type PublicKeySet struct {
	IdentityKey  []byte
	SignedKey    []byte
	Ephemeralkey []byte
	Signature    []byte
}

func test() {
	//TODO understand GCM https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//for some reason when using this mode we don't need to mac or data
	alice := CreateKeys()
	alicePublic := convertKeys(alice)
	bob := CreateKeys()
	bobPublic := convertKeys(bob)
	textMsg := "hello bob, how are you?"

	ask, _ := CreateSharedKey(alice, bobPublic)

	iv, fmsg, _ := Encrypt(ask[:], []byte(textMsg))

	bsk, _ := CreateSharedKey(bob, alicePublic)
	dmsg, _ := Decrypt(bsk[:], iv, fmsg)
	fmt.Println("bob decipher")
	fmt.Println(string(dmsg))

	fmt.Println("ok for now")
}

func CreateKeys() KeySet {
	curve := ecdh.X25519()
	pubkey, pv, _ := ed25519.GenerateKey(rand.Reader)
	preKey, _ := curve.GenerateKey(rand.Reader)
	ek, _ := curve.GenerateKey(rand.Reader)
	sinature := ed25519.Sign(pv, preKey.PublicKey().Bytes())
	return KeySet{
		pv,
		pubkey,
		preKey,
		ek,
		sinature,
	}
}

func isValidKey(keys PublicKeySet) bool {
	return ed25519.Verify(keys.IdentityKey, keys.SignedKey, keys.Signature)
}

func convertKeys(other KeySet) PublicKeySet {
	return PublicKeySet{
		IdentityKey:  other.PublicKey,
		SignedKey:    other.PreKey.PublicKey().Bytes(),
		Ephemeralkey: other.Ek.PublicKey().Bytes(),
		Signature:    other.Signature,
	}
}

func CreateSharedKey(own KeySet, other PublicKeySet) ([sha256.Size]byte, error) {
	empty := [sha256.Size]byte{}
	curve := ecdh.X25519()
	if isValidKey(other) == false {
		return empty, errors.New("Invalid signed key")
	}

	otherIK, err := curve.NewPublicKey(other.SignedKey)
	if err != nil {
		return empty, err
	}
	otherEK, err := curve.NewPublicKey(other.Ephemeralkey)
	if err != nil {
		return empty, err
	}

	dhi, err := own.PreKey.ECDH(otherEK)
	if err != nil {
		return empty, err
	}
	dhe, err := own.Ek.ECDH(otherIK)
	if err != nil {
		return empty, err
	}
	combined := make([]byte, len(dhe))
	for i := 0; i < len(dhi); i++ {
		combined[i] = dhi[i] ^ dhe[i] //TODO find a better way to combine both
	}
	return sha256.Sum256(combined), nil
}

// returns the Initial Vector and the encrypted message
func Encrypt(sharedKey []byte, message []byte) (InitialVector, []byte, error) {
	empty := []byte{}
	aesBlock, err := aes.NewCipher(sharedKey[:])
	if err != nil {
		return [ivSize]byte{}, empty, err
	}
	gcmCipher, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return [ivSize]byte{}, empty, err
	}
	initialVector := [ivSize]byte{}
	io.ReadFull(rand.Reader, initialVector[:])
	ciphertext := gcmCipher.Seal(nil, initialVector[:], message, nil)
	return initialVector, ciphertext, nil
}

func Decrypt(sharedKey []byte, initialVector InitialVector, message []byte) ([]byte, error) {
	empty := []byte{}
	aesBlock, err := aes.NewCipher(sharedKey[:])
	if err != nil {
		return empty, err
	}
	gcmCipher, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return empty, err
	}
	return gcmCipher.Open(nil, initialVector[:], message, nil)
}
