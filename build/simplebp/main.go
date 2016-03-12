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

    config := simplebp.NewConfig(srcDir, bootstrap.BuildDir)

    bootstrap.Main(ctx, config)
}
