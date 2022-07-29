package wrap

import (
	"crypto/rsa"
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v4"
)

///
type JwtHSAuth struct {
	key    []byte
	method *jwt.SigningMethodHMAC // SigningMethodHS{256,384,512}
}

func NewHSAuth(key string, code uint) (auth *JwtHSAuth, err error) {
	auth = &JwtHSAuth{key: []byte(key)}

	switch code {
	case 256:
		auth.method = jwt.SigningMethodHS256
	case 384:
		auth.method = jwt.SigningMethodHS384
	case 512:
		auth.method = jwt.SigningMethodHS512
	default:
		return nil, fmt.Errorf("invalid code")
	}

	return auth, nil
}

func (auth *JwtHSAuth) Sign(data map[string]any) (str string, err error) {
	var (
		token  *jwt.Token
		claims jwt.MapClaims
	)

	claims = make(jwt.MapClaims, len(data))
	for k, v := range data { // TODO: can't do type assertion, how to avoid copy
		claims[k] = v
	}

	token = jwt.NewWithClaims(auth.method, claims)
	return token.SignedString(auth.key)
}

func (auth *JwtHSAuth) Parse(str string) (data map[string]any, err error) {
	var (
		ok     bool
		token  *jwt.Token
		claims jwt.MapClaims
	)

	// options ...ParserOption
	token, err = jwt.Parse(str, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.key, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	data = make(map[string]any, len(claims))
	for k, v := range claims {
		data[k] = v
	}

	return data, nil
}

///
type JwtRSAAuth struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	method     *jwt.SigningMethodRSA // SigningMethodRSA{256,384,512}
}

func NewRSAAuth(privateKeyFile, publicKeyFile string, code uint) (auth *JwtRSAAuth, err error) {
	var bts []byte

	auth = &JwtRSAAuth{}
	if bts, err = os.ReadFile(privateKeyFile); err != nil {
		return nil, err
	}
	if auth.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(bts); err != nil {
		return nil, err
	}

	if bts, err = os.ReadFile(publicKeyFile); err != nil {
		return nil, err
	}
	if auth.publicKey, err = jwt.ParseRSAPublicKeyFromPEM(bts); err != nil {
		return nil, err
	}

	switch code {
	case 256:
		auth.method = jwt.SigningMethodRS256
	case 384:
		auth.method = jwt.SigningMethodRS384
	case 512:
		auth.method = jwt.SigningMethodRS512
	default:
		return nil, fmt.Errorf("invalid code")
	}

	return auth, nil
}

func (auth *JwtRSAAuth) Sign(data map[string]any) (str string, err error) {
	var (
		token  *jwt.Token
		claims jwt.MapClaims
	)

	claims = make(jwt.MapClaims, len(data))
	for k, v := range data { // TODO: can't do type assertion, how to avoid copy
		claims[k] = v
	}

	token = jwt.NewWithClaims(auth.method, claims)
	return token.SignedString(auth.privateKey)
}

func (auth *JwtRSAAuth) Parse(str string) (data map[string]any, err error) {
	var (
		ok     bool
		token  *jwt.Token
		claims jwt.MapClaims
	)

	// options ...ParserOption
	token, err = jwt.Parse(str, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.publicKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	if claims, ok = token.Claims.(jwt.MapClaims); !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	data = make(map[string]any, len(claims))
	for k, v := range claims {
		data[k] = v
	}

	return data, nil
}
