package binding

import (
	"errors"
	"fmt"
	"jupiter-plugin/pkg/filepath"
	"jupiter-plugin/pkg/random"
	"io"
	"mime/multipart"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

type File struct {
	*multipart.FileHeader
}

type Files []File

// 验证图片
func (f File) IsImage() bool {
	fileType := f.Header[HEADER_ContentType][0]
	if fileType == MIMEGIF || fileType == MIMEPNG || fileType == MIMEJPEN {
		return true
	}
	return false
}

// 验证文档
func (f File) IsDoc() bool {
	fileType := f.Header[HEADER_ContentType][0]
	if fileType == MIMEDOC || fileType == MIMEDOCX || fileType == MIMEPDF {
		return true
	}
	return false
}

// 验证文件包
func (f File) IsRarOr7z() bool {
	fileType := f.Header[HEADER_ContentType][0]
	if fileType == MIME7z || fileType == MIMErar {
		return true
	}
	return false
}

// 保存文件目录 (反回文件路径)
func (f File) SaveUploadedFile(dst string) (file string, err error) {
	src, err := f.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	hType, ok := MIMETypes[f.Header[HEADER_ContentType][0]]
	if !ok || hType == "" {
		return "", errors.New("SaveUploadedFile Type then nil")
	}
	// 生成文件路径
	t := time.Now()
	file = fmt.Sprint(t.Year(), "/", t.YearDay(), t.Hour(), t.Minute(), t.Second(), random.Intn(5), hType)
	localFile := dst + file
	// 建立所有文件夹
	if err := os.MkdirAll(path.Dir(localFile), 0700); err != nil {
		return file, err
	}
	// 保存文件
	out, err := os.Create(localFile)
	if err != nil {
		return file, err
	}
	defer out.Close()

	io.Copy(out, src)
	return file, nil
}

// 删除上传的文件
func (f File) RemoveFile(dst string, files ...string) []error {
	for i, file := range files {
		files[i] = dst + file
	}
	return filepath.DeleteFiles(files...)
}

func (f Files) Len() int {
	return len(f)
}

func (f Files) Size() int64 {
	var size int64
	for _, file := range f {
		size = size + file.Size
	}
	return size
}

var tFile = reflect.TypeOf(File{})
var tFiles = reflect.TypeOf(Files{})

func mapFile(ptr interface{}, file map[string][]*multipart.FileHeader) error {
	return mapFileByTag(ptr, file, "form")
}

func mapFileByTag(ptr interface{}, file map[string][]*multipart.FileHeader, tag string) error {
	typ := reflect.TypeOf(ptr).Elem()
	val := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}
		inputFieldName := typeField.Tag.Get(tag)
		inputFieldNameList := strings.Split(inputFieldName, ",")
		inputFieldName = inputFieldNameList[0]

		inputValue := reflect.ValueOf(file[inputFieldName])
		if inputValue.Len() < 1 {
			continue
		}
		inputFieldValue := inputValue.Index(0)
		if !inputFieldValue.IsValid() {
			continue
		}
		if typeField.Type.ConvertibleTo(tFile) {
			structField.Field(0).Set(inputFieldValue)
			continue
		}

		if typeField.Type.ConvertibleTo(tFiles) {
			vfiles := make([]File, 0)
			for i = 0; i < inputValue.Len(); i++ {
				vfile := File{}
				rvFile := reflect.ValueOf(vfile)
				rvFile.Set(inputValue.Index(i))
				vfiles = append(vfiles, vfile)
			}
			files := Files{}
			files = vfiles
			structField.Set(reflect.ValueOf(files))
		}
	}
	return nil
}
