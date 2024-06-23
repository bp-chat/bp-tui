package main

import (
	"crypto/ed25519"
	"fmt"
	"log"
	"strings"
)

func test() {
	pbka, pka, siga := createKeys("eVmg28qWmcsNLvlW3BkRfNeco47GD49m")
	pbkb, pkb, sigb := createKeys("clMRAq1NcKOXfK10tohluvYgqsXFA3mI")
	log.Printf("Alice pvk:\n%x\n\n", pka)
	log.Printf("Alice puk:\n%x\n\n", pbka)
	log.Printf("Alice sig:\n%x\n\n", siga)
	fmt.Printf("\n\n\n")
	log.Printf("Bob pvk:\n%x\n\n", pkb)
	log.Printf("Bob puk:\n%x\n\n", pbkb)
	log.Printf("Bob sig:\n%x\n\n", sigb)
	fmt.Printf("\n\n\n")
}

func createKeys(seed string) (ed25519.PublicKey, ed25519.PrivateKey, []byte) {
	seedReader := strings.NewReader(seed)
	pbka, pka, err := ed25519.GenerateKey(seedReader)
	if err != nil {
		log.Fatalf("Could not create key\n%v", err)
	}
	sig := ed25519.Sign(pka, pbka)
	return pbka, pka, sig
}
