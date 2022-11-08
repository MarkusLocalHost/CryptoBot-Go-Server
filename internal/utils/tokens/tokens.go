package tokens

import (
	"aidanwoods.dev/go-paseto"
	"errors"
	"time"
)

func VerifyToken(token string) (string, error) {
	hexData := "eedc7728fd313726d753ed892f5f9129a9eadbedfce193b6d0fd9e923712a957"
	symK, _ := paseto.V4SymmetricKeyFromHex(hexData)

	parser := paseto.NewParserWithoutExpiryCheck()
	decrypted, err := parser.ParseV4Local(symK, token, nil)
	if err != nil {
		return "", err
	}

	// todo issuer notbefore issuedat

	// check time
	expString, err := decrypted.GetString("exp")
	if err != nil {
		return "", err
	}
	expTime, err := time.Parse(time.RFC3339, expString)
	if err != nil {
		return "", err
	}
	if time.Now().After(expTime) {
		return "", errors.New("expired token")
	}

	// get data
	jsonString, err := decrypted.GetString("data")
	if err != nil {
		return "", err
	}

	return jsonString, nil
}
