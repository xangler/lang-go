package main

import (
	"github.com/learn-go/xmq/pkg"
	_ "github.com/learn-go/xmq/pkg/server"
)

func main() {
	pkg.PkgServer.Start()
}
