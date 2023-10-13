package wxapkg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
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

func Unpack(file, output string) error {
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

	fis := make([]*FileIndex, 0, header.FileCount)
	for i := 0; i < int(header.FileCount); i++ {
		// get filename length
		var namelen uint32
		if err := binary.Read(r, binary.BigEndian, &namelen); err != nil {
			return err
		}

		fi := new(FileIndex)
		// get filename
		n := make([]byte, namelen)
		_, err := io.ReadAtLeast(r, n, int(namelen))
		if err != nil {
			return err
		}

		fi.Name = string(n)
		if err = binary.Read(r, binary.BigEndian, &fi.Offset); err != nil {
			return err
		}
		if err = binary.Read(r, binary.BigEndian, &fi.Size); err != nil {
			return err
		}

		fis = append(fis, fi)
	}

	for _, fi := range fis {
		path := filepath.Join(output, fi.Name)
		if err := CheckDir(path); err != nil {
			return err
		}
		// get file contents
		content := SafetyGetData(data, fi.Offset, fi.Size)
		// format content
		content = Format(path, content)
		if err := os.WriteFile(path, content, 0o600); err != nil {
			return err
		}
	}

	return nil
}

func SafetyGetData(data []byte, offset, size uint32) []byte {
	if len(data) < (int(offset + size)) {
		return nil
	}
	return data[offset : offset+size]
}

func CheckDir(path string) error {
	path = filepath.Dir(path)
	stat, err := os.Stat(path)
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
	return os.MkdirAll(path, os.ModePerm)
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
