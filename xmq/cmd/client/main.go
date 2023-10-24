package main

import (
	"github.com/learn-go/xmq/pkg"
	_ "github.com/learn-go/xmq/pkg/client"
)

func main() {
	pkg.PkgClient.Start()
}
