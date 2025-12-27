package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"math/big"
)

var (
	otpCharset = "0123456789"
	maxOTP     = big.NewInt(int64(len(otpCharset)))
)

// RandOTP generates a random OTP of given length
func RandOTP(length int) string {
	otp := make([]byte, length)
	for i := range length {
		num, err := rand.Int(rand.Reader, maxOTP)
		if err != nil {
			panic(err)
		}
		otp[i] = otpCharset[num.Int64()]
	}
	return string(otp)
}

// HashOTP hashes the given OTP using SHA-256 returning its hex representation
func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(hash[:])
}

// CompareOTPHash compares the given OTP with its hash
func CompareOTPHash(otp, hexHash string) bool {
	actualHash, err := hex.DecodeString(hexHash)
	if err != nil {
		return false
	}
	otpHash := sha256.Sum256([]byte(otp))
	return subtle.ConstantTimeCompare(otpHash[:], actualHash) == 1 // prevent timing attacks
}
