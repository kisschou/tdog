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
	"github.com/wenzhenxi/gorsa"
	ParentCrc32 "hash/crc32"
	"net/url"
	"strconv"
)

type Crypt struct {
	Str string
}

func (h *Crypt) Md5() string {
	return fmt.Sprintf("%x", ParentMd5.Sum([]byte(h.Str)))
}

func (h *Crypt) Sha1() string {
	return fmt.Sprintf("%x", ParentSha1.Sum([]byte(h.Str)))
}

func (h *Crypt) Sha256() string {
	hash := ParentSha256.New()
	hash.Write([]byte(h.Str))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

func (h *Crypt) Sha512() string {
	hash := ParentSha512.New()
	hash.Write([]byte(h.Str))
	sum := hash.Sum(nil)
	return hex.EncodeToString(sum)
}

func (h *Crypt) Crc32() string {
	return strconv.Itoa(int(ParentCrc32.ChecksumIEEE([]byte(h.Str))))
}

func (h *Crypt) Base64Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(h.Str))
}

func (h *Crypt) Base64Decode() string {
	data, err := base64.StdEncoding.DecodeString(h.Str)
	if err != nil {
		return ""
	}
	return string(data)
}

func (h *Crypt) UrlBase64Encode() string {
	return base64.URLEncoding.EncodeToString([]byte(h.Str))
}

func (h *Crypt) UrlBase64Decode() string {
	data, err := base64.URLEncoding.DecodeString(h.Str)
	if err != nil {
		return ""
	}
	return string(data)
}

func (h *Crypt) Urlencode() string {
	return url.QueryEscape(h.Str)
}

func (h *Crypt) Urldecode() string {
	data, err := url.QueryUnescape(h.Str)
	if err != nil {
		return ""
	}
	return data
}

func (h *Crypt) BiuPwdNewBuilder(password string) (salt string, newPassword string) {
	salt = NewUtil().RandomStr(16)
	h.Str = password
	h.Str = h.Sha512() + NewConfig().Get("hex_key").ToString() + salt
	newPassword = h.Sha512()
	return
}

func (h *Crypt) BiuPwdBuilder(salt string, password string) (newPassword string) {
	h.Str = password
	h.Str = h.Sha512() + NewConfig().Get("hex_key").ToString() + salt
	newPassword = h.Sha512()
	return
}

// 生成公私钥
func (h *Crypt) GenerateRsaKey(bits int) (publicKey, privateKey string) {
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
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
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
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	// 生成公钥
	publicKey = string(pem.EncodeToMemory(&publicBlock))
	return
}

// 公钥加密
func (h *Crypt) RsaPubEncode(pubKey string) string {
	pubEncrypt, err := gorsa.PublicEncrypt(h.Str, pubKey)
	if err != nil {
		return ""
	}
	return pubEncrypt
}

// 私钥解密
func (h *Crypt) RsaPriDecode(priKey string) string {
	priDecrypt, err := gorsa.PriKeyDecrypt(h.Str, priKey)
	if err != nil {
		return ""
	}
	return priDecrypt
}

// 私钥加密
func (h *Crypt) RsaPriEncode(priKey string) string {
	priEncode, err := gorsa.PriKeyEncrypt(h.Str, priKey)
	if err != nil {
		return ""
	}
	return priEncode
}

// 公钥解密
func (h *Crypt) RsaPubDecode(pubKey string) string {
	pubDecode, err := gorsa.PublicDecrypt(h.Str, pubKey)
	if err != nil {
		return ""
	}
	return pubDecode
}
