module github.com/carapace-sh/carapace-spec/cmd

go 1.24.0

replace github.com/carapace-sh/carapace-spec => ../

require (
	github.com/carapace-sh/carapace v1.17.0-alpha
	github.com/carapace-sh/carapace-selfupdate v0.0.10
	github.com/carapace-sh/carapace-spec v0.0.0-00010101000000-000000000000
	github.com/invopop/jsonschema v0.14.0
	github.com/spf13/cobra v1.10.2
	go.yaml.in/yaml/v3 v3.0.5
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.2 // indirect
	github.com/carapace-sh/carapace-shlex/v2 v2.0.0-alpha // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pb33f/ordered-map/v2 v2.3.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.yaml.in/yaml/v4 v4.0.0-rc.2 // indirect
)

replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.1.0
