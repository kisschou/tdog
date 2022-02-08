package tdog

import (
	"encoding/json"
	"strings"
	"time"

	tc "github.com/kisschou/TypeConverter"
)

type (
	// JwtHeader 用来表明签名的加密算法 token 类型等
	JwtHeader struct {
		Type      string // 类型
		Algorithm string // 加密算法
	}

	// JwtSignature 通过 header 生明的加密方法生成 签名
	JwtSignature string

	// Jwt jwt数据
	Jwt struct {
		header    string
		payload   string
		signature string
	}
)

func (header *JwtHeader) New() *JwtHeader {
	return &JwtHeader{
		Type:      "JWT",
		Algorithm: "HS256",
	}
}

// NewJwt 初始化一个Jwt.
func NewJwt() *Jwt {
	return new(Jwt)
}

// Build 使用 data 的数据, 和16位的 iv 作为秘钥生成jwt字符串.
func (jwt *Jwt) Build(data map[string]interface{}, iv string) string {
	// header
	jsonData, _ := json.Marshal(new(JwtHeader).New())
	jwt.header = NewCrypt(string(jsonData)).Base64Encode()

	// payload
	payload := make(map[string]interface{})
	payload["data"] = data
	payload["ita"] = time.Now().Unix() + 7200
	jsonData, _ = json.Marshal(payload)
	jwt.payload, _ = NewCrypt(string(jsonData)).AesEncrypt([]byte(iv))

	// signature
	jsonData, _ = json.Marshal(payload)
	jwt.signature = NewCrypt(string(jsonData) + "_cryptWithSalt-" + iv).Sha256()

	return jwt.header + "." + jwt.payload + "." + jwt.signature
}

func walk(input string) *Jwt {
	jwtData := strings.Split(input, ".")

	defer Recover()

	return &Jwt{
		header:    jwtData[0],
		payload:   jwtData[1],
		signature: jwtData[2],
	}
}

// valid 使用 iv 作为秘钥, 校验 input 字符串.
// 跳过过期时间校验.
func (jwt *Jwt) valid(input string, iv string) bool {
	jwt = walk(input)
	if jwt == nil {
		return false
	}

	// valid header .
	baseJwtHeader, _ := json.Marshal(new(JwtHeader).New())
	if jwt.header != NewCrypt(string(baseJwtHeader)).Base64Encode() {
		return false
	}

	// valid signature .
	dt, _ := NewCrypt(jwt.payload).AesDecrypt([]byte(iv))
	if NewCrypt(dt+"_cryptWithSalt-"+iv).Sha256() != jwt.signature {
		return false
	}

	return true
}

// Valid 使用 iv 作为秘钥, 校验 input 字符串.
func (jwt *Jwt) Valid(input string, iv string) bool {
	jwt = walk(input)
	if jwt == nil {
		return false
	}

	// valid header .
	baseJwtHeader, _ := json.Marshal(new(JwtHeader).New())
	if jwt.header != NewCrypt(string(baseJwtHeader)).Base64Encode() {
		return false
	}

	// valid exp .
	dt, _ := NewCrypt(jwt.payload).AesDecrypt([]byte(iv))
	dm := make(map[string]interface{}, 0)
	_ = json.Unmarshal([]byte(dt), &dm)
	exp := tc.New(dm["ita"]).Int64
	if time.Now().Unix() > exp {
		return false
	}

	// valid signature .
	if NewCrypt(dt+"_cryptWithSalt-"+iv).Sha256() != jwt.signature {
		return false
	}

	return true
}

// Get 使用 iv 作为秘钥, 解析 input 字符串, 并从其中的数据集合中获取下标为 key .
func (jwt *Jwt) Get(input string, key string, iv string) (value interface{}) {
	if jwt.valid(input, iv) {
		jwt = walk(input)

		dt, _ := NewCrypt(jwt.payload).AesDecrypt([]byte(iv))
		dm := make(map[string]interface{}, 0)
		_ = json.Unmarshal([]byte(dt), &dm)
		data := make(map[string]interface{}, 0)
		if dm["data"] != nil {
			data = dm["data"].(map[string]interface{})
		}

		if NewUtil().Isset(key, data) {
			value = data[key]
		} else {
			value = nil
		}
	}
	return
}

// GetData 使用 iv 作为秘钥, 解析 input 字符串, 取出其中的数据集合 .
func (jwt *Jwt) GetData(input string, iv string) (data map[string]interface{}) {
	if jwt.valid(input, iv) {
		jwt = walk(input)

		dt, _ := NewCrypt(jwt.payload).AesDecrypt([]byte(iv))
		dm := make(map[string]interface{}, 0)
		_ = json.Unmarshal([]byte(dt), &dm)
		if dm["data"] != nil {
			data = dm["data"].(map[string]interface{})
		}
	}
	return
}

// Refresh .
func (jwt *Jwt) Refresh(input string, iv string) string {
	// header
	jsonData, _ := json.Marshal(new(JwtHeader).New())
	jwt.header = NewCrypt(string(jsonData)).Base64Encode()

	// payload
	payload := make(map[string]interface{})
	payload["data"] = jwt.GetData(input, iv)
	payload["ita"] = time.Now().Unix() + 7200
	jsonData, _ = json.Marshal(payload)
	jwt.payload, _ = NewCrypt(string(jsonData)).AesEncrypt([]byte(iv))

	// signature
	jsonData, _ = json.Marshal(payload)
	jwt.signature = NewCrypt(string(jsonData) + "_cryptWithSalt-" + iv).Sha256()

	return jwt.header + "." + jwt.payload + "." + jwt.signature
}
