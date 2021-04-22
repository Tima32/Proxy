package main

import "math/rand"

//GeneratePublicPrivateKey - generate public & private key
func GeneratePublicPrivateKey() (privateKey, publicKey uint64) {
	privateKey, publicKey = GeneratePublicPrivateKeyCustom(modulusC, baseC)
	return
}

//GeneratePublicPrivateKeyCustom - generate custom public & private key
func GeneratePublicPrivateKeyCustom(modulus, base uint32) (privateKey, publicKey uint64) {
	//var private_key, public_key uint64

	privateKey = rand.Uint64()
	publicKey = compute(uint64(base), uint64(privateKey), uint64(modulus))
	return
}

//FindSucretKey - Find Sucret key
func FindSucretKey(publicKeyB, privateKeyA uint64) uint32 {
	return uint32(compute(publicKeyB, privateKeyA, modulusC))
}

//FindSucretKeyCustom - Find Sucret key
func FindSucretKeyCustom(publicKeyB, privateKeyA uint64, modulus uint32) uint32 {
	return uint32(compute(publicKeyB, privateKeyA, uint64(modulus)))
}

//------------------[private]---------------

func compute(a, m, n uint64) uint64 {
	var r, y uint64
	y = 1

	for m > 0 {
		r = m % 2

		//fast exponention
		if r == 1 {
			y = (y * a) % n
		}
		a = a * a % n
		m = m / 2
	}
	return y
}

const modulusC = 2375746586
const baseC = 575866576
