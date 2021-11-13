package profiler

import (
	"log"
	"os"
	"path/filepath"

	"git.yusiwen.cn/yusiwen/netprofiler/utils"
)

type PostLoadFunc func() error

type FileProfiler struct {
	Name     string
	Files    []string
	PostLoad PostLoadFunc
}

func (fp *FileProfiler) Save(profile, location string) error {
	for _, f := range fp.Files {
		dstPath := filepath.Join(location, profile, fp.Name)
		os.MkdirAll(dstPath, os.ModePerm)
		dst := filepath.Join(dstPath, filepath.Base(f))
		_, err := utils.Copy(f, dst)
		if err != nil {
			return err
		}
		log.Printf("save '%s' to '%s'\n", f, dst)
	}
	return nil
}

func (fp *FileProfiler) Load(profile, location string) error {
	for _, f := range fp.Files {
		srcPath := filepath.Join(location, profile, fp.Name)
		src := filepath.Join(srcPath, filepath.Base(f))
		err := utils.Move(src, f)
		if err != nil {
			return err
		}
	}
	if fp.PostLoad != nil {
		err := fp.PostLoad()
		if err != nil {
			return err
		}
	}
	return nil
}
