package addpiece_ext

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"
)

var IllegalExtInfoFormat = errors.New("illegal DD_Addpiece_Param format")

type DD_Addpiece_Param struct {
	PieceCid      cid.Cid
	RemoteFileUrl string
}

func (ei *DD_Addpiece_Param) String() string {
	return fmt.Sprintf("DD_EXT_INFO|%s|%s", ei.PieceCid.String(), ei.RemoteFileUrl)
}

func Cast(s string) (DD_Addpiece_Param, error) {
	s, found := strings.CutPrefix(s, "DD_EXT_INFO|")
	if !found {
		return DD_Addpiece_Param{}, IllegalExtInfoFormat
	}

	_pieceCid, remoteFileUrl, found := strings.Cut(s, "|")
	if !found {
		return DD_Addpiece_Param{}, IllegalExtInfoFormat
	}

	pieceCid, err := cid.Decode(_pieceCid)
	if err != nil {
		return DD_Addpiece_Param{}, IllegalExtInfoFormat
	}

	return DD_Addpiece_Param{pieceCid, remoteFileUrl}, nil
}

func GetExtInfoFromCtx(ctx context.Context) (DD_Addpiece_Param, error) {
	var info DD_Addpiece_Param
	if _pieceCid, ok := ctx.Value("pieceCID").(string); ok {
		pieceCid, err := cid.Decode(_pieceCid)
		if err != nil {
			return DD_Addpiece_Param{}, IllegalExtInfoFormat
		}

		info.PieceCid = pieceCid
	} else {
		return DD_Addpiece_Param{}, IllegalExtInfoFormat
	}

	if remoteFileUrl, ok := ctx.Value("remoteFileUrl").(string); ok {
		info.RemoteFileUrl = remoteFileUrl
	} else {
		return DD_Addpiece_Param{}, IllegalExtInfoFormat
	}

	return info, nil
}

func WithValue(ctx context.Context, pieceCid cid.Cid, remoteFileUrl string) context.Context {
	ctx = context.WithValue(ctx, "pieceCID", pieceCid.String())
	ctx = context.WithValue(ctx, "remoteFileUrl", remoteFileUrl)
	return ctx
}
