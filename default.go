package mopan

import (
	"strings"

	"github.com/google/uuid"
)

const (
	HomeUrl       = "https://mopan.sc.189.cn"
	MoPanProxyUrl = HomeUrl + "/mopanproxy"
	// EnterPriseUrl = HomeUrl + "/enterprise"

	MoPanProxyUpdload = MoPanProxyUrl + "/fileupload"
	MoPanProxyFamily  = MoPanProxyUrl + "/family"
	MoPanProxyAuthUrl = MoPanProxyUrl + "/auth"
)

const DefaultPublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzDzuXrkkYjcb0N2/oiZX
xNNaYM0PKPZE0aXeCJWo/VBSFE3+Q/1m5u1hI7y6U6OP8Kj5w6BStHU4lqs3/OOW
G8GYKqGD6pz+b8Vp6m22lk8/7lL9375w2siAz+xSWwovAIKTfbMRwUsmJWoGI2vx
rwok6jJoWacP6GcsI335cD7fNsHSOFYTb7SCjKWvAowsHhAWu7W8oP7bB3HE3Xth
6Wy/gbZl/4Hp9rJU8w44/1Hc6O+uzfw4ZNtE0E4cIsK40XifW5SSokpCQkIlPNKH
RuzjIuGQRbCjvl682M/DixSouc4whOcOB6Rf102p2XaKvrmmT1OXCA4dFkRa5rjA
pQIDAQAB
-----END PUBLIC KEY-----
`

// const DefaultPublicKey2 = `
// -----BEGIN PUBLIC KEY-----
// MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAx2bxNtzCYin2B6qWRIwPNLP6E4arzRlpPrlb2UIjWq91pJe8yN2
// QwGigxFDTfkcrv32gqlTJX2xvMr1O7RzA+oIPfA1xfJTzfQHf2HPZ8w67A4WNZHxQmWqcDdUcy6JCKzks1TLGsAH5v17dK/
// AazM2u6n5OvFrqQMnXr/raZhJRVUg3YVXW6Ppbw7fewX2n1DosC+xLU19fpyHSb/YW/9dlDMJ4tvTHrxTxpT8OOM5/bdl5q
// eUN8bBsZht1l97Iyp1Od0oFDbBaorFUsyVEnVa7r5fuFlYSoLgLiCXnMNTLpJF4GbSvEG2vXAmTLrlJ+qYWXBL7O1AJU6tZ
// KchY4wIDAQAB
// -----END PUBLIC KEY-----
// `

// 默认设备信息
// 根据windwos信息构建
var DefaultDeviceInfo = DeviceInfo{
	DeviceNo: "1104a897925070c638d",

	MpRemoteType:    "3",
	MpRemoteChannel: "100",

	MpVersion:     "1.1.201",
	MpVersionCode: 145,

	MpDeviceSerialNum: strings.ReplaceAll(uuid.NewString(), "-", ""),
	MpManufcaturer:    "Windows端",
	MpModel:           "",

	MpOs:         "windows",
	MpOsVersion:  "31",
	MpOsVersion2: "12",
}
