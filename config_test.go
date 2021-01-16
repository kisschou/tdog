package tdog

import (
	"os"
	"testing"
)

var ConfigModule Config

// 单元测试
func TestConfigGet(t *testing.T) {
	type test struct {
		key  string
		want string
	}
	_ = os.Setenv("CONFIG_PATH", "/Users/kisschou/data/golang/src/all-service/sources/BasicService/config")

	tests := map[string]test{
		"测试默认文件中的key获取": {key: "app_port", want: "8003"},
		"测试带文件指向的key获取": {key: "error.ERROR_UNAUTHOZED", want: "Unauthorized!"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			gotData := ConfigModule.Get(tc.key).String()
			if gotData != tc.want {
				t.Errorf("excepted: %#v, got: %#v", tc.want, gotData)
			}
		})
	}
}

// 基准测试
func BenchmarkConfigGet(b *testing.B) {
	type test struct {
		key  string
		want string
	}
	_ = os.Setenv("CONFIG_PATH", "/Users/kisschou/data/golang/src/all-service/sources/BasicService/config")

	tests := map[string]test{
		"测试默认文件中的key获取": {key: "app_port", want: "8003"},
		"测试带文件指向的key获取": {key: "error.ERROR_UNAUTHOZED", want: "Unauthorized!"},
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range tests {
			_ = ConfigModule.Get(tc.key).String()
		}
	}
}
