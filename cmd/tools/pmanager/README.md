# PManager

> Stateless password manager

```bash
usage: 
    ./main [service-name] [login]
stdin:
    [master-key]EOL
```

EOL - End of Line (Enter)

## Example

```bash
$ go run . service-name login
> master-key
62ZD5+xzT+eQkqFjNJqLryOsLSxSUzfCMEHlt6Y4dEo=
```

```bash
$ echo "master-key" | go run . service-name login
> 62ZD5+xzT+eQkqFjNJqLryOsLSxSUzfCMEHlt6Y4dEo=
```
