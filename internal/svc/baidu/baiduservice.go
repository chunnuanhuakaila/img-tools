package baidu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"img-tools/internal/errors"
	"img-tools/internal/svc"
	"img-tools/internal/types"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var currIndexs = map[string]int{}
var accessInfo = map[int]string{}
var accessExpire = map[int]int64{}

type BaiduRsp struct {
	LogID     int    `json:"log_id"`
	Image     string `json:"image"`
	ErrorMsg  string `json:"error_msg"`
	ErrorCode int    `json:"error_code"`
}

// 随机选取一组ak/sk
func chooseFromBaidu(svcCtx *svc.ServiceContext, action string) (string, string, int) {
	list := svcCtx.Config.BaiduAI
	currIndexs[action]++
	if currIndexs[action] >= len(list) {
		currIndexs[action] = 0
	}
	return list[currIndexs[action]].AK, list[currIndexs[action]].SK, currIndexs[action]
}

func GetAccessToken(ctx context.Context, svc *svc.ServiceContext, action string) string {
	ak, sk, index := chooseFromBaidu(svc, action)
	if accessInfo[index] != "" && accessExpire[index] > time.Now().Unix() {
		return accessInfo[index]
	}
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", ak, sk)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		logx.WithContext(ctx).Errorf("get access token error:%v", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.WithContext(ctx).Errorf("get access token error:%v", err)
		return ""
	}
	accessTokenObj := map[string]any{}
	if err = json.Unmarshal([]byte(body), &accessTokenObj); err != nil {
		logx.WithContext(ctx).Errorf("json unmarshal access token error:%v", err)
		return ""
	}
	tenDaysLater := time.Now().AddDate(0, 0, 10)
	timestamp := tenDaysLater.Unix()
	accessExpire[index] = timestamp
	return accessTokenObj["access_token"].(string)
}

func EnhancePhoto(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	token := GetAccessToken(ctx, svcCtx, req.Action)
	if token == "" {
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	url := "https://aip.baidubce.com/rest/2.0/image-process/v1/image_definition_enhance?access_token=" + token
	payload := strings.NewReader("image=" + neturl.QueryEscape(req.BinaryData))
	baiduReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("new request error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	baiduReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	baiduReq.Header.Add("Accept", "application/json")

	return explainRsp(ctx, baiduReq)
}

func Recover(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	var rect []interface{}
	err := json.Unmarshal([]byte(req.ExtraInfo), &rect)
	if err != nil || len(rect) < 1 {
		return rsp, errors.APIError{Code: errors.ParamsInvalid, Msg: "矩形参数不合法"}
	}
	token := GetAccessToken(ctx, svcCtx, req.Action)
	if token == "" {
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	url := "https://aip.baidubce.com/rest/2.0/image-process/v1/inpainting?access_token=" + token
	payload, err := json.Marshal(map[string]interface{}{
		"rectangle": rect,
		"image":     req.BinaryData,
	})
	if err != nil {
		return rsp, errors.APIError{Code: errors.ParamsInvalid, Msg: "图片解析失败"}
	}
	baiduReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		logx.WithContext(ctx).Errorf("new request error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	baiduReq.Header.Add("Content-Type", "application/json")
	baiduReq.Header.Add("Accept", "application/json")

	return explainRsp(ctx, baiduReq)
}

func StretchRecovery(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	token := GetAccessToken(ctx, svcCtx, req.Action)
	if token == "" {
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	url := "https://aip.baidubce.com/rest/2.0/image-process/v1/stretch_restore?access_token=" + token
	payload := strings.NewReader("image=" + neturl.QueryEscape(req.BinaryData))
	baiduReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("new request error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	baiduReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	baiduReq.Header.Add("Accept", "application/json")

	return explainRsp(ctx, baiduReq)
}

func OverResolution(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	token := GetAccessToken(ctx, svcCtx, req.Action)
	if token == "" {
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	url := "https://aip.baidubce.com/rest/2.0/image-process/v1/image_quality_enhance?access_token=" + token
	payload := strings.NewReader("image=" + neturl.QueryEscape(req.BinaryData))
	baiduReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("new request error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	baiduReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	baiduReq.Header.Add("Accept", "application/json")

	return explainRsp(ctx, baiduReq)
}

func ConvertPhoto(ctx context.Context, svcCtx *svc.ServiceContext, req *types.ImgReq) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	token := GetAccessToken(ctx, svcCtx, req.Action)
	if token == "" {
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	url := "https://aip.baidubce.com/rest/2.0/image-process/v1/colourize?access_token=" + token
	payload := strings.NewReader("image=" + neturl.QueryEscape(req.BinaryData))

	baiduReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("new request error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	baiduReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	baiduReq.Header.Add("Accept", "application/json")

	return explainRsp(ctx, baiduReq)
}

func explainRsp(ctx context.Context, baiduReq *http.Request) (*types.ImgRsp, error) {
	rsp := &types.ImgRsp{}
	client := &http.Client{}
	res, err := client.Do(baiduReq)
	if err != nil {
		logx.WithContext(ctx).Errorf("request error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	if res.StatusCode != 200 {
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		logx.WithContext(ctx).Errorf("request error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	var ret BaiduRsp
	err = json.Unmarshal(body, &ret)
	if err != nil {
		logx.WithContext(ctx).Errorf("explain struct error:%v", err)
		return rsp, errors.APIError{Code: errors.Unavilable, Msg: errors.ErrMsgMap[errors.Unavilable]}
	}
	if ret.ErrorCode != 0 {
		return rsp, errors.APIError{Code: ret.ErrorCode, Msg: ret.ErrorMsg}
	}
	rsp.BinaryData = ret.Image
	return rsp, nil
}
