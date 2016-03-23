package main

import (
	"flag"
	"path/filepath"

	"github.com/google/blueprint"
	"github.com/google/blueprint/bootstrap"

	"github.com/TKilbourn/simplebp"
)

func main() {
	flag.Parse()

	srcDir := filepath.Dir(flag.Arg(0))

	ctx := blueprint.NewContext()

	ctx.RegisterModuleType("cc_binary", simplebp.NewCcBinary)
	ctx.RegisterModuleType("cc_shared_lib", simplebp.NewCcSharedLib)
	ctx.RegisterModuleType("run_script", simplebp.NewScript)

	config := simplebp.NewConfig(srcDir, bootstrap.BuildDir)

	bootstrap.Main(ctx, config)
}
