# temporary
[godoc](https://godoc.org/github.com/linfangrong/temporary)

## 场景
将io.Reader转换成io.ReadSeeker。
1. 数据大小超过maxBufferSize时, 写入到临时文件ioutil.TempFile(fileDir, filePattern).
2. 提供异步操作流程NewAsyncTemporary/NewMustCloseReaderAsyncTemporary, 在调用其他方法之前一定要先使用Await().
3. 数据结果操作.
	* Size()  数据大小
	* Type()  TemporaryBuffer | TemporaryFile
	* Name()  当类型为文件(TemporaryFile)时, 返回临时文件名.
	* Bytes() 当类型为缓存(TemporaryBuffer)时, 返回[]byte.

## temporary.Temporary
```go
type Temporary interface {
	Await() error
	Size() int64
	Type() string
	Name() string
	Bytes() []byte
	io.Reader
	io.Seeker
	io.Closer
}
```

## 同步方式
```go
func NewTemporary(reader io.Reader, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary, err error)
```

### Example
```go
package main

import (
	"github.com/linfangrong/temporary"
)

func main() {
	var (
		temp temporary.Temporary
		err  error
	)
	if temp, err = temporary.NewTemporary(
		reader,      // io.Reader
		2*1024*1204, // maxBufferSize 超过这个长度, 会转化成临时文件
		"/tmp",      // fileDir 临时文件目录
		"*",         // filePattern 文件名生成方式
	); err != nil {
		// TODO error
		return
	}
	defer temp.Close()
	//
	// temp 可以当做io.ReadSeeker
	// Size() int64
	// io.Reader
	// io.Seeker
	// io.Closer
	//
	// 判断存储方式的类型
	// Type() string
	// Name() string   Type: TemporaryFile
	// Bytes() []byte  Type: TemporaryBuffer

	...
}
```

## 异步方式
```go
func NewAsyncTemporary(reader io.Reader, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary)
func NewMustCloseReaderAsyncTemporary(readcloser io.ReadCloser, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary)
```
