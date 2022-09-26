module github.com/rsteube/carapace-spec

go 1.18

require (
	github.com/rsteube/carapace v0.24.1
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require github.com/inconshreveable/mousetrap v1.0.0 // indirect

replace github.com/spf13/pflag => github.com/rsteube/carapace-pflag v0.0.4
