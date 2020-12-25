package chromePwd

import (
	"testing"
)

func TestGetChromePwd(t *testing.T) {
	_, err := GetChromePwd()
	t.Log(err)
}

func TestEncrypt(t *testing.T) {
	Enc, err := Encrypt([]byte("123456"))
	t.Log(Enc, err)
	dec, err := Decrypt(Enc)
	t.Log(string(dec), err)
}
