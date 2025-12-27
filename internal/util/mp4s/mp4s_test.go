package mp4s_test

import (
	"os"
	"testing"
	"time"

	"github.com/syscc/Emby-Go/internal/util/mp4s"
)

func TestGenWithDuration(t *testing.T) {
	d := time.Hour*2 + time.Minute*23 + time.Second*21 + time.Millisecond*90
	bytes := mp4s.GenWithDuration(d)
	os.WriteFile("test.mp4", bytes, os.ModePerm)
}
