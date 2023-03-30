package chain

import "testing"

func TestGetChainList(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(GetChain(2).Name)
		})
	}
}
