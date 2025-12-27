package openlist_test

import (
	"log"
	"testing"

	"github.com/syscc/Emby-Go/internal/config"
	"github.com/syscc/Emby-Go/internal/service/openlist"
)

func TestWalkFsList(t *testing.T) {
	err := config.ReadFromFile("../../../config.yml")
	if err != nil {
		t.Fatal(err)
		return
	}

	walker := openlist.WalkFsList("/", 4)
	page, err := walker.Next()
	for err == nil {
		log.Println("page: ", page)
		page, err = walker.Next()
	}
	if err == openlist.ErrWalkEOF {
		return
	}
	t.Fatal(err)
}
