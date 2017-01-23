package apiserver

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   ``,
		DirModTime: time.Unix(1484979746, 0),
		ChildFiles: []*embedded.EmbeddedFile{},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`files/insp`, &embedded.EmbeddedBox{
		Name: `files/insp`,
		Time: time.Unix(1484979746, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{},
	})
}
