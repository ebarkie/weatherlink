// Copyright (c) 2016 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"fmt"
	"io"
	"log"
	"strings"
)

// Setup loggers that can be overridden by the user.
//
// When using the standard logger it would be useful to use SetOutput
// to toggle between os.Std* and io.Discard, For example:
//
// Info.SetOutput(io.Discard)
//
// When using other loggers like Logrus SetOutput can be used to set
// the output to its io.Writer interface and SetFlags can be used to
// disable timestamps (since it has its own).  For example:
//
// weatherlink.Debug.SetOutput(log.WriterLevel(logrus.DebugLevel))
// weatherlink.Debug.SetFlags(log.Lshortfile)

// Loggers.
var (
	Trace = log.New(io.Discard, "[TRCE]", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	Debug = log.New(io.Discard, "[DBUG]", log.LstdFlags|log.Lshortfile)
	Info  = log.New(io.Discard, "[INFO]", log.LstdFlags)
	Warn  = log.New(io.Discard, "[WARN]", log.LstdFlags|log.Lshortfile)
	Error = log.New(io.Discard, "[ERRO]", log.LstdFlags|log.Lshortfile)
)

// Sdump returns a variable as a string.  It includes field names and pointers,
// if applicable.
var Sdump = func(i ...interface{}) (s string) {
	return fmt.Sprintf(strings.Repeat("%+v\n", len(i)), i...)
}
