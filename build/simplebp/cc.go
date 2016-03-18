package simplebp

import (
	"github.com/google/blueprint"
	"github.com/google/blueprint/pathtools"
	"path/filepath"
	"strings"
)

var (
	defaultCFlags = []string{
		"-Wall",
		"-std=c99",
		"-O2",
	}
	defaultLdFlags = []string{}

	pctx   = blueprint.NewPackageContext("bp/build/simplebp")
	ccRule = pctx.StaticRule("cc",
		blueprint.RuleParams{
			Command:     "$cc -MMD -MF $out.d $cFlags $incPaths -c $in -o $out",
			Depfile:     "$out.d",
			Deps:        blueprint.DepsGCC,
			Description: "CC $out",
		},
		"cFlags", "incPaths")

	linkRule = pctx.StaticRule("link",
		blueprint.RuleParams{
			Command:     "$ld $ldFlags $in -o $out $ldPaths $libs",
			Description: "LINK $out",
		},
		"ldFlags", "ldPaths", "libs")
)

func init() {
	pctx.StaticVariable("cc", "gcc")
	pctx.StaticVariable("ld", "gcc")

	pctx.StaticVariable("defaultCFlags", strings.Join(defaultCFlags, " "))
	pctx.StaticVariable("defaultLdFlags", strings.Join(defaultLdFlags, " "))
}

type BaseProperties struct {
	Srcs    []string // The source inputs
	Cflags  []string // The C flags to use while compiling
	Ldflags []string // The linker flags
}

type BinaryModule struct {
	properties BaseProperties // Base properties shared by all modules
	output     string         // The output artifact for the module
}

type SharedLibProperties struct {
	BaseProperties
	IncludePaths []string // Paths exported to dependers for include files
}

type SharedLibModule struct {
	properties SharedLibProperties
        incPaths   []string
	output     string
	outPath    string // The path exported to dependers for linking
}

func NewCBinary() (blueprint.Module, []interface{}) {
	module := new(BinaryModule)
	properties := &module.properties
	return module, []interface{}{properties}
}

func (m *BinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	config := ctx.Config().(*config)

	srcs := pathtools.PrefixPaths(m.properties.Srcs, ctx.ModuleDir())
	m.output = filepath.Join(config.buildDir, ctx.ModuleDir(), ctx.ModuleName())

	cflags := []string{"${defaultCFlags}"}
	cflags = append(cflags, m.properties.Cflags...)

        deps := new(depsData)
	ctx.VisitDepsDepthFirst(func(module blueprint.Module) {
		gatherDepData(module, ctx, deps)
	})

	objs := compileSrcsToObjs(ctx, srcs, cflags, deps.includePaths, config.buildDir)

        ldflags := []string{"${defaultLdFlags}"}
        ldflags = append(ldflags, m.properties.Ldflags...)

	compileObjsToOutput(ctx, objs, ldflags, deps.linkPaths, deps.libraryNames, deps.outputPaths, []string{m.output})

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:      blueprint.Phony,
		Outputs:   []string{ctx.ModuleName()},
		Implicits: []string{m.output},
	})
}

func NewCSharedLib() (blueprint.Module, []interface{}) {
	module := new(SharedLibModule)
	properties := &module.properties
	return module, []interface{}{properties}
}

func (m *SharedLibModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	config := ctx.Config().(*config)

	srcs := pathtools.PrefixPaths(m.properties.Srcs, ctx.ModuleDir())
        m.incPaths = pathtools.PrefixPaths(m.properties.IncludePaths, ctx.ModuleDir())
	m.outPath = filepath.Join(config.buildDir, ctx.ModuleDir())
        m.output = filepath.Join(m.outPath, "lib"+ctx.ModuleName()+".so")

	cflags := []string{"${defaultCFlags}", "-fPIC"}
	cflags = append(cflags, m.properties.Cflags...)

	deps := &depsData{}
	ctx.VisitDepsDepthFirst(func(module blueprint.Module) {
		gatherDepData(module, ctx, deps)
	})

	objs := compileSrcsToObjs(ctx, srcs, cflags, deps.includePaths, config.buildDir)


        ldflags := []string{"${defaultLdFlags}", "-shared"}
        ldflags = append(ldflags, m.properties.Ldflags...)

	compileObjsToOutput(ctx, objs, ldflags, deps.linkPaths, deps.libraryNames, deps.outputPaths, []string{m.output})

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:      blueprint.Phony,
		Outputs:   []string{"lib"+ctx.ModuleName()+".so"},
		Implicits: []string{m.output},
	})
}

type depsData struct {
	includePaths []string
	linkPaths    []string
	libraryNames []string
	outputPaths  []string
}

func gatherDepData(module blueprint.Module, ctx blueprint.ModuleContext, deps *depsData) {
	libModule, ok := module.(*SharedLibModule)
	if !ok {
		// TODO: report an error
		return
	}
	deps.includePaths = append(deps.includePaths, libModule.incPaths...)
	deps.linkPaths = append(deps.linkPaths, libModule.outPath)
	deps.libraryNames = append(deps.libraryNames, ctx.OtherModuleName(module))
	deps.outputPaths = append(deps.outputPaths, libModule.output)
}

func compileSrcsToObjs(ctx blueprint.ModuleContext, srcs []string, flags []string, includePaths []string, buildDir string) []string {
	flagStr := strings.Join(flags, " ")

        incPathFlags := make([]string, len(includePaths))
        for i, path := range includePaths {
            incPathFlags[i] = "-I"+path
        }
	incStr := strings.Join(incPathFlags, " ")

	objs := make([]string, len(srcs))
	for i, s := range srcs {
		objs[i] = filepath.Join(buildDir, pathtools.ReplaceExtension(s, "o"))
		ctx.Build(pctx, blueprint.BuildParams{
			Rule:    ccRule,
			Inputs:  []string{s},
			Outputs: []string{objs[i]},
			Args: map[string]string{
				"cFlags":   flagStr,
				"incPaths": incStr,
			},
		})
	}
	return objs
}

func compileObjsToOutput(ctx blueprint.ModuleContext, objs []string, flags []string, linkPaths []string, libNames []string, libOutputs []string, out []string) {
        flagStr := strings.Join(flags, " ")

        linkPathFlags := make([]string, len(linkPaths))
        for i, path := range linkPaths {
            linkPathFlags[i] = "-L"+path
        }
        linkPathStr := strings.Join(linkPathFlags, " ")

        libNameFlags := make([]string, len(libNames))
        for i, name := range libNames {
            libNameFlags[i] = "-l"+name
        }
        libNameStr := strings.Join(libNameFlags, " ")

	ctx.Build(pctx, blueprint.BuildParams{
		Rule:      linkRule,
		Inputs:    objs,
		Outputs:   out,
		Implicits: libOutputs,
		Args: map[string]string{
			"ldFlags": flagStr,
                        "ldPaths": linkPathStr,
			"libs":    libNameStr,
		},
	})
}
