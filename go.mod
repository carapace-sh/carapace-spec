module github.com/rsteube/carapace-spec

go 1.18

require (
	github.com/invopop/jsonschema v0.6.0
	github.com/rsteube/carapace v0.25.1
	github.com/spf13/cobra v1.6.0
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
)

replace github.com/spf13/pflag => github.com/rsteube/carapace-pflag v0.0.4
