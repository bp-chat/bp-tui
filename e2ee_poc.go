package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type actorKeys struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	preKey     *ecdh.PrivateKey
	ek         *ecdh.PrivateKey
	signature  []byte
}

func test() {
	//TODO understand GCM
	alice := createKeys()
	bob := createKeys()
	textMsg := "hello bob, how are you?"

	isValidSignature := ed25519.Verify(bob.publicKey, bob.preKey.PublicKey().Bytes(), bob.signature)
	if isValidSignature == false {
		fmt.Println("Could not verify signature")
	}

	adhe, _ := alice.ek.ECDH(bob.preKey.PublicKey())
	adhi, _ := alice.preKey.ECDH(bob.ek.PublicKey())
	ask := sha256.Sum256(append(adhe, adhi...))

	abc, _ := aes.NewCipher(ask[:])
	ac, _ := cipher.NewGCM(abc)
	aiv := make([]byte, 12)
	io.ReadFull(rand.Reader, aiv)
	ciphertext := ac.Seal(nil, aiv, []byte(textMsg), nil)
	fmsg := append(aiv, ciphertext...)

	bdhe, _ := bob.ek.ECDH(alice.preKey.PublicKey())
	bdhi, _ := bob.preKey.ECDH(alice.ek.PublicKey())
	// I don't like that it's necessary to preserve order but it's ok... I guess...
	bsk := sha256.Sum256(append(bdhi, bdhe...))
	bbc, _ := aes.NewCipher(bsk[:])

	bc, _ := cipher.NewGCM(bbc)
	riv := fmsg[:12]
	dmsg, _ := bc.Open(nil, riv, fmsg[12:], nil)
	fmt.Println("bob decipher")
	fmt.Println(string(dmsg))

	fmt.Println("ok for now")
}

func createKeys() actorKeys {
	curve := ecdh.X25519()
	pubkey, pv, _ := ed25519.GenerateKey(rand.Reader)
	preKey, _ := curve.GenerateKey(rand.Reader)
	ek, _ := curve.GenerateKey(rand.Reader)
	sinature := ed25519.Sign(pv, preKey.PublicKey().Bytes())
	return actorKeys{
		pv,
		pubkey,
		preKey,
		ek,
		sinature,
	}
}
