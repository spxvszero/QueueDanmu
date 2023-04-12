package bili

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

const AcceptHeader = "Accept"
const ContentTypeHeader = "Content-Type"
const AuthorizationHeader = "Authorization"
const JsonType = "application/json"
const BiliVersion = "1.0"
const HmacSha256 = "HMAC-SHA256"
const BiliTimestampHeader = "x-bili-timestamp"
const BiliSignatureMethodHeader = "x-bili-signature-method"
const BiliSignatureNonceHeader = "x-bili-signature-nonce"
const BiliAccessKeyIdHeader = "x-bili-accesskeyid"
const BiliSignVersionHeader = "x-bili-signature-version"
const BiliContentMD5Header = "x-bili-content-md5"

//CreateSignature 生成Authorization加密串
func CreateSignature(header *CommonHeader, accessKeySecret string) string {
	sStr := header.ToSortedString()
	return HmacSHA256(accessKeySecret, sStr)
}

func Md5(str string) (md5str string) {
	data := []byte(str)
	has := md5.Sum(data)
	md5str = fmt.Sprintf("%x", has)
	return md5str
}

func HmacSHA256(key string, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

type CommonHeader struct {
	ContentType       string
	ContentAcceptType string
	Timestamp         string
	SignatureMethod   string
	SignatureVersion  string
	Authorization     string
	Nonce             string
	AccessKeyId       string
	ContentMD5        string
}

// ToMap 所有字段转map<string, string>
func (h *CommonHeader) ToMap() map[string]string {
	return map[string]string{
		BiliTimestampHeader:       h.Timestamp,
		BiliSignatureMethodHeader: h.SignatureMethod,
		BiliSignatureNonceHeader:  h.Nonce,
		BiliAccessKeyIdHeader:     h.AccessKeyId,
		BiliSignVersionHeader:     h.SignatureVersion,
		BiliContentMD5Header:      h.ContentMD5,
		AuthorizationHeader:       h.Authorization,
		ContentTypeHeader:         h.ContentType,
		AcceptHeader:              h.ContentAcceptType,
	}
}

// ToSortMap 参与加密的字段转map<string, string>
func (h *CommonHeader) ToSortMap() map[string]string {
	return map[string]string{
		BiliTimestampHeader:       h.Timestamp,
		BiliSignatureMethodHeader: h.SignatureMethod,
		BiliSignatureNonceHeader:  h.Nonce,
		BiliAccessKeyIdHeader:     h.AccessKeyId,
		BiliSignVersionHeader:     h.SignatureVersion,
		BiliContentMD5Header:      h.ContentMD5,
	}
}

//ToSortedString 生成需要加密的文本
func (h *CommonHeader) ToSortedString() (sign string) {
	hMap := h.ToSortMap()
	var hSil []string
	for k := range hMap {
		hSil = append(hSil, k)
	}
	sort.Strings(hSil)
	for _, v := range hSil {
		sign += v + ":" + hMap[v] + "\n"
	}
	sign = strings.TrimRight(sign, "\n")
	return
}
