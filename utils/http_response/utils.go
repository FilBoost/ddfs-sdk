package http_response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	logging "github.com/ipfs/go-log/v2"
)

var log = logging.Logger("http_response")

func ProcessRequest(req *http.Request, timeout time.Duration, decode func(req *http.Request, resp *http.Response) error) error {
	cli := http.Client{}
	cli.Timeout = timeout
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	var ignoreClosingBody bool
	defer func() {
		if !ignoreClosingBody {
			_err := resp.Body.Close()
			if _err != nil {
				log.Warnf("[DD] close body failed,access:%v, err:%v", req.URL, _err)
			}
		}
	}()

	if resp.StatusCode >= 500 {
		err = New500Err(resp.StatusCode, resp.Status, req.URL)
		return err
	}

	if 400 <= resp.StatusCode && resp.StatusCode < 500 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("[DD] %v:%v access:%v, err:%v", resp.Status, resp.StatusCode, req.URL, err)
			return fmt.Errorf("%v:%v access:%v, err:%v", resp.Status, resp.StatusCode, req.URL, err)
		}

		if err, ok := UnmarshalResponseErr(body); ok {
			return err
		}

		log.Errorf("[DD] %v:%v access:%v, body: %v", resp.Status, resp.StatusCode, req.URL, string(body))
		return fmt.Errorf("%v:%v access:%v, body: %v", resp.Status, resp.StatusCode, req.URL, string(body))
	}

	if decode == nil {
		return nil
	}

	ignoreClosingBody = true
	err = decode(req, resp)
	if err != nil {
		return err
	}
	return nil
}

func ProcessRequestAndDecodeResponse(req *http.Request, timeout time.Duration, responseInfo interface{}) error {
	return ProcessRequest(req, timeout, func(req *http.Request, resp *http.Response) error {
		defer func() {
			_err := resp.Body.Close()
			if _err != nil {
				log.Warnf("[DD] close body failed,access:%v, err:%v", req.URL, _err)
			}
		}()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("[DD] read body failed err:%v", err)
			return err
		}

		var resInfo ResponseInfo
		err = json.Unmarshal(body, &resInfo)
		if err != nil {
			log.Errorf("[DD] unmarshal http body failed err:%v", err)
			return err
		}

		if !IsOk(resInfo.Code) {
			return fmt.Errorf("[DD] information of reponse: %v", string(body))
		}

		err = json.Unmarshal(resInfo.Data, responseInfo)
		if err != nil {
			log.Errorf("[DD] unmarshal http data failed err:%v", err)
			return err
		}
		return nil
	})
}

func ProcessRequestAndIgnoreResponse(req *http.Request, timeout time.Duration) error {
	return ProcessRequest(req, timeout, nil)
}

func MakeRequest(method string, host string, _url string, rawQuery string, data []byte) (*http.Request, error) {
	hUrl, err := url.ParseRequestURI(host)
	if err != nil {
		log.Errorf("[DD] parse url failed url: %v,err:%v", host, err)
		return nil, err
	}

	goodUrl, err := hUrl.Parse(_url)
	if err != nil {
		log.Errorf("[DD] parse url failed url: %v,err:%v", _url, err)
		return nil, err
	}

	goodUrl.RawQuery = rawQuery
	var body io.Reader
	if data != nil {
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, goodUrl.String(), body)
	if err != nil {
		log.Errorf("[DD] NewRequest failed err:%v", err)
		return nil, err
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	return req, nil
}
