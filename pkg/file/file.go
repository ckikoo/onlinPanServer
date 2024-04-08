package file

import "io"

// AbstractFile 是一个抽象文件类，包含了 Writer、Reader 和 Seeker 接口的方法
type AbstractFile struct {
	io.Reader
	io.Seeker
	io.Closer
}
