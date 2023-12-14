package main

import "unicode"

func verifyPasswordStrength(password string) bool {
	/*
		- Should have minimum of 8 characters
		- Should have atleast 1 number
		- Should have atleast 1 symbol
		- Should have atleast 1 uppercase letter
		- Should have atleast 1 lowercase letter
	*/
	hasMin8Chars := len(password) >= 8
	hasNumber := false
	hasSymbol := false
	hasUppercase := false
	hasLowercase := false

	for _, c := range password {
		if unicode.IsNumber(c) {
			hasNumber = true
		}
		if !unicode.IsNumber(c) && !unicode.IsLetter(c) {
			hasSymbol = true
		}
		if unicode.IsUpper(c) {
			hasUppercase = true
		}
		if unicode.IsLower(c) {
			hasLowercase = true
		}
	}

	return hasMin8Chars && hasNumber && hasSymbol && hasUppercase && hasLowercase
}
