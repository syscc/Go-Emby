package localtree_test

import (
	"testing"

	"github.com/syscc/Emby-Go/internal/service/openlist/localtree"
)

func TestSynchronizer_InitSnapshot(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "1", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := localtree.NewSynchronizer("../../../../openlist-local-tree", 50)
			if err := s.InitSnapshot(); (err != nil) != tt.wantErr {
				t.Errorf("Synchronizer.InitSnapshot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
