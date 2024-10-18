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
J64mESAVm2o-8q6nWaywsoQbQvn8pf7U74O-Vr9HwSDu
```

```bash
$ echo "master-key" | go run . -salt="service-name" -work=24
J64mESAVm2o-8q6nWaywsoQbQvn8pf7U74O-Vr9HwSDu
```
