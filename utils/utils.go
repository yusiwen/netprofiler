package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func(source *os.File) {
		err := source.Close()
		if err != nil {
			log.Printf("failed to close [%v]: %v\n", source, err)
		}
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer func(destination *os.File) {
		err := destination.Close()
		if err != nil {
			log.Printf("failed to close [%v]: %v\n", destination, err)
		}
	}(destination)

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func Move(src, dst string) error {
	return os.Rename(src, dst)
}

// CopyDirectory recursively copies a src directory to a destination.
func CopyDirectory(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := createDir(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if _, err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		// `go test` fails on Windows even with this `if` supposedly
		// protecting the `syscall.Stat_t` and `os.Lchown` calls (not
		// available on windows). why?
		/*
			if runtime.GOOS != "windows" {
					stat, ok := fileInfo.Sys().(*syscall.Stat_t)
					if !ok {
						return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
					}
					if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
						return err
					}
			}
		*/

		isSymlink := entry.Type()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Type()); err != nil {
				return err
			}
		}
	}
	return nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func createDir(dir string, perm os.FileMode) error {
	if exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}
