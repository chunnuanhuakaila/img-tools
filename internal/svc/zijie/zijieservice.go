package zijie

import (
	"context"
	"img-tools/internal/errors"
	"img-tools/internal/svc"
	"img-tools/internal/types"
	"net/url"

	"github.com/volcengine/volc-sdk-golang/service/visual"
	"github.com/volcengine/volc-sdk-golang/service/visual/model"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	ConvertAction         = "ConvertPhotoV2"
	EnhanceAction         = "EnhancePhotoV2"
	OverResolutionAction  = "OverResolutionV2"
	StretchRecoveryAction = "StretchRecovery"
	StyleConversionAction = "ImageStyleConversion"
	RecoverAction         = "RecoverAction"
)

var ImageActionMap = map[string]string{
	"enhance":         EnhanceAction,
	"convert":         ConvertAction,
	"overResolution":  OverResolutionAction,
	"stretchRecovery": StretchRecoveryAction,
	"styleConversion": StyleConversionAction,
	"recover":         RecoverAction,
}

var ImageWidthLimit = map[string]string{
	EnhanceAction:        "50,2128",
	ConvertAction:        "50,2160",
	OverResolutionAction: "50,2000",
	//RecoverAction:        "50,2000",
}

var ImageHeightLimit = map[string]string{
	EnhanceAction:        "50,4046",
	ConvertAction:        "50,4046",
	OverResolutionAction: "50,2000",
	//RecoverAction:        "50,2000",
}

var currIndexs = map[string]int{}

// 随机选取一组ak/sk
func chooseFromZijie(svcCtx *svc.ServiceContext, action string) (string, string) {
	list := svcCtx.Config.ZijieAI
	currIndexs[action]++
	if currIndexs[action] >= len(list) {
		currIndexs[action] = 0
	}
	return list[currIndexs[action]].AK, list[currIndexs[action]].SK
}

func ConvertPhoto(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	ak, sk := chooseFromZijie(svcCtx, ConvertAction)
	visual.DefaultInstance.Client.SetAccessKey(ak)
	visual.DefaultInstance.Client.SetSecretKey(sk)

	form := &model.ConvertPhotoV2Request{
		ReqKey:           "lens_opr",
		BinaryDataBase64: []string{req.BinaryData},
		ImageUrls:        []string{req.URL},
		IfColor:          2,
	}
	resp, _, err := visual.DefaultInstance.ConvertPhotoV2(form)
	if err != nil {
		logx.WithContext(ctx).Errorf("req failed:%v", err)
		return rsp, err
	}
	if resp.Code != 10000 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	if len(resp.Data.BinaryDataBase64) < 1 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	rsp.BinaryData = resp.Data.BinaryDataBase64[0]
	return rsp, nil
}

func EnhancePhoto(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	ak, sk := chooseFromZijie(svcCtx, ConvertAction)
	visual.DefaultInstance.Client.SetAccessKey(ak)
	visual.DefaultInstance.Client.SetSecretKey(sk)

	form := &model.EnhancePhotoV2Request{
		ReqKey:             "lens_lqir",
		BinaryDataBase64:   []string{req.BinaryData},
		ResolutionBoundary: req.ExtraInfo,
	}
	resp, _, err := visual.DefaultInstance.EnhancePhotoV2(form)
	if err != nil {
		logx.WithContext(ctx).Errorf("req failed:%v", err)
		return rsp, err
	}
	if resp.Code != 10000 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	if len(resp.Data.BinaryDataBase64) < 1 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	rsp.BinaryData = resp.Data.BinaryDataBase64[0]
	return rsp, nil
}

func OverResolution(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	ak, sk := chooseFromZijie(svcCtx, ConvertAction)
	visual.DefaultInstance.Client.SetAccessKey(ak)
	visual.DefaultInstance.Client.SetSecretKey(sk)

	imageBase64List := []string{
		req.BinaryData,
	}
	resp, _, err := visual.DefaultInstance.OverResolutionV2(imageBase64List)
	if err != nil {
		logx.WithContext(ctx).Errorf("req failed:%v", err)
		return rsp, err
	}
	if resp.Code != 10000 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	if len(resp.Data.BinaryDataBase64) < 1 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	rsp.BinaryData = resp.Data.BinaryDataBase64[0]
	return rsp, nil
}

func StretchRecovery(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	ak, sk := chooseFromZijie(svcCtx, ConvertAction)
	visual.DefaultInstance.Client.SetAccessKey(ak)
	visual.DefaultInstance.Client.SetSecretKey(sk)

	form := url.Values{}
	form.Add("image_base64", req.BinaryData)
	resp, _, err := visual.DefaultInstance.StretchRecovery(form)
	if err != nil {
		logx.WithContext(ctx).Errorf("req failed:%v", err)
		return rsp, err
	}
	if resp.Code != 10000 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	rsp.BinaryData = resp.Data.Image
	return rsp, nil
}

func StyleConversion(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	if req.ExtraInfo != "jzcartoon" && req.ExtraInfo != "watercolor_cartoon" {
		return rsp, errors.APIError{Code: errors.ParamsInvalid, Msg: "风格暂不支持"}
	}
	ak, sk := chooseFromZijie(svcCtx, ConvertAction)
	visual.DefaultInstance.Client.SetAccessKey(ak)
	visual.DefaultInstance.Client.SetSecretKey(sk)
	form := url.Values{}
	form.Add("image_base64", req.BinaryData)
	form.Add("type", req.ExtraInfo)
	resp, _, err := visual.DefaultInstance.ImageStyleConversion(form)
	if err != nil {
		logx.WithContext(ctx).Errorf("req failed:%v", err)
		return rsp, err
	}
	if resp.Code != 10000 {
		return rsp, errors.APIError{Code: resp.Code, Msg: resp.Message}
	}
	rsp.BinaryData = resp.Data.Image
	return rsp, nil
}
