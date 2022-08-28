module github.com/rsteube/carapace-spec

go 1.18

replace github.com/spf13/pflag => github.com/cornfeedhobo/pflag v1.1.0

require (
	github.com/rsteube/carapace v0.20.2
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require github.com/inconshreveable/mousetrap v1.0.0 // indirect
