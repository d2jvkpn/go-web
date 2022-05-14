package misc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

///
type Output struct {
	Key    string `mapstructure:"key"`    // key in json
	Header string `mapstructure:"header"` // header key
	Tmpl   string `mapstructure:"tmpl"`   // header value template
}

type RequestTmpl struct {
	Name            string            `mapstructure:"name"`
	Method          string            `mapstructure:"method"`
	Path            string            `mapstructure:"path"`
	Params          map[string]string `mapstructure:"params"`
	Body            string            `mapstructure:"body"`
	NoPublicHeaders bool              `mapstructure:"no_public_headers"`
	Headers         map[string]string `mapstructure:"headers"`
	Outputs         []Output          `mapstructure:"outputs"`
}

func (output *Output) Get(bts []byte) (value string, err error) {

	value = gjson.GetBytes(bts, output.Key).Str
	if output.Tmpl != "" {
		value = strings.Replace(output.Tmpl, "{}", value, -1)
	}

	return value, nil
}

///
type RequestTmpls struct {
	Url     string            `mapstructure:"url"`
	Headers map[string]string `mapstructure:"headers"`
	Prelude []RequestTmpl     `mapstructure:"prelude"`
	APIs    []RequestTmpl     `mapstructure:"apis"`
	// header  http.Header
	client *http.Client
	apiMap map[string]*RequestTmpl
}

func LoadRequestTmpls(name, fp string) (item *RequestTmpls, err error) {
	conf := viper.New()
	conf.SetConfigName(name)
	conf.SetConfigFile(fp)
	if err = conf.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ReadInConfig(): %q, %v", fp, err)
	}

	item = new(RequestTmpls)
	if err = conf.Unmarshal(item); err != nil {
		return nil, err
	}

	item.client = new(http.Client)
	// item.header = make(http.Header)
	item.apiMap = make(map[string]*RequestTmpl, len(item.APIs))
	for i := range item.APIs {
		api := &item.APIs[i]
		m := strings.ToUpper(api.Method)
		if api.Method != "GET" && api.Method != "POST" {
			return nil, fmt.Errorf("invalid method: %s", api.Method)
		}
		api.Method = m

		item.apiMap[api.Name] = api
	}

	//	for k, v := range item.Headers {
	//		item.header.Set(k, v)
	//	}

	return
}

func (item *RequestTmpls) request(tmpl *RequestTmpl, prelude bool) (
	statusCode int, body string, err error) {
	var (
		bts    []byte
		v, p   string
		out    bytes.Buffer
		reader io.Reader
		req    *http.Request
		res    *http.Response
	)

	p = item.Url + tmpl.Path
	if len(tmpl.Params) > 0 {
		strs := make([]string, 0, len(tmpl.Params))
		for k, v := range tmpl.Params {
			strs = append(strs, fmt.Sprintf("%s=%s", k, v))
		}
		p += "?" + strings.Join(strs, "&")
	}

	if len(tmpl.Body) > 0 {
		reader = bytes.NewBufferString(tmpl.Body)
	}

	if req, err = http.NewRequest(tmpl.Method, p, reader); err != nil {
		return 0, "", err
	}

	if !tmpl.NoPublicHeaders {
		for k, v := range item.Headers {
			req.Header.Add(k, v)
		}
	}
	for k, v := range tmpl.Headers {
		req.Header.Add(k, v)
	}

	if res, err = item.client.Do(req); err != nil {
		return 0, "", err
	}
	statusCode = res.StatusCode
	defer res.Body.Close()

	bts, err = io.ReadAll(res.Body)

	isJSON := strings.Contains(res.Header.Get("Content-Type"), "application/json")
	if len(bts) > 0 && isJSON {
		if e := json.Indent(&out, bts, "", "  "); e == nil {
			body = string(out.Bytes())
		} else {
			body = string(bts)
		}
	} else {
		body = string(bts)
	}

	if err != nil {
		return
	}
	if !prelude {
		return
	}

	for i := range tmpl.Outputs {
		output := &tmpl.Outputs[i]
		if v, err = output.Get(bts); err != nil {
			return
		}
		// item.header.Set(output.Header, v)
		item.Headers[output.Header] = v
	}

	return
}

func (item *RequestTmpls) Request(name string) (
	statusCode int, body string, err error) {
	var (
		ok   bool
		tmpl *RequestTmpl
	)

	if tmpl, ok = item.apiMap[name]; !ok {
		return -3, "", fmt.Errorf("api not found")
	}

	for i := range item.Prelude {
		r := &item.Prelude[i]
		if statusCode, _, err = item.request(r, true); err != nil {
			return -2, "", err
		}
		if statusCode != http.StatusOK {
			return -1, "", fmt.Errorf("%s statusCode: %d", r.Name, statusCode)
		}
	}

	return item.request(tmpl, false)
}
