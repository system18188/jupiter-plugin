package filepath

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Excludes is handly filepath match manipulation
type Excludes []string

// Validate helps to evalute the pattern are valid
// `Match` method is from error and focus on match
func (e *Excludes) Validate() error {
	for _, pattern := range *e {
		if _, err := filepath.Match(pattern, "abc/b/c"); err != nil {
			return fmt.Errorf("unable to evalute pattern: %v", pattern)
		}
	}
	return nil
}

// Match evalutes given file with available patterns returns true if matches
// otherwise false. `Match` internally uses the `filepath.Match`
//
// Note: `Match` ignore pattern errors, use `Validate` method to ensure
// you have correct exclude patterns
func (e *Excludes) Match(file string) bool {
	for _, pattern := range *e {
		if match, _ := filepath.Match(pattern, file); match {
			return match
		}
	}
	return false
}

// 检查文件是否存在
func Exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// IsFileExists return true is file or directory is exists, otherwise returns
// false. It also take cares of symlink path as well
func IsFileExists(filename string) bool {
	_, err := os.Lstat(filename)
	return err == nil
}

// IsDirEmpty returns true if the given directory is empty also returns true if
// directory not exists. Otherwise returns false
func IsDirEmpty(path string) bool {
	if !IsFileExists(path) {
		// directory not exists
		return true
	}
	dir, _ := os.Open(path)
	defer CloseQuietly(dir)
	results, _ := dir.Readdir(1)
	return len(results) == 0
}

// IsDir returns true if the given `path` is directory otherwise returns false.
// Also returns false if path is not exists
func IsDir(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.IsDir()
}

// ApplyFileMode applies the given file mode to the target{file|directory}
func ApplyFileMode(target string, mode os.FileMode) error {
	err := os.Chmod(target, mode)
	if err != nil {
		return fmt.Errorf("unable to apply mode: %v, to given file or directory: %v", mode, target)
	}
	return nil
}

// LineCnt counts no. of lines on file
func LineCnt(fileName string) int {
	f, err := os.Open(fileName)
	if err != nil {
		return 0
	}
	defer CloseQuietly(f)

	return LineCntr(f)
}

