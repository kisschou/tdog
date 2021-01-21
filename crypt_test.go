package tdog

import "testing"

func TestCrypt_Md5(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试md5方法",
			fields: fields{
				"123456",
			},
			want: "e10adc3949ba59abbe56e057f20f883e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Md5(); got != tt.want {
				t.Errorf("Crypt.Md5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Sha1(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试sha1方法",
			fields: fields{
				"123456",
			},
			want: "7c4a8d09ca3762af61e59520943dc26494f8941b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Sha1(); got != tt.want {
				t.Errorf("Crypt.Sha1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Sha256(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试sha256方法",
			fields: fields{
				"123456",
			},
			want: "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Sha256(); got != tt.want {
				t.Errorf("Crypt.Sha256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Sha512(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试sha512方法",
			fields: fields{
				"123456",
			},
			want: "ba3253876aed6bc22d4a6ff53d8406c6ad864195ed144ab5c87621b6c233b548baeae6956df346ec8c17f5ea10f35ee3cbc514797ed7ddd3145464e2a0bab413",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Sha512(); got != tt.want {
				t.Errorf("Crypt.Sha512() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Crc32(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试crc32方法",
			fields: fields{
				"123456",
			},
			want: "158520161",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Crc32(); got != tt.want {
				t.Errorf("Crypt.Crc32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Base64Encode(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试Base64Encode方法",
			fields: fields{
				"123456",
			},
			want: "MTIzNDU2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Base64Encode(); got != tt.want {
				t.Errorf("Crypt.Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Base64Decode(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试Base64Decode方法",
			fields: fields{
				"MTIzNDU2",
			},
			want: "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Base64Decode(); got != tt.want {
				t.Errorf("Crypt.Base64Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_UrlBase64Encode(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试UrlBase64Encode方法",
			fields: fields{
				"123456",
			},
			want: "MTIzNDU2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.UrlBase64Encode(); got != tt.want {
				t.Errorf("Crypt.UrlBase64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_UrlBase64Decode(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试UrlBase64Decode方法",
			fields: fields{
				"MTIzNDU2",
			},
			want: "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.UrlBase64Decode(); got != tt.want {
				t.Errorf("Crypt.UrlBase64Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Urlencode(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试Urlencode方法",
			fields: fields{
				"http://www.kisschou.com/?a=de1",
			},
			want: "http%3A%2F%2Fwww.kisschou.com%2F%3Fa%3Dde1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Urlencode(); got != tt.want {
				t.Errorf("Crypt.Urlencode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_Urldecode(t *testing.T) {
	type fields struct {
		Str string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "测试Urldecode方法",
			fields: fields{
				"http%3A%2F%2Fwww.kisschou.com%2F%3Fa%3Dde1",
			},
			want: "http://www.kisschou.com/?a=de1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if got := h.Urldecode(); got != tt.want {
				t.Errorf("Crypt.Urldecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypt_BiuPwdBuilder(t *testing.T) {
	type fields struct {
		Str string
	}
	type args struct {
		salt     string
		password string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantNewPassword string
	}{
		// TODO: Add test cases.
		{
			name:   "测试BiuPwdBuilder方法",
			fields: fields{},
			args: args{
				salt:     "66BHE7IUvXoomLQ3",
				password: "111111",
			},
			wantNewPassword: "b7f6e7f529674d9f2645fa5ec146122696d6fd9d74a794561bbf237fa7db1f5f11d7fbf22fb94ed6954e5fbdf756d2f23ea434574d7f22439d0dd94e0f72c901",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Crypt{
				Str: tt.fields.Str,
			}
			if gotNewPassword := h.BiuPwdBuilder(tt.args.salt, tt.args.password); gotNewPassword != tt.wantNewPassword {
				t.Errorf("Crypt.BiuPwdBuilder(%v, %v) = %v, want %v", tt.args.salt, tt.args.password, gotNewPassword, tt.wantNewPassword)
			}
		})
	}
}
