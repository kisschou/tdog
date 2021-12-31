package tdog

import (
	"fmt"
	"os"
	"testing"
)

// 常规 -->

func TestConfigGet(t *testing.T) {
	configTdog := NewConfig()

	// 指定路径、文件、前缀、批量获取
	trueList := map[string]string{
		"host":    "127.0.0.1",
		"port":    "3306",
		"user":    "root",
		"pass":    "root",
		"db":      "member_db",
		"charset": "utf8mb4",
		"prefix":  "",
	}
	rlist := configTdog.SetPath("./tests/config").SetFile("database").SetPrefix("member_service").GetMulti("host", "port", "user", "pass", "db", "charset", "prefix")
	for k, resultImpl := range rlist {
		if trueList[k] != resultImpl.ToString() {
			t.Errorf("fail. want: %v, got: %v", trueList[k], resultImpl.ToString())
		}
	}

	// 不指定路径、文件
	trueList = map[string]string{
		"host":    "127.0.0.1",
		"port":    "3306",
		"user":    "root",
		"pass":    "root",
		"db":      "oss_db",
		"charset": "utf8mb4",
		"prefix":  "",
	}
	rlist = configTdog.SetPrefix("oss_service").GetMulti("host", "port", "user", "pass", "db", "charset", "prefix")
	for k, resultImpl := range rlist {
		if trueList[k] != resultImpl.ToString() {
			t.Errorf("fail. want: %v, got: %v", trueList[k], resultImpl.ToString())
		}
	}

	// 单个获取
	os.Setenv("CONFIG_PATH", "./tests/config")
	configTdog = NewConfig()
	var tests = []struct {
		input string
		want  string
	}{
		{"app_name", "Demo"},
		{"hex_key", "1UgZfL<=Au3M9dQrcB7yzzd==8?xlZ:T3oP1id>jKWi4Bc9Kz<RbTSXqgPWJXxes945rzQ4ojlsqN25:95yb1-DsbAu1k6mUSs27CXDmQu-mx0lhH82-rA5=<NO0Z4:x"},
		{"cache.member_service.host", "127.0.0.1"},
		{"cache.member_service.port", "6379"},
		{"cache.member_service.pass", "WoBuShiRen"},
	}
	for _, test := range tests {
		if got := configTdog.Get(test.input).ToString(); got != test.want {
			t.Errorf("fail. want: %v, got: %v", test.want, got)
		}
	}
}

// <--

// 基准 -->

// BenchmarkConfigGet .
// => go test -bench=ConfigGet -run=none -benchmem
func BenchmarkConfigGet(b *testing.B) {
	configTdog := NewConfig()

	for i := 0; i < b.N; i++ {
		/*
			// 指定路径、文件、前缀、批量获取
			trueList := map[string]string{
				"host":    "127.0.0.1",
				"port":    "3306",
				"user":    "root",
				"pass":    "root",
				"db":      "member_db",
				"charset": "utf8mb4",
				"prefix":  "",
			}
			rlist := configTdog.SetPath("./tests/config").SetFile("database").SetPrefix("member_service").GetMulti("host", "port", "user", "pass", "db", "charset", "prefix")
			for k, resultImpl := range rlist {
				if trueList[k] != resultImpl.ToString() {
					fmt.Println("fail. want: ", trueList[k], ", got: ", resultImpl.ToString())
				}
			}

			// 不指定路径、文件
			trueList = map[string]string{
				"host":    "127.0.0.1",
				"port":    "3306",
				"user":    "root",
				"pass":    "root",
				"db":      "oss_db",
				"charset": "utf8mb4",
				"prefix":  "",
			}
			rlist = configTdog.SetPrefix("oss_service").GetMulti("host", "port", "user", "pass", "db", "charset", "prefix")
			for k, resultImpl := range rlist {
				if trueList[k] != resultImpl.ToString() {
					fmt.Println("fail. want: ", trueList[k], ", got: ", resultImpl.ToString())
				}
			}
		*/

		// 单个获取
		os.Setenv("CONFIG_PATH", "./tests/config")
		configTdog = NewConfig()
		var tests = []struct {
			input string
			want  string
		}{
			{"app_name", "Demo"},
			{"hex_key", "1UgZfL<=Au3M9dQrcB7yzzd==8?xlZ:T3oP1id>jKWi4Bc9Kz<RbTSXqgPWJXxes945rzQ4ojlsqN25:95yb1-DsbAu1k6mUSs27CXDmQu-mx0lhH82-rA5=<NO0Z4:x"},
			{"cache.member_service.host", "127.0.0.1"},
			{"cache.member_service.port", "6379"},
			{"cache.member_service.pass", "WoBuShiRen"},
		}
		for _, test := range tests {
			if got := configTdog.Get(test.input).ToString(); got != test.want {
				fmt.Println("fail. want: ", test.want, ", got: ", got)
			}
		}
	}
}

// <--