// LineCntr counts no. of lines for given reader
func LineCntr(r io.Reader) int {
	buf := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return count
		}

		count += bytes.Count(buf[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count
}

// Walk method extends filepath.Walk to also follows symlinks.
// Always returns the path of the file or directory also path
// is inline to name of symlink
func Walk(srcDir string, walkFn filepath.WalkFunc) error {
	return doWalk(srcDir, srcDir, walkFn)
}

func doWalk(fname string, linkName string, walkFn filepath.WalkFunc) error {
	fsWalkFn := func(path string, info os.FileInfo, err error) error {
		var name string
		name, err = filepath.Rel(fname, path)
		if err != nil {
			return err
		}

		path = filepath.Join(linkName, name)

		if err == nil && info.Mode()&os.ModeSymlink == os.ModeSymlink {
			var symlinkPath string
			symlinkPath, err = filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			// https://github.com/golang/go/blob/master/src/path/filepath/path.go#L392
			info, err = os.Lstat(symlinkPath)
			if err != nil {
				return walkFn(path, info, err)
			}

			if info.IsDir() {
				return doWalk(symlinkPath, path, walkFn)
			}
		}

		return walkFn(path, info, err)
	}

	return filepath.Walk(fname, fsWalkFn)
}

// CopyFile copies the given source file into destination
func CopyFile(dest, src string) (int64, error) {
	if !IsFileExists(src) {
		return 0, fmt.Errorf("source file does not exists: %v", src)
	}

	baseName := filepath.Base(src)
	if !strings.HasSuffix(dest, baseName) {
		dest = filepath.Join(dest, baseName)
	}

	if IsFileExists(dest) {
		return 0, fmt.Errorf("destination file already exists: %v", dest)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return 0, fmt.Errorf("unable to create dest file: %v", dest)
	}
	defer CloseQuietly(destFile)

	srcFile, err := os.Open(src)
	if err != nil {
		return 0, fmt.Errorf("unable to open source file: %v", src)
	}
	defer CloseQuietly(srcFile)

	copiedBytes, err := io.Copy(destFile, srcFile)
	if err != nil {
		return 0, fmt.Errorf("unable to copy file from %v to %v (%v)",
			src, dest, err)
	}

	return copiedBytes, nil
}

// CopyDir copies entire directory, sub directories and files into destination
// and it excludes give file matches
func CopyDir(dest, src string, excludes Excludes) error {
	if !IsFileExists(src) {
		return fmt.Errorf("source dir does not exists: %v", src)
	}

	src = filepath.Clean(src)
	srcInfo, _ := os.Lstat(src)
	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not directory: %v", src)
	}

	baseName := filepath.Base(src)
	if !strings.HasSuffix(dest, baseName) {
		dest = filepath.Join(dest, baseName)
	}

	if IsFileExists(dest) {
		return fmt.Errorf("destination dir already exists: %v", dest)
	}

	if err := excludes.Validate(); err != nil {
		return err
	}

	return Walk(src, func(srcPath string, info os.FileInfo, err error) error {
		if excludes.Match(filepath.Base(srcPath)) {
			if info.IsDir() {
				// excluding directory
				return filepath.SkipDir
			}
			// excluding file
			return nil
		}

		relativeSrcPath := strings.TrimLeft(srcPath[len(src):], string(filepath.Separator))
		destPath := filepath.Join(dest, relativeSrcPath)

		if info.IsDir() {
			// directory permissions is not preserved from source
			return MkDirAll(destPath, 0755)
		}

		// copy source into destination
		if _, err = CopyFile(destPath, srcPath); err != nil {
			return err
		}

		// Apply source permision into target as well
		// so file permissions are preserved
		return ApplyFileMode(destPath, info.Mode())
	})
}

// DeleteFiles method deletes give files or directories. ensure your supplying
// appropriate paths.
func DeleteFiles(files ...string) (errs []error) {
	for _, f := range files {
		if !IsFileExists(f) {
			errs = append(errs, fmt.Errorf("path does not exists: %s", f))
			continue
		}

		var err error
		if IsDir(f) {
			err = os.RemoveAll(f)
		} else {
			err = os.Remove(f)
		}

		if err != nil {
			errs = append(errs, err)
		}
	}

	return
}

// DirsPath method returns all directories absolute path from given base path recursively.
func DirsPath(basePath string, recursive bool) (pdirs []string, err error) {
	if recursive {
		err = Walk(basePath, func(srcPath string, info os.FileInfo, err error) error {
			if info.IsDir() {
				pdirs = append(pdirs, srcPath)
			}
			return nil
		})
		return
	}

	var list []os.FileInfo
	list, err = ioutil.ReadDir(basePath)
	if err != nil {
		return
	}

	for _, v := range list {
		if v.IsDir() {
			pdirs = append(pdirs, filepath.Join(basePath, v.Name()))
		}
	}

	return
}

// FilesPath method returns all files absolute path from given base path recursively.
func FilesPath(basePath string, recursive bool) (files []string, err error) {
	if recursive {
		err = Walk(basePath, func(srcPath string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, srcPath)
			}
			return nil
		})
		return
	}

	var list []os.FileInfo
	list, err = ioutil.ReadDir(basePath)
	if err != nil {
		return
	}

	for _, v := range list {
		if !v.IsDir() {
			files = append(files, filepath.Join(basePath, v.Name()))
		}
	}

	return
}

// StripExt method returns name of the file without extension.
//    E.g.: index.html => index
func StripExt(name string) string {
	if IsStrEmpty(name) {
		return name
	}

	idx := strings.LastIndexByte(name, '.')
	if idx > 0 {
		return name[:idx]
	}

	return name
}

// IsStrEmpty returns true if strings is empty otherwise false
func IsStrEmpty(v string) bool {
	return len(strings.TrimSpace(v)) == 0
}
