package img

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"img-tools/internal/errors"
	"img-tools/internal/svc"
	"img-tools/internal/svc/baidu"
	"img-tools/internal/svc/zijie"
	"img-tools/internal/types"
	"img-tools/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEditLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditLogic {
	return &EditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EditLogic) Edit(req *types.ImgReq) (resp *types.ImgRsp, err error) {
	resp = &types.ImgRsp{}
	if (req.URL == "" && req.BinaryData == "") || req.Token == "" {
		return resp, errors.APIError{Code: errors.ParamsInvalid, Msg: errors.ErrMsgMap[errors.ParamsInvalid]}
	}
	// 验证token
	if err = l.checkToken(req); err != nil {
		return resp, err
	}
	// 验证action
	action, ok := zijie.ImageActionMap[req.Action]
	if !ok {
		return resp, errors.APIError{Code: errors.ParamsInvalid, Msg: "暂不支持该编辑类型"}
	}
	if req.BinaryData != "" {
		req.BinaryData = strings.SplitN(req.BinaryData, ",", 2)[1]
	}
	// 验证图片大小
	if req.BinaryData != "" {
		width, height, size, err := utils.GetImgInfoFromByte(req.BinaryData)
		if err == nil {
			if size > 5*1024*1024 {
				return resp, errors.APIError{Code: errors.ParamsInvalid, Msg: "图片不能超过5M"}
			}
			if limit, ok := zijie.ImageWidthLimit[action]; ok {
				tmp := strings.Split(limit, ",")
				if width < utils.InterfaceToInt(tmp[0]) || width > utils.InterfaceToInt(tmp[1]) {
					return resp, errors.APIError{Code: errors.ParamsInvalid, Msg: "图片尺寸超出限制"}
				}
			}
			if limit, ok := zijie.ImageHeightLimit[action]; ok {
				tmp := strings.Split(limit, ",")
				if height < utils.InterfaceToInt(tmp[0]) || height > utils.InterfaceToInt(tmp[1]) {
					return resp, errors.APIError{Code: errors.ParamsInvalid, Msg: "图片尺寸超出限制"}
				}
			}
		}
	}

	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	t := random.Intn(2)
	t = 1
	switch action {
	case zijie.ConvertAction:
		if t == 0 {
			return zijie.ConvertPhoto(l.ctx, l.svcCtx, req)
		} else {
			return baidu.ConvertPhoto(l.ctx, l.svcCtx, req)
		}
	case zijie.EnhanceAction:
		if t == 0 {
			return zijie.EnhancePhoto(l.ctx, l.svcCtx, req)
		} else {
			return baidu.EnhancePhoto(l.ctx, l.svcCtx, req)
		}
	case zijie.OverResolutionAction:
		if t == 0 {
			return zijie.OverResolution(l.ctx, l.svcCtx, req)
		} else {
			return baidu.OverResolution(l.ctx, l.svcCtx, req)
		}
	case zijie.StretchRecoveryAction:
		if t == 0 {
			return zijie.StretchRecovery(l.ctx, l.svcCtx, req)
		} else {
			return baidu.StretchRecovery(l.ctx, l.svcCtx, req)
		}
	case zijie.StyleConversionAction:
		return zijie.StyleConversion(l.ctx, l.svcCtx, req)
	case zijie.RecoverAction:
		return baidu.Recover(l.ctx, l.svcCtx, req)
	}

	return
}

func (l *EditLogic) checkToken(req *types.ImgReq) error {
	// 超过5分钟认为过期
	if (time.Now().Unix() - req.TimeStamp) > 5*60 {
		l.Logger.Errorf("timestamp is outof 5 minites.")
		return errors.APIError{Code: errors.TokenInvalid, Msg: errors.ErrMsgMap[errors.TokenInvalid]}
	}
	var str string
	if req.BinaryData != "" {
		str = fmt.Sprintf("%s,%s,%d", req.BinaryData[:1000], req.Action, req.TimeStamp)
	} else {
		str = fmt.Sprintf("%s,%s,%d", req.URL, req.Action, req.TimeStamp)
	}
	hasher := md5.New()
	io.WriteString(hasher, str)
	sum := hasher.Sum(nil)
	md5String := hex.EncodeToString(sum)
	if md5String != req.Token {
		l.Logger.Errorf("token is uninvalid.")
		return errors.APIError{Code: errors.TokenInvalid, Msg: errors.ErrMsgMap[errors.TokenInvalid]}
	}
	return nil
}
