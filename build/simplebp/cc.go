package simplebp

import (
    "path/filepath"
    "github.com/google/blueprint"
    "github.com/google/blueprint/pathtools"
)

var (
    pctx = blueprint.NewPackageContext("bp/build/simplebp")

    cc = pctx.StaticVariable("cc", "gcc")
    ld = pctx.StaticVariable("ld", "gcc")
    cFlags = pctx.StaticVariable("cFlags", "-Wall -std=c99 -O2")
    ldFlags = pctx.StaticVariable("ldFlags", "")
    libs = pctx.StaticVariable("libs", "")

    ccRule = pctx.StaticRule("cc",
            blueprint.RuleParams{
                Command: "$cc -MMD -MF $out.d $cFlags -c $in -o $out",
                Depfile: "$out.d",
                Deps: blueprint.DepsGCC,
                Description: "CC $out",
            })

    linkRule = pctx.StaticRule("link",
            blueprint.RuleParams{
                Command: "$ld $ldFlags $in -o $out $libs",
                Description: "LINK $out",
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

    objs := make([]string, 0, len(srcs))

    for _, s := range srcs {
        obj := filepath.Join(config.buildDir, pathtools.ReplaceExtension(s, "o"))
        ctx.Build(pctx, blueprint.BuildParams{
            Comment: "compile a C source file",
            Rule: ccRule,
            Inputs: []string{s},
            Outputs: []string{obj},
        })
        objs = append(objs, obj)
    }

    out := pathtools.PrefixPaths([]string{ctx.ModuleName()},
            filepath.Join(config.buildDir, ctx.ModuleDir()))

    ctx.Build(pctx, blueprint.BuildParams{
        Comment: "build a C binary",
        Rule: linkRule,
        Inputs: objs,
        Outputs: out,
    })
}
