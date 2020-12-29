// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fs

import (
	"os"
)

// DirFS returns a file system for the tree of files rooted at the directory dir.
//
// Note that DirFS("/prefix") only guarantees that the Open calls it makes to the
// operating system will begin with "/prefix": DirFS("/prefix").Open("file") is the
// same as os.Open("/prefix/file"). So if /prefix/file is a symbolic link pointing outside
// the /prefix tree, then using DirFS does not stop the access any more than using
// os.Open does. DirFS is therefore not a general substitute for a chroot-style security
// mechanism when the directory tree contains arbitrary content.
func DirFS(dir string) FS {
	return dirFS(dir)
}

type dirFS string

func (dir dirFS) Open(name string) (File, error) {
	if !ValidPath(name) {
		return nil, &PathError{Op: "open", Path: name, Err: ErrInvalid}
	}
	f, err := os.Open(string(dir) + "/" + name)
	if err != nil {
		return nil, err // nil fs.File
	}
	return &file{f}, nil
}

type file struct {
	*os.File
}

func (f *file) Stat() (FileInfo, error) {
	if f == nil {
		return nil, ErrInvalid
	}
	fi, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return fileInfo{fi}, nil
}

func (f *file) ReadDir(n int) ([]DirEntry, error) {
	if f == nil {
		return nil, ErrInvalid
	}
	fileInfos, err := f.File.Readdir(n)
	dirEntries := make([]DirEntry, len(fileInfos))
	for i, fi := range fileInfos {
		dirEntries[i] = dirEntry{info: fileInfo{fi}}
	}
	return dirEntries, err
}

type dirEntry struct {
	info fileInfo
}

func (d dirEntry) Name() string            { return d.info.Name() }
func (d dirEntry) IsDir() bool             { return d.info.IsDir() }
func (d dirEntry) Type() FileMode          { return d.info.Mode().Type() }
func (d dirEntry) Info() (FileInfo, error) { return d.info, nil }

type fileInfo struct {
	os.FileInfo
}

func (fi fileInfo) Mode() FileMode { return FileMode(fi.FileInfo.Mode()) }
