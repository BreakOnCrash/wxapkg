package wxapkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

const (
	WXAPkgMagic    = 0xBE
	WXAPkgEndMagic = 0xED
)

type Header struct {
	Magic       byte
	Reserved    uint32 // 预留的4字节
	IndexLength uint32 // 索引长度
	BodyLength  uint32 // 数据长度
	EndMagic    byte
	FileCount   uint32 // 文件总数
}

type FileIndex struct {
	Name   string
	Offset uint32
	Size   uint32
}

var (
	ErrInvalidWXAPkg = errors.New("invalid wxapkg file")
)

func Unpack(file, output string, format bool, printf func(format string, a ...interface{})) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	r := bytes.NewReader(data)
	header := new(Header)
	if err := binary.Read(r, binary.BigEndian, header); err != nil {
		return err
	}

	if header.Magic != WXAPkgMagic ||
		header.EndMagic != WXAPkgEndMagic ||
		header.Reserved != 0 {
		return ErrInvalidWXAPkg
	}

	fis := make(chan *FileIndex)
	go func() {
		defer close(fis)

		for i := 0; i < int(header.FileCount); i++ {
			// get filename length
			var namelen uint32
			if err = binary.Read(r, binary.BigEndian, &namelen); err != nil {
				return
			}

			fi := new(FileIndex)
			// get filename
			n := make([]byte, namelen)
			_, err := io.ReadAtLeast(r, n, int(namelen))
			if err != nil {
				return
			}

			fi.Name = string(n)
			if err = binary.Read(r, binary.BigEndian, &fi.Offset); err != nil {
				return
			}
			if err = binary.Read(r, binary.BigEndian, &fi.Size); err != nil {
				return
			}

			fis <- fi
		}
	}()

	checkDir := CacheCheckDir()
	for fi := range fis {
		path := filepath.Join(output, fi.Name)
		if err = checkDir(path); err != nil {
			return err
		}
		// get file contents
		content := SafetyGetData(data, fi.Offset, fi.Size)
		// format content
		if format {
			content = Format(path, content)
		}
		if printf != nil {
			printf("\t file:%s, content-length:%d\n", fi.Name, len(content))
		}
		if err := os.WriteFile(path, content, 0o600); err != nil {
			return err
		}
	}
	return err
}

func SafetyGetData(data []byte, offset, size uint32) []byte {
	if len(data) < (int(offset + size)) {
		return nil
	}
	return data[offset : offset+size]
}

func CacheCheckDir() func(path string) error {
	var cache sync.Map

	return func(path string) (err error) {
		if _, ok := cache.Load(path); ok { // path checked
			return nil
		}

		defer func() {
			if err == nil {
				cache.Store(path, struct{}{})
			}
		}()

		var stat fs.FileInfo
		path = filepath.Dir(path)
		stat, err = os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				goto NotExist
			}
			return err
		}

		if stat.IsDir() {
			return nil
		}

	NotExist:
		err = os.MkdirAll(path, os.ModePerm)
		return
	}
}

func Format(ext string, data []byte) []byte {
	ext = filepath.Ext(ext)
	switch ext {
	case ".json":
		return PrettyJson(data)
	case ".html", ".htm":
		// remove last byte "0x00"
		if bytes.HasSuffix(data, []byte{0x00}) {
			data = data[:len(data)-1]
		}
		return PrettyHtml(data)
	case ".js":
		return PrettyJavaScript(data)
	default:
		return data
	}
}
