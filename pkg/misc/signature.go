package misc

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type SigningUrlMd5 struct {
	secrete, key string
	lowcase      bool
}

func NewSigningUrlMd5(secrete, key string, lowcase bool) SigningUrlMd5 {
	return SigningUrlMd5{secrete: secrete, key: key, lowcase: lowcase}
}

func (sign *SigningUrlMd5) signValue(param map[string]string) (value string) {
	var (
		k, format string
		keys      []string
		pairs     []string
	)

	keys = make([]string, 0, len(param))

	for k = range param {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	pairs = make([]string, 0, len(param))
	for _, k = range keys {
		pairs = append(pairs, k+param[k])
	}

	if format = "%X"; sign.lowcase {
		format = "%x"
	}

	value = fmt.Sprintf(format, md5.Sum([]byte(sign.secrete+strings.Join(pairs, "")+sign.secrete)))

	return value
}

func (sign *SigningUrlMd5) Sign(param map[string]string) (query string) {
	var (
		value string
		pairs []string
	)

	value = sign.signValue(param)
	pairs = make([]string, 0, len(param)+1)

	for k, v := range param {
		pairs = append(pairs, url.QueryEscape(k)+"="+url.QueryEscape(v))
	}
	pairs = append(pairs, url.QueryEscape(sign.key)+"="+url.QueryEscape(value))

	return strings.Join(pairs, "&")
}

func (sign *SigningUrlMd5) Verify(query string) (err error) {
	var (
		value  string
		param  map[string]string
		values url.Values
	)

	if values, err = url.ParseQuery(query); err != nil {
		return err
	}

	if value = values.Get(sign.key); value == "" || len(value) != 32 {
		return fmt.Errorf("invalid signature")
	}

	param = make(map[string]string, len(values)-1)

	for k := range values {
		if k == sign.key {
			continue
		}
		if len(values[k]) == 0 {
			param[k] = ""
		} else {
			param[k] = values[k][0]
		}
	}

	if value != sign.signValue(param) {
		return fmt.Errorf("signature not match")
	}

	return nil
}
