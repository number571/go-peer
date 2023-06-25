# Encryptor

> Encrypt messages by asymmetric keys

```bash
usage: 
    ./main [e|d] [pubkey-file|privkey-file]
stdin:
    [data-value]EOF
```

EOF - End of File (Ctrl+D)

## Example

Generate keys

```bash
make genkey # default key size = 4096 bits
# OR
make genkey N=1024 # with custom set key size
```

Encrypt/Decrypt

```bash
$ echo "hello, world" | ./main e pub.key > encrypted.msg
$ cat encrypted.msg | ./main d priv.key
```
