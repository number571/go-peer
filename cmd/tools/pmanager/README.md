# PManager

> Stateless password manager

```bash
usage: 
    ./main -salt=[service-name] -work=[diff-size]
stdin:
    [master-key]EOL
```

EOL - End of Line (Enter)

## Example

```bash
$ go run . -salt="service-name" -work=24
master-key
EvCqIyPVq9ydjspox6GqN63ggT0xrUfNFnFgwAy1odQ=
```

```bash
$ echo "master-key" | go run . -salt="service-name" -work=24
EvCqIyPVq9ydjspox6GqN63ggT0xrUfNFnFgwAy1odQ=
```
