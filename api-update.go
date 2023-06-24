package mopan

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type UpdloadFileParam struct {
	ParentFolderId string
	FileName       string
	FileSize       int64
	File           io.ReadSeeker
}

type InitMultiUploadData struct {
	FileDataExists Bool   `json:"fileDataExists"`
	UploadFileID   string `json:"uploadFileId"`
	UploadType     Int    `json:"uploadType"`
	UploadHost     string `json:"uploadHost"`

	PartSize int      `json:"-"`
	PartInfo []string `json:"-"`
}

// 初始化上传
func (c *MoClient) InitMultiUpload(ctx context.Context, file UpdloadFileParam, paramOption []ParamOption, option ...RestyOption) (*InitMultiUploadData, error) {
	partSize := 10485760
	count := int(math.Ceil(float64(file.FileSize) / float64(partSize)))

	// 优先计算所需信息
	fileMd5 := md5.New()
	silceMd5 := md5.New()
	silceMd5Hexs := make([]string, 0, count)
	silceMd5Base64s := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		silceMd5.Reset()
		if _, err := io.CopyN(io.MultiWriter(fileMd5, silceMd5), file.File, int64(partSize)); err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, err
		}
		md5Byte := silceMd5.Sum(nil)
		silceMd5Hexs = append(silceMd5Hexs, strings.ToUpper(hex.EncodeToString(md5Byte)))
		silceMd5Base64s = append(silceMd5Base64s, fmt.Sprint(i, "-", base64.StdEncoding.EncodeToString(md5Byte)))
	}
	if _, err := file.File.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	fileMd5Hex := strings.ToUpper(hex.EncodeToString(fileMd5.Sum(nil)))
	sliceMd5Hex := fileMd5Hex
	if file.FileSize > int64(partSize) {
		sliceMd5Hex = strings.ToUpper(Md5Hex(strings.Join(silceMd5Hexs, "\n")))
	}

	param := Json{
		"parentFolderId": file.ParentFolderId,
		"fileName":       file.FileName,
		"fileSize":       file.FileSize,
		"fileMd5":        fileMd5Hex,
		"sliceMd5":       sliceMd5Hex,
		"sliceSize":      partSize,

		"limitrate": "10240000", // 限制速度??
		"source":    1,
	}
	ApplyParamOption(param, paramOption...)

	var resp InitMultiUploadData
	_, err := c.Request(MoPanProxyUpdload+"/service/initMultiUpload", param, &resp,
		append(option, func(request *resty.Request) {
			request.SetContext(ctx)
		})...)
	if err != nil {
		return nil, err
	}
	resp.PartSize = partSize
	resp.PartInfo = silceMd5Base64s
	return &resp, nil
}

type GetUploadedPartsInfoData struct {
	UploadFileID     string `json:"uploadFileId"`
	UploadedPartList string `json:"uploadedPartList"`
}

// 查询分片上传情况
func (c *MoClient) GetUploadedPartsInfo(uploadFileId string, option ...RestyOption) (*GetUploadedPartsInfoData, error) {
	param := Json{
		"uploadFileId": uploadFileId,
	}

	var resp GetUploadedPartsInfoData
	_, err := c.Request(MoPanProxyUpdload+"/service/getUploadedPartsInfo", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

type GetMultiUploadData struct {
	HTTPMethod string `json:"httpMethod"`
	HTTPURL    string `json:"httpURL"`

	ContentType   string `json:"contentType"`
	Authorization string `json:"authorization"`
	Date          string `json:"date"`
	Limitrate     string `json:"limitrate"`

	PartMD5    string `json:"partMD5"`
	PartNumber int    `json:"partNumber"`

	UploadID   string `json:"uploadId"`
	ExpireTime Time3  `json:"expireTime"`
}

func (m *GetMultiUploadData) Headers() map[string]string {
	return map[string]string{
		"Content-Type":  m.ContentType,
		"Authorization": m.Authorization,
		"Date":          m.Date,
		"x-amz-limit":   "rate=" + m.Limitrate,
		"Content-Md5":   m.PartMD5,
	}
}

func (m *GetMultiUploadData) NewRequest(ctx context.Context, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, m.HTTPURL, body)
	if err != nil {
		return nil, err
	}
	for k, v := range m.Headers() {
		req.Header.Set(k, v)
	}
	return req, nil
}

// 获取分片上传信息
func (c *MoClient) GetAllMultiUploadUrls(uploadFileId string, partInfo []string, option ...RestyOption) ([]GetMultiUploadData, error) {
	param := Json{
		"uploadFileId": uploadFileId,
		"partInfo":     strings.Join(partInfo, ","),
	}

	var resp []GetMultiUploadData
	_, err := c.Request(MoPanProxyUpdload+"/service/getAllMultiUploadUrls", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type CommitMultiUploadData struct {
	CreateDate string `json:"createDate"`
	FileMd5    string `json:"fileMd5"`
	FileName   string `json:"fileName"`
	FileSize   string `json:"fileSize"`
	Rev        string `json:"rev"`
	UserFileID string `json:"userFileId"`
	UserID     string `json:"userId"`
}

// 提交上传文件
func (c *MoClient) CommitMultiUploadFile(uploadFileId string, paramOption []ParamOption, option ...RestyOption) (*CommitMultiUploadData, error) {
	param := Json{
		"uploadFileId": uploadFileId,
		"opertype":     3,

		"isLog": "其他",
		// "filmingTime": "2006-01-02 15:04:05",
	}
	ApplyParamOption(param, paramOption...)

	var resp CommitMultiUploadData
	_, err := c.Request(MoPanProxyUpdload+"/service/commitMultiUploadFile", param, &resp, option...)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
