package main

import (
	"C"
	"aidanwoods.dev/go-paseto"
	"time"
)

//export createToken
func createToken(data *C.char) *C.char {
	hexData := "eedc7728fd313726d753ed892f5f9129a9eadbedfce193b6d0fd9e923712a957"

	symK, _ := paseto.V4SymmetricKeyFromHex(hexData)

	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetIssuer("CryptocurrencyBot")
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	token.SetString("data", C.GoString(data))

	encrypted := token.V4Encrypt(symK, nil)

	return C.CString(encrypted)
}

func main() {

}
