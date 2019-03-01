package temporary

import (
	"io"
)

/**
 * 1. 超过一定大小、将临时缓冲转成临时文件
 * 2. TODO 异步读取源信息
 * 3. TODO 异步读取源信息后关闭
 **/

type Temporary interface {
	io.Reader
	io.Seeker
	io.Closer
}

type temporaryItemer interface {
	Size() int64
	io.Writer
	io.Reader
	io.Seeker
	io.Closer
}

type temporary struct {
	item        temporaryItemer
	itemConvert bool

	maxBufferSize int64
	fileDir       string
	filePattern   string
}

func NewTemporary(reader io.Reader, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary, err error) {
	var temp *temporary = &temporary{
		item:        newTemporaryBuffer(),
		itemConvert: false,

		maxBufferSize: maxBufferSize,
		fileDir:       fileDir,
		filePattern:   filePattern,
	}
	if _, err = io.Copy(temp, reader); err != nil {
		return
	}
	return temp, nil
}

func (temp *temporary) toTemporaryFile() (err error) {
	if temp.itemConvert {
		return
	}
	var (
		seekPosition int64
		tf           *temporaryFile
	)
	if seekPosition, err = temp.item.Seek(0, io.SeekCurrent); err != nil {
		return
	}
	if tf, err = newTemporaryFile(temp.fileDir, temp.filePattern); err != nil {
		return
	}
	if _, err = temp.item.Seek(0, io.SeekStart); err != nil {
		return
	}
	if _, err = io.Copy(tf, temp.item); err != nil {
		return
	}
	temp.item.Close()
	if _, err = tf.Seek(seekPosition, io.SeekStart); err != nil {
		return
	}
	temp.item = tf
	temp.itemConvert = true
	return
}

func (temp *temporary) Write(p []byte) (n int, err error) {
	if temp.item.Size()+int64(len(p)) > temp.maxBufferSize {
		if err = temp.toTemporaryFile(); err != nil {
			return
		}
	}
	if n, err = temp.item.Write(p); err != nil {
		if err == ErrBufferTooLarge { // 内存分配出错
			if err = temp.toTemporaryFile(); err != nil {
				return
			}
			var m int
			m, err = temp.item.Write(p[n:])
			n += m
			if err != nil {
				return
			}
			return
		}
		return
	}
	return
}

func (temp *temporary) Read(p []byte) (n int, err error) {
	return temp.item.Read(p)
}

func (temp *temporary) Seek(offset int64, whence int) (abs int64, err error) {
	return temp.item.Seek(offset, whence)
}

func (temp *temporary) Close() (err error) {
	return temp.item.Close()
}
