package tdog

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func Test_Build(t *testing.T) {
	// 随机16位秘钥
	iv := NewUtil().RandomStr(16, 1, 2, 3)

	// init .
	jwtTdog := NewJwt()

	// memberId
	var trueMemberId int64 = 6683020668019466240

	// data
	testData := map[string]interface{}{
		"memberId": trueMemberId,
	}

	// generator .
	jwt := jwtTdog.Build(testData, iv)

	if len(jwt) < 1 {
		t.Fatal("Generator Error: Null result.")
	}

	if !jwtTdog.Valid(jwt, iv) {
		t.Fatal("Generator Error: Error result. JWT: " + jwt + " ; IV: " + iv)
	}

	// valid get data .
	dt := jwtTdog.GetData(jwt, iv)
	assert.Equal(t, int64(dt["memberId"].(float64)), trueMemberId, "get data is not same result.")

	// valid get by key .
	getMemberId := int64(jwtTdog.Get(jwt, "memberId", iv).(float64))
	assert.Equal(t, getMemberId, trueMemberId, "get memberId is not same result.")
}

func Test_Refresh(t *testing.T) {
	// 秘钥
	iv := "C3hr9iiK2qsYFpeH"

	// 数据
	jwt := "eyJUeXBlIjoiSldUIiwiQWxnb3JpdGhtIjoiSFMyNTYifQ==.6UJ+68nKGPm8kl3p070upJIPCA3TJhJjHXzDyOOVkJqRW4+fjVUc6Tg/B34P0SPaLtF19+9Ijz8PRQVz+SFC+w==.43915a05d590b38b570752f960ba91265454ee8edd01ccc8c9caf35a1ef9af63"

	// init .
	jwtTdog := NewJwt()

	// memberId
	var trueMemberId int64 = 6683020668019466240

	// Refresh .
	jwt = jwtTdog.Refresh(jwt, iv)

	if len(jwt) < 1 {
		t.Fatal("Generator Error: Null result.")
	}

	if !jwtTdog.Valid(jwt, iv) {
		t.Fatal("Generator Error: Error result. JWT: " + jwt + " ; IV: " + iv)
	}

	// valid get data .
	dt := jwtTdog.GetData(jwt, iv)
	assert.Equal(t, int64(dt["memberId"].(float64)), trueMemberId, "get data is not same result.")

	// valid get by key .
	getMemberId := int64(jwtTdog.Get(jwt, "memberId", iv).(float64))
	assert.Equal(t, getMemberId, trueMemberId, "get memberId is not same result.")
}
