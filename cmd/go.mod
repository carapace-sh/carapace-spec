module github.com/carapace-sh/carapace-spec/cmd

go 1.23.1

replace github.com/carapace-sh/carapace-spec => ../

require (
	github.com/carapace-sh/carapace v1.3.3
	github.com/carapace-sh/carapace-spec v0.0.0-00010101000000-000000000000
	github.com/invopop/jsonschema v0.12.0
	github.com/spf13/cobra v1.8.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/carapace-sh/carapace-shlex v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
)
