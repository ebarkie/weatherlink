// Copyright (c) 2016-2017 Eric Barkie. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package weatherlink

import (
	"io/ioutil"
	"log"
)

// Setup loggers that can be overridden by the user.
//
// When using the standard logger it would be useful to use SetOutput
// to toggle between os.Std* and ioutil.Discard, For example:
//
// Info.SetOutput(ioutil.Discard)
//
// When using other loggers like Logrus SetOutput can be used to set
// the output to its io.Writer interface and SetFlags can be used to
// disable timestamps (since it has its own).  For example:
//
// weatherlink.Debug.SetOutput(log.WriterLevel(logrus.DebugLevel))
// weatherlink.Debug.SetFlags(log.Lshortfile)

// Loggers
var (
	Trace = log.New(ioutil.Discard, "[TRCE]", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	Debug = log.New(ioutil.Discard, "[DBUG]", log.LstdFlags|log.Lshortfile)
	Info  = log.New(ioutil.Discard, "[INFO]", log.LstdFlags)
	Warn  = log.New(ioutil.Discard, "[WARN]", log.LstdFlags|log.Lshortfile)
	Error = log.New(ioutil.Discard, "[ERRO]", log.LstdFlags|log.Lshortfile)
)
