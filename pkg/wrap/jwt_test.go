package wrap

import (
	"crypto/rsa"
	"fmt"
	"os"
	"testing"

	jwt "github.com/golang-jwt/jwt/v4"

	. "github.com/stretchr/testify/require"
)

func TestJwtAlgs(t *testing.T) {
	fmt.Println(">> TestJwtAlgs:", jwt.GetAlgorithms())
}

func TestJwtHSAuth(t *testing.T) {
	var (
		str   string
		err   error
		data  map[string]any
		data2 map[string]any
		auth  *JwtHSAuth
	)

	auth, err = NewHSAuth("123456", 256)
	NoError(t, err)

	data = map[string]any{
		"key1": "value1",
		"key2": 42.24,
	}

	str, err = auth.Sign(data)
	NoError(t, err)
	fmt.Println(">>> signed token:", str)

	data2, err = auth.Parse(str)
	NoError(t, err)
	fmt.Println(">>> parsed token:", data2)

	Equal(t, data["key1"], data2["key1"])
	Equal(t, data["key2"], data2["key2"])
}

// $ openssl req -newkey rsa:2048 -new -nodes -x509 -days 365 -subj "/C=US/ST=New Sweden/L=Stockholm/O=.../OU=.../CN=.../emailAddress=..." -keyout rsa_private.pem -out rsa_public.pem
func TestRSAPem(t *testing.T) {
	var (
		bts        []byte
		sig        string
		err        error
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
	)

	bts, err = os.ReadFile("rsa_private.pem")
	NoError(t, err)
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bts)
	NoError(t, err)
	fmt.Println(">>> privateKey:", privateKey)

	bts, err = os.ReadFile("rsa_public.pem")
	NoError(t, err)
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(bts)
	NoError(t, err)
	fmt.Println(">>> publicKey:", publicKey)

	method := jwt.SigningMethodRS256
	sig, err = method.Sign("abcdefg", privateKey)
	NoError(t, err)
	fmt.Println(">> Signature:", sig)

	err = method.Verify("abcdefg", sig, publicKey)
	NoError(t, err)
}

func TestJwtRSAAuth(t *testing.T) {
	var (
		str   string
		err   error
		data  map[string]any
		data2 map[string]any
		auth  *JwtRSAAuth
	)

	auth, err = NewRSAAuth("rsa_private.pem", "rsa_public.pem", 256)
	NoError(t, err)

	data = map[string]any{
		"key1": "value1",
		"key2": 42.24,
	}

	str, err = auth.Sign(data)
	NoError(t, err)
	fmt.Println(">>> signed token:", str)

	data2, err = auth.Parse(str)
	NoError(t, err)
	fmt.Println(">>> parsed token:", data2)

	Equal(t, data["key1"], data2["key1"])
	Equal(t, data["key2"], data2["key2"])
}
