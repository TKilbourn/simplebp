package simplebp

type config struct {
    srcDir string
    buildDir string
}

func NewConfig(srcDir string, buildDir string) interface{} {
    return &config{srcDir, buildDir}
}
