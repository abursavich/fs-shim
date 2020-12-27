// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.16

// Package fs defines basic interfaces to a file system.
// A file system can be provided by the host operating system
// but also by other packages.
package fs

import (
	"io/fs"
	"os"
)

// fs.go

// An FS provides access to a hierarchical file system.
//
// The FS interface is the minimum implementation required of the file system.
// A file system may implement additional interfaces,
// such as ReadFileFS, to provide additional or optimized functionality.
type FS = fs.FS

// A File provides access to a single file.
// The File interface is the minimum implementation required of the file.
// A file may implement additional interfaces, such as
// ReadDirFile, ReaderAt, or Seeker, to provide additional or optimized functionality.
type File = fs.File

// A ReadDirFile is a directory file whose entries can be read with the ReadDir method.
// Every directory file should implement this interface.
// (It is permissible for any file to implement this interface,
// but if so ReadDir should return an error for non-directories.)
type ReadDirFile = fs.ReadDirFile

// A DirEntry is an entry read from a directory
// (using the ReadDir function or a ReadDirFile's ReadDir method).
type DirEntry = fs.DirEntry

// A FileInfo describes a file and is returned by Stat.
type FileInfo = fs.FileInfo

// A FileMode represents a file's mode and permission bits.
// The bits have the same definition on all systems, so that
// information about files can be moved from one system
// to another portably. Not all bits apply to all systems.
// The only required bit is ModeDir for directories.
type FileMode = fs.FileMode

// PathError records an error and the operation and file path that caused it.
type PathError = fs.PathError

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
	return fs.ValidPath(name)
}

// Generic file system errors.
// Errors returned by file systems can be tested against these errors
// using errors.Is.
var (
	ErrInvalid    = fs.ErrInvalid    // "invalid argument"
	ErrPermission = fs.ErrPermission // "permission denied"
	ErrExist      = fs.ErrExist      // "file already exists"
	ErrNotExist   = fs.ErrNotExist   // "file does not exist"
	ErrClosed     = fs.ErrClosed     // "file already closed"
)

// The defined file mode bits are the most significant bits of the FileMode.
// The nine least-significant bits are the standard Unix rwxrwxrwx permissions.
// The values of these bits should be considered part of the public API and
// may be used in wire protocols or disk representations: they must not be
// changed, although new bits might be added.
const (
	// The single letters are the abbreviations
	// used by the String method's formatting.
	ModeDir        = fs.ModeDir        // d: is a directory
	ModeAppend     = fs.ModeAppend     // a: append-only
	ModeExclusive  = fs.ModeExclusive  // l: exclusive use
	ModeTemporary  = fs.ModeTemporary  // T: temporary file; Plan 9 only
	ModeSymlink    = fs.ModeSymlink    // L: symbolic link
	ModeDevice     = fs.ModeDevice     // D: device file
	ModeNamedPipe  = fs.ModeNamedPipe  // p: named pipe (FIFO)
	ModeSocket     = fs.ModeSocket     // S: Unix domain socket
	ModeSetuid     = fs.ModeSetuid     // u: setuid
	ModeSetgid     = fs.ModeSetgid     // g: setgid
	ModeCharDevice = fs.ModeCharDevice // c: Unix character device, when ModeDevice is set
	ModeSticky     = fs.ModeSticky     // t: sticky
	ModeIrregular  = fs.ModeIrregular  // ?: non-regular file; nothing else is known about this file

	ModeType = fs.ModeType // Mask for the type bits. For regular files, none will be set.

	ModePerm = fs.ModePerm // Unix permission bits
)

// glob.go

// A GlobFS is a file system with a Glob method.
type GlobFS = fs.GlobFS

// Glob returns the names of all files matching pattern or nil
// if there is no matching file. The syntax of patterns is the same
// as in path.Match. The pattern may describe hierarchical names such as
// /usr/*/bin/ed (assuming the Separator is '/').
//
// Glob ignores file system errors such as I/O errors reading directories.
// The only possible returned error is path.ErrBadPattern, reporting that
// the pattern is malformed.
//
// If fs implements GlobFS, Glob calls fs.Glob.
// Otherwise, Glob uses ReadDir to traverse the directory tree
// and look for matches for the pattern.
func Glob(fsys FS, pattern string) (matches []string, err error) {
	return fs.Glob(fsys, pattern)
}

// readdir.go

// ReadDirFS is the interface implemented by a file system
// that provides an optimized implementation of ReadDir.
type ReadDirFS = fs.ReadDirFS

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
//
// If fs implements ReadDirFS, ReadDir calls fs.ReadDir.
// Otherwise ReadDir calls fs.Open and uses ReadDir and Close
// on the returned file.
func ReadDir(fsys FS, name string) ([]DirEntry, error) {
	return fs.ReadDir(fsys, name)
}

// readfile.go

// ReadFileFS is the interface implemented by a file system
// that provides an optimized implementation of ReadFile.
type ReadFileFS = fs.ReadFileFS

