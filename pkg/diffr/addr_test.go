package diffr

import "testing"

func TestListenAndDisplay(t *testing.T) {
	tests := []struct {
		addr        string
		port        int
		wantListen  string
		wantDisplay string
	}{
		{"http://localhost", 8675, "localhost:8675", "http://localhost:8675"},
		{"http://127.0.0.1", 9000, "127.0.0.1:9000", "http://127.0.0.1:9000"},
		{"127.0.0.1", 8675, "127.0.0.1:8675", "http://127.0.0.1:8675"},
		{"0.0.0.0", 8675, "0.0.0.0:8675", "http://0.0.0.0:8675"},
		{"localhost:9000", 8675, "localhost:9000", "http://localhost:9000"},
	}
	for _, tt := range tests {
		listen, display, err := listenAndDisplay(tt.addr, tt.port)
		if err != nil {
			t.Fatalf("%q: %v", tt.addr, err)
		}
		if listen != tt.wantListen || display != tt.wantDisplay {
			t.Fatalf("%q: listen=%q display=%q, want %q %q", tt.addr, listen, display, tt.wantListen, tt.wantDisplay)
		}
	}
}
