package https

import (
	"testing"
)

func TestNewCACert(t *testing.T) {
	cert, err := NewCACert("test.shanhu.io")
	if err != nil {
		t.Fatalf("NewCACert() got error: %s", err)
	}

	if _, err := cert.X509KeyPair(); err != nil {
		t.Fatalf("convert to tls cert got error: %s", err)
	}
}

func TestMakeRSACertWithNoHost(t *testing.T) {
	_, err := MakeRSACert(new(RSACertConfig))
	if err == nil {
		t.Errorf("expect error with not host, got nil")
	}
}
