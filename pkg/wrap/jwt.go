package wrap

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

type JwtHSAuth struct {
	key    []byte
	method *jwt.SigningMethodHMAC // SigningMethodHS{256,384,512}
}

func NewHSAuth(key string, hsCode uint) (auth *JwtHSAuth, err error) {
	auth = &JwtHSAuth{key: []byte(key)}

	switch hsCode {
	case 256:
		auth.method = jwt.SigningMethodHS256
	case 384:
		auth.method = jwt.SigningMethodHS384
	case 512:
		auth.method = jwt.SigningMethodHS512
	default:
		return nil, fmt.Errorf("invalid hsCode")
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

func GinJwtHSAuth(auth *JwtHSAuth, handle func(*gin.Context, map[string]any) error,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			head string
			data map[string]any
			err  error
		)

		if head = ctx.Request.Header.Get("Authorization"); head == "" {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !strings.HasPrefix(head, "Bearer ") {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		if data, err = auth.Parse(head[7:]); err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if handle != nil {
			if err = handle(ctx, data); err != nil {
				ctx.Abort()
				return
			}
			// for k, v := range data { ctx.Set(k, v) }
		}

		ctx.Next()
	}
}
