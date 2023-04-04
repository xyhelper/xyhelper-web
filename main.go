package main

import (
	_ "github.com/xyhelper/xyhelper-web/internal/packed"

	_ "github.com/cool-team-official/cool-admin-go/contrib/drivers/sqlite"

	_ "github.com/xyhelper/xyhelper-web/modules"

	"github.com/gogf/gf/v2/os/gctx"

	"github.com/xyhelper/xyhelper-web/internal/cmd"
)

func main() {
	// gres.Dump()
	cmd.Main.Run(gctx.New())
}
