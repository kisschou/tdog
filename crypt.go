// Copyright 2012 Kisschou. All rights reserved.
// Based on the path package, Copyright 2011 The Go Authors.
// Use of this source code is governed by a MIT-style license that can be found
// at https://github.com/kisschou/tdog/blob/master/LICENSE.

package tdog

import (
	ParentMd5 "crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	ParentSha1 "crypto/sha1"
	ParentSha256 "crypto/sha256"
	ParentSha512 "crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	ParentCrc32 "hash/crc32"
	"net/url"
	"strconv"

	"github.com/wenzhenxi/gorsa"
)

/**
 * The module for crypt handler.
 *
 * @Author: Kisschou
 * @Build: 2021-05-06
 */
type crypt struct {
	input string
}

// NewCrypt init a new crypt module
func NewCrypt(input string) *crypt {
	return &crypt{input: input}
}

// Md5 md5加密
func (h *crypt) Md5() string {
	return fmt.Sprintf("%x", ParentMd5.Sum([]byte(h.input)))
}

// Sha1 sha1加密
func (h *crypt) Sha1() string {
	return fmt.Sprintf("%x", ParentSha1.Sum([]byte(h.input)))
}

// Sha256 sha256加密
func (h *crypt) Sha256() string {
	hash := ParentSha256.New()
	hash.Write([]byte(h.input))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

// Sha512 sha512加密
func (h *crypt) Sha512() string {
	hash := ParentSha512.New()
	hash.Write([]byte(h.input))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

// Crc32 循环冗余校验
func (h *crypt) Crc32() string {
	return strconv.Itoa(int(ParentCrc32.ChecksumIEEE([]byte(h.input))))
}

// Base64Encode base64加密
func (h *crypt) Base64Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(h.input))
}

// Base64Decode base64解密
func (h *crypt) Base64Decode() string {
	data, err := base64.StdEncoding.DecodeString(h.input)
	if err != nil {
		return ""
	}
	return string(data)
}

// UrlBase64Encode base64链接加密
func (h *crypt) UrlBase64Encode() string {
	return base64.URLEncoding.EncodeToString([]byte(h.input))
}

// UrlBase64Decode base64链接解密
func (h *crypt) UrlBase64Decode() string {
	data, err := base64.URLEncoding.DecodeString(h.input)
	if err != nil {
		return ""
	}
	return string(data)
}

// Urlencode url编码
func (h *crypt) Urlencode() string {
	return url.QueryEscape(h.input)
}

// Urldecode url解码
func (h *crypt) Urldecode() string {
	data, err := url.QueryUnescape(h.input)
	if err != nil {
		return ""
	}
	return data
}

func (h *crypt) BiuPwdNewBuilder(password string) (salt string, newPassword string) {
	salt = NewUtil().RandomStr(16)
	h.input = password
	h.input = h.Sha512() + NewConfig().Get("hex_key").ToString() + salt
	newPassword = h.Sha512()
	return
}

func (h *crypt) BiuPwdBuilder(salt string, password string) (newPassword string) {
	h.input = password
	h.input = h.Sha512() + NewConfig().Get("hex_key").ToString() + salt
	newPassword = h.Sha512()
	return
}

// GenerateRsaKey 生成公私钥
// @param bits int RSA密钥的字位数
// @return publicKey string 公钥
// @return privateKey string 私钥
func (h *crypt) GenerateRsaKey(bits int) (publicKey, privateKey string) {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		publicKey, privateKey = "", ""
		return
	}
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(priKey)
	//使用pem格式对x509输出的内容进行编码
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA PRIVATE KEY", Bytes: X509PrivateKey}
	// 生成私钥
	privateKey = string(pem.EncodeToMemory(&privateBlock))
	//获取公钥的数据
	pubKey := priKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&pubKey)
	if err != nil {
		publicKey, privateKey = "", ""
		return
	}
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "PUBLIC KEY", Bytes: X509PublicKey}
	// 生成公钥
	publicKey = string(pem.EncodeToMemory(&publicBlock))
	return
}

// RsaPubEncode 公钥加密
// @param publicKey string 公钥
// @return string 加密结果
func (h *crypt) RsaPubEncode(pubKey string) string {
	pubEncrypt, err := gorsa.PublicEncrypt(h.input, pubKey)
	if err != nil {
		return ""
	}
	return pubEncrypt
}

// RsaPriDecode 私钥解密
// @param priKey string 公钥
// @return string 解密结果
func (h *crypt) RsaPriDecode(priKey string) string {
	priDecrypt, err := gorsa.PriKeyDecrypt(h.input, priKey)
	if err != nil {
		return ""
	}
	return priDecrypt
}

// RsaPriEncode 私钥加密
// @param priKey string 公钥
// @return string 加密结果
func (h *crypt) RsaPriEncode(priKey string) string {
	priEncode, err := gorsa.PriKeyEncrypt(h.input, priKey)
	if err != nil {
		return ""
	}
	return priEncode
}

// RsaPubDecode 公钥解密
// @param publicKey string 公钥
// @return string 解密结果
func (h *crypt) RsaPubDecode(pubKey string) string {
	pubDecode, err := gorsa.PublicDecrypt(h.input, pubKey)
	if err != nil {
		return ""
	}
	return pubDecode
}
