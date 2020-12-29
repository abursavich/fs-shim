// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fs defines basic interfaces to a file system.
// A file system can be provided by the host operating system
// but also by other packages.
package fs

import (
	"os"
)

// An FS provides access to a hierarchical file system.
//
// The FS interface is the minimum implementation required of the file system.
// A file system may implement additional interfaces,
// such as ReadFileFS, to provide additional or optimized functionality.
type FS interface {
	// Open opens the named file.
	//
	// When Open returns an error, it should be of type *PathError
	// with the Op field set to "open", the Path field set to name,
	// and the Err field describing the problem.
	//
	// Open should reject attempts to open names that do not satisfy
	// ValidPath(name), returning a *PathError with Err set to
	// ErrInvalid or ErrNotExist.
	Open(name string) (File, error)
}

// ValidPath reports whether the given path name
// is valid for use in a call to Open.
// Path names passed to open are unrooted, slash-separated
// sequences of path elements, like “x/y/z”.
// Path names must not contain a “.” or “..” or empty element,
// except for the special case that the root directory is named “.”.
//
// Paths are slash-separated on all systems, even Windows.
// Backslashes must not appear in path names.
func ValidPath(name string) bool {
	if name == "." {
		// special case
		return true
	}

	// Iterate over elements in name, checking each.
	for {
		i := 0
		for i < len(name) && name[i] != '/' {
			if name[i] == '\\' {
				return false
			}
			i++
		}
		elem := name[:i]
		if elem == "" || elem == "." || elem == ".." {
			return false
		}
		if i == len(name) {
			return true // reached clean ending
		}
		name = name[i+1:]
	}
}

// A File provides access to a single file.
// The File interface is the minimum implementation required of the file.
// A file may implement additional interfaces, such as
// ReadDirFile, ReaderAt, or Seeker, to provide additional or optimized functionality.
type File interface {
	Stat() (FileInfo, error)
	Read([]byte) (int, error)
	Close() error
}

// The defined file mode bits are the most significant bits of the FileMode.
// The nine least-significant bits are the standard Unix rwxrwxrwx permissions.
// The values of these bits should be considered part of the public API and
// may be used in wire protocols or disk representations: they must not be
// changed, although new bits might be added.
const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        = FileMode(os.ModeDir)        // d: is a directory
	ModeAppend     = FileMode(os.ModeAppend)     // a: append-only
	ModeExclusive  = FileMode(os.ModeExclusive)  // l: exclusive use
	ModeTemporary  = FileMode(os.ModeTemporary)  // T: temporary file; Plan 9 only
	ModeSymlink    = FileMode(os.ModeSymlink)    // L: symbolic link
	ModeDevice     = FileMode(os.ModeDevice)     // D: device file
	ModeNamedPipe  = FileMode(os.ModeNamedPipe)  // p: named pipe (FIFO)
	ModeSocket     = FileMode(os.ModeSocket)     // S: Unix domain socket
	ModeSetuid     = FileMode(os.ModeSetuid)     // u: setuid
	ModeSetgid     = FileMode(os.ModeSetgid)     // g: setgid
	ModeCharDevice = FileMode(os.ModeCharDevice) // c: Unix character device, when ModeDevice is set
	ModeSticky     = FileMode(os.ModeSticky)     // t: sticky
	ModeIrregular  = FileMode(os.ModeIrregular)  // ?: non-regular file; nothing else is known about this file

	ModeType = FileMode(os.ModeType) // Mask for the type bits. For regular files, none will be set.

	ModePerm = FileMode(os.ModePerm) // Unix permission bits
)

// Generic file system errors.
// Errors returned by file systems can be tested against these errors
// using errors.Is.
var (
	ErrInvalid    = os.ErrInvalid    // "invalid argument"
	ErrPermission = os.ErrPermission // "permission denied"
	ErrExist      = os.ErrExist      // "file already exists"
	ErrNotExist   = os.ErrNotExist   // "file does not exist"
	ErrClosed     = os.ErrClosed     // "file already closed"
)

// PathError records an error and the operation and file path that caused it.
type PathError = os.PathError
