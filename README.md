# temporary
[godoc](https://godoc.org/github.com/linfangrong/temporary)

## 同步方式
```go
func NewAsyncTemporary(reader io.Reader, maxBufferSize int64, fileDir string, filePattern string) (_ Temporary)
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
