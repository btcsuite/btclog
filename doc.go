// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

/*
Package btclog defines an interface and default implementation for subsystem
logging.

Log level verbosity may be modified at runtime for each individual subsystem
logger.

The default implementation in this package must be created by the Backend type.
Backends can write to any io.Writer, including multi-writers created by
io.MultiWriter.  Multi-writers allow log output to be written to many writers,
including standard output and log files.
*/
package btclog
