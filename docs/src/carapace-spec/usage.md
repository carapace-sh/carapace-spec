# Usage

## Bash
```sh
# ~/.bashrc
source <(carapace-spec example/pkill.yaml)
```
![](./usage/bash.png)

## Elvish
```sh
# ~/.config/elvish/rc.elv
eval (carapace-spec example/pkill.yaml|slurp)
```
![](./usage/elvish.png)

## Fish
```sh
# ~/.config/fish/config.fish
carapace-spec example/pkill.yaml | source
```
![](./usage/fish.png)

## Nushell
> update config.nu according to [Multiple Completer](http://www.nushell.sh/cookbook/external_completers.html#multiple-completer))
```sh
#~/.config/nushell/config.nu
carapace-spec example/pkill.yaml
```
![](./usage/nushell.png)

## Oil
```sh
# ~/.config/oil/oshrc
source <(carapace-spec example/pkill.yaml)
```
![](./usage/oil.png)

## Powershell
```sh
# ~/.config/powershell/Microsoft.PowerShell_profile.ps1
carapace-spec example/pkill.yaml | Out-String | Invoke-Expression
```
![](./usage/powershell.png)

# Tcsh
```sh
## ~/.tcshrc
eval `carapace-spec example/pkill.yaml`
```
![](./usage/tcsh.png)

## Xonsh
```sh
# ~/.config/xonsh/rc.xsh
exec($(carapace-spec example/pkill.yaml))
```
![](./usage/xonsh.png)

## Zsh
```sh
# ~/.zshrc
source <(carapace-spec example/pkill.yaml)
```
![](./usage/zsh.png)
