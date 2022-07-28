package wrap

import (
	"fmt"
	"testing"

	jwt "github.com/golang-jwt/jwt/v4"
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

	if auth, err = NewHSAuth("123456", 256); err != nil {
		t.Fatal(err)
	}

	data = map[string]any{
		"key1": "value1",
		"key2": 42.24,
	}

	if str, err = auth.Sign(data); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> signed token:", str)

	if data2, err = auth.Parse(str); err != nil {
		t.Fatal(err)
	}
	fmt.Println(">>> parsed token:", data2)
}
