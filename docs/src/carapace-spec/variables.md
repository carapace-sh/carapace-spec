# Variables

Variables are replaced using [drone/envsubst](https://github.com/drone/envsubst).

- `${C_ARG<position>}` positional arguments `[0..n]`
- `${C_FLAG_<flagname>}` flag values (if modified)
- `${C_PART<position>}` parts of the current word during multipart completion `[0..n]`
- `${C_VALUE}` the word currently being completed

```yaml
name: myvar
flags:
  --suffix=: file suffixes
completion:
  flag:
    suffix: ["$list(,)", ".go", "go.sum", "go.mod", ".md", "LICENSE"]
  positional:
    - ["$files([${C_FLAG_SUFFIX//,/, }])"] # replace `,` with `, ` for valid array syntax
    - ["${C_FLAG_SUFFIX:-default}", "${C_ARG0}"] # use default if flag is not set
```

![](./variables.cast)
