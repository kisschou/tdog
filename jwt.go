package tdog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	// JwtHeader 用来表明签名的加密算法 token 类型等
	JwtHeader struct {
		Type      string // 类型
		Algorithm string // 加密算法
	}

	// JwtPayload Payload 记录你需要的信息。 其中应该包含 Claims
	JwtPayload map[string]interface{}

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

// New USAGE:
// jwt := new(core.Jwt)
// data := make(map[string]interface{})
// data["username"] = username
// data["password"] = password
// jwt.New(data)
func (jwt *Jwt) New(data JwtPayload) string {

	// header
	jwtHeader := make(map[string]string)
	header := new(JwtHeader)
	header = header.New()
	jwtHeader["type"] = header.Type
	jwtHeader["alg"] = header.Algorithm
	jsonData, _ := json.Marshal(jwtHeader)
	jwt.header = NewCrypt(string(jsonData)).Base64Encode()

	// payload
	payload := make(map[string]interface{})
	payload["data"] = data
	payload["ita"] = time.Now().Unix()
	payload["exp"] = 7200
	jsonData, _ = json.Marshal(payload)
	jwt.payload = NewCrypt(string(jsonData)).Base64Encode()

	// signature
	jsonData, _ = json.Marshal(payload)
	jwt.signature = NewCrypt(string(jsonData) + NewConfig().Get("hex_key").ToString()).Sha256()

	return jwt.header + "." + jwt.payload + "." + jwt.signature
}

func (jwt *Jwt) Walk(data string) *Jwt {
	jwtData := strings.Split(data, ".")
	return &Jwt{
		header:    jwtData[0],
		payload:   jwtData[1],
		signature: jwtData[2],
	}
}

func (jwt *Jwt) Check(data string) bool {
	jwtData := strings.Split(data, ".")
	if len(jwtData) != 3 {
		return false
	}

	jwt = jwt.Walk(data)

	// check header.
	header := new(JwtHeader)
	header = header.New()
	jwtHeader := make(map[string]string)
	_ = json.Unmarshal([]byte(NewCrypt(jwt.header).Base64Decode()), &jwtHeader)
	if jwtHeader["type"] != header.Type || jwtHeader["alg"] != header.Algorithm {
		return false
	}

	// check payload.
	jwtPayload := make(map[string]interface{})
	_ = json.Unmarshal([]byte(NewCrypt(jwt.payload).Base64Decode()), &jwtPayload)
	ita, _ := strconv.Atoi(fmt.Sprintf("%1.0f", jwtPayload["ita"]))
	exp, _ := strconv.Atoi(fmt.Sprintf("%1.0f", jwtPayload["exp"]))
	ita = ita + exp
	if ita < int(time.Now().Unix()) {
		return false
	}

	// check signature.
	if jwt.signature != NewCrypt(NewCrypt(jwt.payload).Base64Decode()+NewConfig().Get("hex_key").ToString()).Sha256() {
		return false
	}

	return true
}

func (jwt *Jwt) Refresh(authorization string) string {
	if jwt.Check(authorization) {
		return authorization
	}
	jwt = jwt.Walk(authorization)
	// check payload.
	jwtPayload := make(map[string]interface{})
	_ = json.Unmarshal([]byte(NewCrypt(jwt.payload).Base64Decode()), &jwtPayload)
	return jwt.New(jwtPayload["data"].(map[string]interface{}))
}

func (jwt *Jwt) Get(data string, key string) (value interface{}) {
	if !jwt.Check(data) {
		return ""
	}
	jwt = jwt.Walk(data)
	jwtPayload := make(map[string]interface{})
	_ = json.Unmarshal([]byte(NewCrypt(jwt.payload).Base64Decode()), &jwtPayload)
	list := jwtPayload["data"].(map[string]interface{})
	if _, ok := list[key]; ok {
		value = list[key]
		return
	}
	return
}

func (jwt *Jwt) GetData(data string) map[string]interface{} {
	if !jwt.Check(data) {
		return nil
	}
	jwt = jwt.Walk(data)
	jwtPayload := make(map[string]interface{})
	_ = json.Unmarshal([]byte(NewCrypt(jwt.payload).Base64Decode()), &jwtPayload)
	return jwtPayload["data"].(map[string]interface{})
}