// ReadFile reads the named file from the file system fs and returns its contents.
// A successful call returns a nil error, not io.EOF.
// (Because ReadFile reads the whole file, the expected EOF
// from the final Read is not treated as an error to be reported.)
//
// If fs implements ReadFileFS, ReadFile calls fs.ReadFile.
// Otherwise ReadFile calls fs.Open and uses Read and Close
// on the returned file.
func ReadFile(fsys FS, name string) ([]byte, error) {
	return fs.ReadFile(fsys, name)
}

// stat.go

// A StatFS is a file system with a Stat method.
type StatFS = fs.StatFS

// Stat returns a FileInfo describing the named file from the file system.
//
// If fs implements StatFS, Stat calls fs.Stat.
// Otherwise, Stat opens the file to stat it.
func Stat(fsys FS, name string) (FileInfo, error) {
	return fs.Stat(fsys, name)
}

// sub.go

// A SubFS is a file system with a Sub method.
type SubFS = fs.SubFS

// Sub returns an FS corresponding to the subtree rooted at fsys's dir.
//
// If fs implements SubFS, Sub calls returns fsys.Sub(dir).
// Otherwise, if dir is ".", Sub returns fsys unchanged.
// Otherwise, Sub returns a new FS implementation sub that,
// in effect, implements sub.Open(dir) as fsys.Open(path.Join(dir, name)).
// The implementation also translates calls to ReadDir, ReadFile, and Glob appropriately.
//
// Note that Sub(os.DirFS("/"), "prefix") is equivalent to os.DirFS("/prefix")
// and that neither of them guarantees to avoid operating system
// accesses outside "/prefix", because the implementation of os.DirFS
// does not check for symbolic links inside "/prefix" that point to
// other directories. That is, os.DirFS is not a general substitute for a
// chroot-style security mechanism, and Sub does not change that fact.
func Sub(fsys FS, dir string) (FS, error) {
	return fs.Sub(fsys, dir)
}

// walk.go

// SkipDir is used as a return value from WalkDirFuncs to indicate that
// the directory named in the call is to be skipped. It is not returned
// as an error by any function.
var SkipDir = fs.SkipDir

// WalkDirFunc is the type of the function called by WalkDir to visit
// each each file or directory.
//
// The path argument contains the argument to Walk as a prefix.
// That is, if Walk is called with root argument "dir" and finds a file
// named "a" in that directory, the walk function will be called with
// argument "dir/a".
//
// The directory and file are joined with Join, which may clean the
// directory name: if Walk is called with the root argument "x/../dir"
// and finds a file named "a" in that directory, the walk function will
// be called with argument "dir/a", not "x/../dir/a".
//
// The d argument is the fs.DirEntry for the named path.
//
// The error result returned by the function controls how WalkDir
// continues. If the function returns the special value SkipDir, WalkDir
// skips the current directory (path if d.IsDir() is true, otherwise
// path's parent directory). Otherwise, if the function returns a non-nil
// error, WalkDir stops entirely and returns that error.
//
// The err argument reports an error related to path, signaling that
// WalkDir will not walk into that directory. The function can decide how
// to handle that error; as described earlier, returning the error will
// cause WalkDir to stop walking the entire tree.
//
// WalkDir calls the function with a non-nil err argument in two cases.
//
// First, if the initial os.Lstat on the root directory fails, WalkDir
// calls the function with path set to root, d set to nil, and err set to
// the error from os.Lstat.
//
// Second, if a directory's ReadDir method fails, WalkDir calls the
// function with path set to the directory's path, d set to an
// fs.DirEntry describing the directory, and err set to the error from
// ReadDir. In this second case, the function is called twice with the
// path of the directory: the first call is before the directory read is
// attempted and has err set to nil, giving the function a chance to
// return SkipDir and avoid the ReadDir entirely. The second call is
// after a failed ReadDir and reports the error from ReadDir.
// (If ReadDir succeeds, there is no second call.)
//
// The differences between WalkDirFunc compared to filepath.WalkFunc are:
//
//   - The second argument has type fs.DirEntry instead of fs.FileInfo.
//   - The function is called before reading a directory, to allow SkipDir
//     to bypass the directory read entirely.
//   - If a directory read fails, the function is called a second time
//     for that directory to report the error.
//
type WalkDirFunc = fs.WalkDirFunc

// WalkDir walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root.
//
// All errors that arise visiting files and directories are filtered by fn:
// see the fs.WalkDirFunc documentation for details.
//
// The files are walked in lexical order, which makes the output deterministic
// but requires WalkDir to read an entire directory into memory before proceeding
// to walk that directory.
//
// WalkDir does not follow symbolic links found in directories,
// but if root itself is a symbolic link, its target will be walked.
func WalkDir(fsys FS, root string, fn WalkDirFunc) error {
	return fs.WalkDir(fsys, root, fn)
}

// os

// DirFS returns a file system for the tree of files rooted at the directory dir.
//
// Note that DirFS("/prefix") only guarantees that the Open calls it makes to the
// operating system will begin with "/prefix": DirFS("/prefix").Open("file") is the
// same as os.Open("/prefix/file"). So if /prefix/file is a symbolic link pointing outside
// the /prefix tree, then using DirFS does not stop the access any more than using
// os.Open does. DirFS is therefore not a general substitute for a chroot-style security
// mechanism when the directory tree contains arbitrary content.
func DirFS(dir string) FS {
	return os.DirFS(dir)
}
