# Usage

```sh
# bash
source <(carapace-spec example.yaml)

# elvish
eval (carapace-spec example.yaml|slurp)

# fish
carapace-spec example.yaml | source

# oil
source <(carapace-spec example.yaml)

# nushell (update config.nu according to output)
carapace-spec example.yaml

# powershell
carapace-spec example.yaml | Out-String | Invoke-Expression

# tcsh
eval `carapace-spec example.yaml`

# xonsh
exec($(carapace-spec example.yaml))

# zsh
source <(carapace-spec example.yaml)
```
