## 代码自动生成
goctl api go -api ./img.api -dir .

## 运行方式
go run img.go

## 接口说明
### 图像增强限制：
图片大小及尺寸：最大 5 MB; 输入图片尺寸 width 满足 [50, 2128], height 满足 [50, 4046]
```
curl --location '127.0.0.1:9999/img/edit' \
--header 'token: sss' \  
--header 'Content-Type: application/json' \
--data '{
    "binaryData":"", //文件的md5值
    "action":"enhance",  // 写死
    "timestamp":1111 // 时间戳
}'
```
返回：
```
{
    "code": 0,
    "msg": "",
    "binaryData": "",
    "url":"" //预留
}
```

### 图像超分限制：
图片大小及尺寸：最大 5 MB，宽必须在[50, 2160]范围内, 长必须在[50, 4096]范围内。
```
curl --location '127.0.0.1:9999/img/edit' \
--header 'token: sss' \  
--header 'Content-Type: application/json' \
--data '{
    "binaryData":"", //文件的md5值
    "action":"overResolution",  // 写死
    "timestamp":1111 // 时间戳

}'
```
返回：
```
{
    "code": 0,
    "msg": "",
    "binaryData": "",
    "url":"" //预留
}
```
### 老照片修复：
1. 图片文件大小：最大 5 MB。
2. 图片分辨率最大2000 x 2000像素。
```
curl --location '127.0.0.1:9999/img/edit' \
--header 'token: sss' \  
--header 'Content-Type: application/json' \
--data '{
    "binaryData":"", //文件的md5值
    "action":"convert",  // 写死
    "timestamp":1111 // 时间戳

}'
```
返回：
```
{
    "code": 0,
    "msg": "",
    "binaryData": "",
    "url":"" //预留
}
```   

### 图片拉伸
图片文件大小：最大 5 MB。
```
curl --location '127.0.0.1:9999/img/edit' \
--header 'token: sss' \  
--header 'Content-Type: application/json' \
--data '{
    "binaryData":"", //文件的md5值
    "action":"stretchRecovery",  // 写死
    "timestamp":1111 // 时间戳

}'
```
返回：
```
{
    "code": 0,
    "msg": "",
    "binaryData": "",
    "url":"" //预留
}
```   

### 图片风格转换
图片文件大小：最大 5 MB。
```
curl --location '127.0.0.1:9999/img/edit' \
--header 'token: sss' \  
--header 'Content-Type: application/json' \
--data '{
    "binaryData":"", //文件的md5值
    "action":"styleConversion",  // 写死
    "timestamp":1111, // 时间戳
    "extraInfo":"watercolor_cartoon" //watercolor_cartoon(水彩风) | jzcartoon(剪纸风)
}'
```
返回：
```
{
    "code": 0,
    "msg": "",
    "binaryData": "",
    "url":"" //预留
}
```

### 图片修复
```
curl --location '127.0.0.1:9999/img/edit' \
--header 'token: sss' \
--header 'Content-Type: application/json' \
--data '{
    "binaryData":"",
    "action":"recover",
    "timestamp":1111,
    "extraInfo":"[{\"width\":1,\"height\":2,\"top\":3,\"left\":4}]"
}'
```
返回：
```
{
    "code": 0,
    "msg": "",
    "binaryData": "",
    "url":"" //预留
}
```

