package impl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"ddfs-sdk/api"
	"ddfs-sdk/utils/http_response"
)

var log = logging.Logger("ddfs-apiImpl")

var _ api.RemoteFileOpt = (*DDFileOpt)(nil)

type DDFileOpt struct {
	file string
	host string
}

const FetchFileUrl = "/api/v0/file_opt/fetch"

const RevertFileStateUrl = "/api/v0/file_opt/revert"

const ConfirmFileUrl = "/api/v0/file_opt/confirm"

func (fo *DDFileOpt) FetchWithConfirm() (io.ReadCloser, uint64, error) {
	return fo.fetch(false)
}

func (fo *DDFileOpt) Fetch() (io.ReadCloser, uint64, error) {
	return fo.fetch(true)
}

func (fo *DDFileOpt) fetch(offerConfirmation bool) (io.ReadCloser, uint64, error) {
	v := url.Values{}
	v.Add("file", fo.file)
	v.Add("offer_confirmation", strconv.FormatBool(offerConfirmation))

	req, err := http_response.MakeRequest(http.MethodGet, fo.host, FetchFileUrl, v.Encode(), nil)
	if err != nil {
		log.Errorf("[DD] MakeRequest failed err: %v", err)
		return nil, 0, err
	}

	var body io.ReadCloser
	var contentLength int64
	err = http_response.ProcessRequest(req, time.Hour*24, func(req *http.Request, resp *http.Response) error {
		if resp.ContentLength < 0 {
			resp.Body.Close()
			return fmt.Errorf("http can't get known ContentLength,access:%v", req.URL)
		}
		body = resp.Body
		contentLength = resp.ContentLength
		return nil
	})
	if err != nil {
		log.Errorf("[DD] request http failed err: %v", err)
		return nil, 0, err
	}

	return body, uint64(contentLength), nil
}

func (fo *DDFileOpt) Confirm(key string) error {
	v := url.Values{}
	v.Add("file", fo.file)
	v.Add("key", key)
	req, err := http_response.MakeRequest(http.MethodPut, fo.host, ConfirmFileUrl, v.Encode(), nil)
	if err != nil {
		log.Errorf("[DD] NewRequest failed err: %v", err)
		return err
	}

	err = http_response.ProcessRequestAndIgnoreResponse(req, time.Second*30)
	if err != nil {
		log.Errorf("[DD] request http failed err: %v", err)
		return err
	}

	return nil
}

func (fo *DDFileOpt) Revert() error {
	v := url.Values{}
	v.Add("file", fo.file)
	req, err := http_response.MakeRequest(http.MethodPut, fo.host, RevertFileStateUrl, v.Encode(), nil)
	if err != nil {
		log.Errorf("[DD] NewRequest failed err: %v", err)
		return err
	}

	err = http_response.ProcessRequestAndIgnoreResponse(req, time.Second*30)
	if err != nil {
		log.Errorf("[DD] request http failed err: %v", err)
		return err
	}

	return nil
}

func NewFileOpt(remoteFileUrl string) (api.RemoteFileOpt, error) {
	ss := strings.Split(remoteFileUrl, "|")
	if len(ss) != 2 {
		return nil, fmt.Errorf("unknown remoteFileUrl: %v", remoteFileUrl)
	}

	return &DDFileOpt{
		host: ss[0],
		file: ss[1],
	}, nil
}
