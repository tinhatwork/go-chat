// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/tinhatwork/go-chat/transport"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	server := transport.NewServer()
	server.ListenAndServe(*addr)
}
