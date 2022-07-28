package wrap

import (
	"fmt"
	"testing"

	jwt "github.com/golang-jwt/jwt/v4"

	. "github.com/stretchr/testify/require"
)

func TestJwt(t *testing.T) {
	fmt.Println(">> TestJwt:", jwt.GetAlgorithms())
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
