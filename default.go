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

// 默认设备信息
// 根据windwos信息构建
var DefaultDeviceInfo = DeviceInfo{
	DeviceNo: "1104a897925070c638d",

	MpRemoteType:    "3",
	MpRemoteChannel: "100",

	MpVersion:     "1.0.3008",
	MpVersionCode: 145,

	MpDeviceSerialNum: strings.ReplaceAll(uuid.NewString(), "-", ""),
	MpManufcaturer:    "Windows端",
	MpModel:           "",

	MpOs:         "Windows",
	MpOsVersion:  "31",
	MpOsVersion2: "12",
}
