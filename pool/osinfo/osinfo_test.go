package osinfo

import "testing"

func TestGetOSInformation(t *testing.T) {
	t.Log(GetOSInformation())
}
