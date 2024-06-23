package main

import (
	"crypto/ecdh"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"log"
)

type idKey struct {
	privateKey *ecdh.PrivateKey
	publicKey  *ecdh.PublicKey
	signature  []byte
}

func test() {
	aliceKey := createKeys()
	bobKey := createKeys()
	fmt.Printf("### Alice ###\n")
	fmt.Printf("pvk:%x\n", aliceKey.privateKey.Bytes())
	fmt.Printf("puk:%x\n", aliceKey.publicKey.Bytes())
	fmt.Printf("sig:%x\n", aliceKey.signature)
	fmt.Printf("\n\n### Bob ###\n")
	fmt.Printf("pvk:%x\n", bobKey.privateKey.Bytes())
	fmt.Printf("puk:%x\n", bobKey.publicKey.Bytes())
	fmt.Printf("sig:%x\n", bobKey.signature)
	fmt.Printf("\n\n")

	ska, _ := aliceKey.privateKey.ECDH(bobKey.publicKey)
	skb, _ := bobKey.privateKey.ECDH(aliceKey.publicKey)
	fmt.Printf("Alice shared:\t\t%x\n", ska)
	fmt.Printf("Bob shared:\t\t%x\n", skb)
}

func createKeys() idKey {
	curve := ecdh.X25519()
	privateKey, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Could not create key\n%v", err)
	}
	publicKey := privateKey.PublicKey()

	sig := ed25519.Sign(append(privateKey.Bytes(), publicKey.Bytes()...), publicKey.Bytes())
	return idKey{
		privateKey,
		publicKey,
		sig,
	}
}
