module github.com/rsteube/carapace-spec

go 1.18

require (
	github.com/rsteube/carapace v0.24.1
	github.com/spf13/cobra v1.6.0
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/inconshreveable/mousetrap v1.0.1 // indirect

replace github.com/spf13/pflag => github.com/rsteube/carapace-pflag v0.0.4
