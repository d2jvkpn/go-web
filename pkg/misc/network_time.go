package misc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// https://en.wikipedia.org/wiki/Network_Time_Protocol
// https://en.wikipedia.org/wiki/File:NTP-Algorithm.svg
type NetworkTimeServer struct {
	delay int64
	*http.Server
}

func (ser *NetworkTimeServer) Run() error {
	return ser.ListenAndServe()
}

func NewNetworkTimeServer(addr string, delay int64) (ser *NetworkTimeServer, err error) {
	if delay < 0 {
		return nil, fmt.Errorf("invalid delay")
	}
	if delay == 0 {
		delay = 10
	}
	ser = new(NetworkTimeServer)
	ser.delay = delay

	mux := http.NewServeMux()

	mux.HandleFunc("/", NTSFunc(delay))

	ser.Server = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 10, // 8M
	}

	return ser, nil
}

func NTSFunc(delay int64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			d   int64
			err error
		)

		d = delay
		data := make(map[string]int64, 2)
		data["t2"] = time.Now().UnixMilli()

		if delayStr := r.URL.Query().Get("delay"); delayStr != "" {
			if d, err = strconv.ParseInt(delayStr, 10, 64); err != nil {
				w.WriteHeader(http.StatusBadRequest)

				json.NewEncoder(w).Encode(
					map[string]interface{}{"code": -1, "msg": "bad request", "data": data},
				)
				return
			}
		}

		// fmt.Println(">>>", d)
		time.Sleep(time.Duration(d) * time.Millisecond)
		data["t3"] = time.Now().UnixMilli()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(
			map[string]interface{}{"code": 0, "msg": "ok", "data": data},
		)
	}
}

type NetworkTimeResult struct {
	T1    time.Time `json:"t1"`
	T2    time.Time `json:"t2"`
	T3    time.Time `json:"t3"`
	T4    time.Time `json:"t4"`
	Sigma int64     `json:"sigma"`
	Delta int64     `json:"delta"`
}

func (result NetworkTimeResult) String() string {
	bts, _ := json.MarshalIndent(result, "", "  ")
	return string(bts)
}

func GetNetworkTime(addr string, delay int64) (result *NetworkTimeResult, err error) {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}

	address := fmt.Sprintf("%s?delay=%d", addr, delay)
	result = new(NetworkTimeResult)
	result.T1 = time.Now()
	var res *http.Response
	if res, err = http.Get(address); err != nil {
		return nil, fmt.Errorf("http.Get(%q): %w", address, err)
	}
	defer res.Body.Close()
	result.T4 = time.Now()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response.StatusCode: %d", res.StatusCode)
	}

	data := struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			T2 int64 `json:"t2"` // unix millisec
			T3 int64 `json:"t3"` // unix millisec
		} `json:"data"`
	}{}

	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode response.Body: %w", err)

	}

	result.T2, result.T3 = time.UnixMilli(data.Data.T2), time.UnixMilli(data.Data.T3)

	d1 := result.T4.UnixMilli() - result.T1.UnixMilli()
	d2 := result.T3.UnixMilli() - result.T2.UnixMilli()
	result.Sigma = (d1 - d2) / 2

	result.Delta = result.T2.UnixMilli() - result.T1.UnixMilli() - result.Sigma
	return result, nil
}
