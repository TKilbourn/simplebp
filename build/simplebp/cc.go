package simplebp

import (
    "github.com/google/blueprint"
    "github.com/google/blueprint/pathtools"
)

var (
    pctx = blueprint.NewPackageContext("bp/build/simplebp")

    cc = pctx.StaticVariable("cc", "gcc")
    cFlags = pctx.StaticVariable("cFlags", "-Wall -std=c99 -O2")

    ccRule = pctx.StaticRule("cc",
            blueprint.RuleParams{
                Command: "$cc $cFlags $in -o $out",
                Description: "CC $out",
            })
)

type ccBinary struct {
    properties struct {
        Srcs []string
    }
}

func NewCcBinary() (blueprint.Module, []interface{}) {
    module := new(ccBinary)
    properties := &module.properties
    return module, []interface{}{properties}
}

func (m *ccBinary) GenerateBuildActions(ctx blueprint.ModuleContext) {
    config := ctx.Config().(*config)
    srcs := pathtools.PrefixPaths(m.properties.Srcs, ctx.ModuleDir())
    out := pathtools.PrefixPaths([]string{ctx.ModuleName()}, config.buildDir)
    ctx.Build(pctx, blueprint.BuildParams{
        Comment: "build a cc_binary",
        Rule: ccRule,
        Inputs: srcs,
        Outputs: out,
    })
}
