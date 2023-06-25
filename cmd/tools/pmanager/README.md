# PManager

> Stateless password manager

```bash
usage: 
    ./main [service-name]
stdin:
    [master-key]EOF
```

EOF - End of File (Ctrl+D)

## Example

```bash
$ echo "master-key" | go run . service-name login
fsJ5+QUz5nv/JK3VdqYWqzqQoMy0pg7FqQQtQKq2cnw=
```
