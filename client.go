package mopan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
)

type RestyOption func(request *resty.Request)
type Json map[string]any

func NewMoClientWithAuthorization(authorization string) *MoClient {
	return NewMoClient().SetAuthorization(authorization)
}

func NewMoClient() *MoClient {
	return &MoClient{
		Client:     resty.New(),
		DeviceInfo: DefaultDeviceInfo.Encrypt(),
	}
}

type MoClient struct {
	Authorization string
	DeviceInfo    string

	Client *resty.Client
}

func (c *MoClient) SetDeviceInfo(info string) *MoClient {
	if info != "" {
		c.DeviceInfo = info
	}
	return c
}

func (c *MoClient) SetAuthorization(authorization string) *MoClient {
	if !strings.HasPrefix("Bearer", authorization) {
		authorization = "Bearer " + authorization
	}
	c.Authorization = authorization
	return c
}

func (c *MoClient) SetClient(client *http.Client) *MoClient {
	c.Client = resty.NewWithClient(client)
	return c
}

func (c *MoClient) SetRestyClient(client *resty.Client) *MoClient {
	c.Client = client
	return c
}

func (c *MoClient) SetProxy(proxy string) *MoClient {
	c.Client.SetProxy(proxy)
	return c
}

func (c *MoClient) request(url string, data Json, resp any, option ...RestyOption) ([]byte, error) {
	secretKey := GetSecretKey()
	encryptedKey := MustRsaEncryptBase64Str(secretKey, DefaultPublicKey)
	req := c.Client.R().SetHeaders(map[string]string{
		"Authorization": c.Authorization,
		"encrypted-key": encryptedKey,
		"remoteInfo":    c.DeviceInfo,
	})

	if data != nil {
		req.SetHeader("Content-Type", "application/json")
		temp, err := c.Client.JSONMarshal(data)
		if err != nil {
			return nil, err
		}
		enc, _ := AesEncryptBase64(temp, []byte(secretKey))
		req.SetBody(enc)
	}

	for _, opt := range option {
		opt(req)
	}
	resp_, err := req.Post(url)
	if err != nil {
		return nil, err
	}

	body := resp_.Body()
	// 解密数据
	if bytes.HasPrefix(body, []byte{'"'}) && bytes.HasSuffix(body, []byte{'"'}) {
		body, err = AesDecryptBase64(bytes.Trim(body, "\""), []byte(secretKey))
		if err != nil {
			return nil, err
		}
	}

	var result Resp
	c.Client.JSONUnmarshal(body, &result)
	if result.Code != 200 {
		return nil, &result
	}

	if resp != nil {
		if err := c.Client.JSONUnmarshal(result.Data, &resp); err != nil {
			return nil, err
		}
	}
	return result.Data, nil
}

func (c *MoClient) Request(url string, data Json, resp any, option ...RestyOption) ([]byte, error) {
	return c.request(url, data, resp, option...)
}

type Resp struct {
	Code    int64           `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Status  bool            `json:"status"`
}

func (r *Resp) Error() string {
	return fmt.Sprintf("Code:%d, Message:%s", r.Code, r.Message)
}
