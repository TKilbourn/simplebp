package simplebp

import (
	"github.com/google/blueprint"
	"github.com/google/blueprint/pathtools"

	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	scriptRule = pctx.StaticRule("script",
		blueprint.RuleParams{
			Command:     "$script $args",
			Description: "RUN $script",
		},
		"script", "args")
)

type ScriptModule struct {
	properties struct {
		Script string
		Inputs []string
		Output string
		Args   string
	}
}

func NewScript() (blueprint.Module, []interface{}) {
	module := new(ScriptModule)
	properties := &module.properties
	return module, []interface{}{properties}
}

type scriptInput struct {
	Name      string
	Basename  string
	Extension string
}

type scriptArgs struct {
	Input  scriptInput
	Output string
}

func (m *ScriptModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	config := ctx.Config().(*config)

	var scriptPath string
	if s := m.properties.Script; strings.HasPrefix(s, "//") {
		scriptPath = filepath.Join(config.srcDir, s[2:])
	} else {
		scriptPath = filepath.Join(ctx.ModuleDir(), s)
	}

	if stat, err := os.Stat(scriptPath); err != nil {
		ctx.ModuleErrorf("Could not stat %v: %v", scriptPath, err)
		return
	} else if stat.Mode()&0111 == 0 {
		ctx.ModuleErrorf("%s is not an executable", scriptPath)
		return
	}

	srcs := pathtools.PrefixPaths(m.properties.Inputs, ctx.ModuleDir())

	argsTmpl, argsErr := template.New("args").Parse(m.properties.Args)
	if argsErr != nil {
		ctx.ModuleErrorf("Could not parse script args: %v", argsErr)
		return
	}
	outTmpl, outErr := template.New("out").Parse(m.properties.Output)
	if outErr != nil {
		ctx.ModuleErrorf("Could not parse output template: %v", outErr)
		return
	}

	for _, s := range srcs {
		args := &scriptArgs{
			Input: scriptInput{
				Name:      s,
				Basename:  strings.TrimSuffix(s, filepath.Ext(s)),
				Extension: filepath.Ext(s),
			},
		}

		outBuf := &bytes.Buffer{}
		if err := outTmpl.Execute(outBuf, args); err != nil {
			ctx.ModuleErrorf("Could not generate output: %v", err)
			return
		}

		args.Output = filepath.Join(config.buildDir, outBuf.String())
		argsBuf := &bytes.Buffer{}
		if err := argsTmpl.Execute(argsBuf, args); err != nil {
			ctx.ModuleErrorf("Could not generate args: %v", err)
			return
		}

		ctx.Build(pctx, blueprint.BuildParams{
			Rule:      scriptRule,
			Inputs:    []string{s},
			Outputs:   []string{args.Output},
			Implicits: []string{scriptPath},
			Args: map[string]string{
				"script": scriptPath,
				"args":   argsBuf.String(),
			},
		})
	}
}
