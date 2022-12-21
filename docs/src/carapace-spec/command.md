# Command

```go
type Command struct {
	Name            string
	Aliases         []string
	Description     string
	Group           string
	Flags           map[string]string
	PersistentFlags map[string]string
	Completion      struct {
		Flag          map[string][]string
		Positional    [][]string
		PositionalAny []string
		Dash          [][]string
		DashAny       []string
	}
	Commands []Command
}
```
