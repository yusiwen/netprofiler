package profiler

import (
	"log"
	"os"
	"path/filepath"

	"git.yusiwen.cn/yusiwen/netprofiler/utils"
)

type PostLoadFunc func() error

type File struct {
	Path          string `json:"path"`
	RootPrivilege bool   `json:"root-privilege"`
}

type FileProfiler struct {
	Name     string       `json:"name"`
	Files    []File       `json:"files"`
	PostLoad PostLoadFunc `json:"-"`
}

func (fp *FileProfiler) Save(profile, location string) error {
	for _, f := range fp.Files {
		dstPath := filepath.Join(location, profile, fp.Name)
		os.MkdirAll(dstPath, os.ModePerm)
		dst := filepath.Join(dstPath, filepath.Base(f.Path))
		_, err := utils.Copy(f.Path, dst)
		if err != nil {
			return err
		}
		log.Printf("save '%s' to '%s'\n", f.Path, dst)
	}
	return nil
}

func (fp *FileProfiler) Load(profile, location string) error {
	for _, f := range fp.Files {
		srcPath := filepath.Join(location, profile, fp.Name)
		src := filepath.Join(srcPath, filepath.Base(f.Path))
		var err error
		if f.RootPrivilege {
			err = utils.CopySudo(src, f.Path)
		} else {
			_, err = utils.Copy(src, f.Path)
		}
		if err != nil {
			return err
		}
		log.Printf("load '%s' to '%s'\n", src, f.Path)
	}
	if fp.PostLoad != nil {
		err := fp.PostLoad()
		if err != nil {
			return err
		}
	}
	return nil
}
