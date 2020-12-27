// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.16

package fstest

import (
	"io/fs"
	"testing/fstest"
)

// mapfs.go

// A MapFS is a simple in-memory file system for use in tests,
// represented as a map from path names (arguments to Open)
// to information about the files or directories they represent.
//
// The map need not include parent directories for files contained
// in the map; those will be synthesized if needed.
// But a directory can still be included by setting the MapFile.Mode's ModeDir bit;
// this may be necessary for detailed control over the directory's FileInfo
// or to create an empty directory.
//
// File system operations read directly from the map,
// so that the file system can be changed by editing the map as needed.
// An implication is that file system operations must not run concurrently
// with changes to the map, which would be a race.
// Another implication is that opening or reading a directory requires
// iterating over the entire map, so a MapFS should typically be used with not more
// than a few hundred entries or directory reads.
type MapFS = fstest.MapFS

// A MapFile describes a single file in a MapFS.
type MapFile = fstest.MapFile

// testfs.go

// TestFS tests a file system implementation.
// It walks the entire tree of files in fsys,
// opening and checking that each file behaves correctly.
// It also checks that the file system contains at least the expected files.
// As a special case, if no expected files are listed, fsys must be empty.
// Otherwise, fsys must only contain at least the listed files: it can also contain others.
//
// If TestFS finds any misbehaviors, it returns an error reporting all of them.
// The error text spans multiple lines, one per detected misbehavior.
//
// Typical usage inside a test is:
//
//	if err := fstest.TestFS(myFS, "file/that/should/be/present"); err != nil {
//		t.Fatal(err)
//	}
//
func TestFS(fsys fs.FS, expected ...string) error {
	return fstest.TestFS(fsys, expected...)
}
