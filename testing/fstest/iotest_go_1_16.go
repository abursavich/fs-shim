// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.16

package fstest

import (
	"io"
	"testing/iotest"
)

// TestReader tests that reading from r returns the expected file content.
// It does reads of different sizes, until EOF.
// If r implements io.ReaderAt or io.Seeker, TestReader also checks
// that those operations behave as they should.
//
// If TestReader finds any misbehaviors, it returns an error reporting them.
// The error text may span multiple lines.
func iotestTestReader(r io.Reader, content []byte) error {
	return iotest.TestReader(r, content)
}
