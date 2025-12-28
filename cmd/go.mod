module github.com/carapace-sh/carapace-spec/cmd

go 1.24

replace github.com/carapace-sh/carapace-spec => ../

require (
	github.com/carapace-sh/carapace v1.11.0
	github.com/carapace-sh/carapace-selfupdate v0.0.5
	github.com/carapace-sh/carapace-spec v0.0.0-00010101000000-000000000000
	github.com/invopop/jsonschema v0.12.0
	github.com/spf13/cobra v1.10.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/carapace-sh/carapace-shlex v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
)

replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.1.0
