## Secpy-Chat

> The Secpy-Chat application

<img src="_images/secpy_chat_logo.png" alt="secpy_chat_logo.png"/>

The application `secpy_chat` allows you to communicate securely (using end-to-end encryption) using HLT and HLE applications. This is an example of how it is possible to write client-safe applications for the Hidden Lake environment without being based on the Go programming language (the main language for writing Hidden Lake applications).

## Config structure

```
"hlt_host" address of the HLT service
"hle_host" address of the HLE service
"friends"  map of {"alias_name":"public_key"}
```

```yaml
hlt_host: localhost:9582
hle_host: localhost:9551
friends: 
  Alice: PubKey{3082020A02820201...3324D10203010001}
```

## How it works

<p align="center"><img src="_images/secpy_chat.gif" alt="secpy_chat.gif"/></p>
<p align="center">Figure 1. Chat node1 with node2.</p>

The application connects to two services at once: [HLE](https://github.com/number571/go-peer/tree/master/cmd/hidden_lake/helpers/encryptor) and [HLT](https://github.com/number571/go-peer/tree/master/cmd/hidden_lake/traffic). The first service makes it possible to encrypt and decrypt messages. The second service allows you to send and receive encrypted messages from the network. In this case, the secpy_chat is guided only by the interfaces of the services, representing the frontend component.

## Example 

Build and run services HLT, HLE
```bash
$ cd examples/secure_messenger
$ make
```

Run client#1
```bash
$ cd examples/secure_messenger/node1
$ python3 main.py
> /friend Alice
# waiting client#2
> hello
> [Bob]: world!
```

Run client#2
```bash
$ cd examples/secure_messenger/node2
$ python3 main.py
> /friend Bob
# waiting client#2
> [Alice]: hello
> world!
```
