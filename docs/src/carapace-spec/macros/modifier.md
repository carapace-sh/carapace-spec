# Modifier

Modifiers change the completion for a position in general.

## chdir

`$chdir(<directory>)` changes the directory.

```yml
["$chdir(/tmp)", "$(pwd)"]
```

## list

`$list(<delimiter>)` completes values as unique list with given delimiter.

```yml
["$list(,)", "a", "b", "c", "d"]
```

## multiparts

`$multiparts(<delimiter>)` completes values splitted on given delimiter separately.

```yml
["$multiparts(/)", "a", "a/b", "a/c", "b", "b/a"]
```

## nospace

`$nospace` prevents space suffix being added to the inserted values.

```yml
["$nospace", "one", "two"]
```
