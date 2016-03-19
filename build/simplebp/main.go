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

	ctx.RegisterModuleType("c_binary", simplebp.NewCBinary)
	ctx.RegisterModuleType("c_shared_lib", simplebp.NewCSharedLib)
	ctx.RegisterModuleType("run_script", simplebp.NewScript)

	config := simplebp.NewConfig(srcDir, bootstrap.BuildDir)

	bootstrap.Main(ctx, config)
}
