# HL

> Hidden Lake

<img src="_images/hl_logo.png" alt="hl_logo.png"/>

The `Hidden Lake` is an anonymous network built on a `micro-service` architecture. At the heart of HL is the core - `HLS` (service), which generates anonymizing traffic and combines many other services (for example, `HLT` and `HLM`). Thus, Hidden Lake is not a whole and monolithic solution, but a composition of several combined services.

## How it works

The scheme of the anonymous Hidden Lake network boils down to the use of three levels: `network` (N), `friendly` (F), and `application` (A).

<p align="center"><img src="_images/hl_scheme.jpg" alt="hl_scheme.jpg"/></p>
<p align="center">Figure 1. The scheme of the anonymous Hidden Lake network.</p>

1. The network layer ensures the `transfer of raw bytes` from one node to another. HLS and HLT services interact at this level. HLT is an auxiliary service that is used to relay messages from HLS. The HLT can be replaced by an HLS node.
2. The friendly layer (also known as anonymizing) performs the function of `anonymizing traffic` by setting up a list of friends and queue control on the side of the HLS service. Cryptographic routing based on public keys works at this level.
3. The application layer boils down to using the final `logic of the application` itself to transmit and receive messages. One of the main tasks of this level is to control data security, provided that intermediate (group) HLS nodes exist/are used.

## Settings and connections

```yaml
# [HLS, HLT] nodes
message_size_bytes: 8192
work_size_bits: 20
key_size_bits: 4096
queue_period_ms: 5000
limit_void_size_bytes: 4096
# 'network_key' parameter is shown in the table

# [HLT] nodes
## [62.233.46.109, 94.103.91.81]
messages_capacity: 1048576 # 2^20 msgs ~= 8GiB
## [185.43.4.253]
messages_capacity: 33554432 # 2^25 msgs ~= 256GiB
```

<table style="width: 100%">
  <tr>
    <th>Available network</th>
    <th>Types of services</th>
  </tr>
  <tr>
    <td>
        <table style="width: 100%">
            <tr>
                <th>ID</th>
                <th>Type</th>
                <th>Version</th>
                <th>Host</th>
                <th>Port</th>
                <th>Network key</th>
                <th>Connections</th>
                <th>Provider</th>
                <th>Country</th>
                <th>City</th>
                <th>Characteristics</th>
                <th>Expired time</th>
            </tr>
            <tr>
                <td>1</td>
                <td>HLTr/HLTs</td>
                <td>v1.5.24</td>
                <td>94.103.91.81</td> 
                <td>9581/9582</td>
                <td>8Jkl93Mdk93md1bz</td>
                <td>[]</td>
                <td><a href="https://vdsina.ru">vdsina.ru</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x4.0GHz, 1.0GB RAM, 30GB HDD</td>
                <td>±eternal</td>
            </tr>
            <tr>
                <td>2</td>
                <td>HLTr</td>
                <td>v1.5.24</td>
                <td>195.133.1.126</td>
                <td>9581</td>
                <td>8Jkl93Mdk93md1bz</td>
                <td>[1]</td>
                <td><a href="https://ruvds.com">ruvds.ru</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x2.2GHz, 0.5GB RAM, 10GB HDD</td>
                <td>±28.07.2027</td>
            </tr>
            <tr>
                <td>3</td>
                <td>HLTr/HLTs</td>
                <td>v1.5.24</td>
                <td>62.233.46.109</td>
                <td>9581/9582</td>
                <td>oi4r9NW9Le7fKF9d</td>
                <td>[]</td>
                <td><a href="https://eternalhost.net">eternalhost.net</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x2.8GHz, 1.0GB RAM, 16GB HDD</td>
                <td>±eternal</td>
            </tr>
            <tr>
                <td>4</td>
                <td>HLTr</td>
                <td>v1.5.24</td>
                <td>193.233.18.245</td>
                <td>9581</td>
                <td>oi4r9NW9Le7fKF9d</td>
                <td>[3]</td>
                <td><a href="https://4vps.su">4vps.su</a></td>
                <td>Russia</td>
                <td>Novosibirsk</td>
                <td>1x2.5GHz, 1.0GB RAM, 5GB VNMe</td>
                <td>±07.08.2027</td>
            </tr>
            <tr>
                <td>5</td>
                <td>HLTs</td>
                <td>v1.5.24</td>
                <td>185.43.4.253</td>
                <td>9582</td>
                <td>j2BR39JfDf7Bajx3</td>
                <td>[]</td>
                <td><a href="https://firstvds.ru">firstvds.ru</a></td>
                <td>Russia</td>
                <td>Moscow</td>
                <td>1x3.1GHz, 2.0GB RAM, 300GB HDD</td>
                <td>±10.12.2024</td>
            </tr>
        </table>
    </td>
    <td>
        <table style="width: 100%">
            <tr>
                <th>Type</th>
                <th>Name</th>
                <th>Default port</th>
            </tr>
            <tr>
                <td>HLS</td>
                <td>node</td>
                <td>9571</td>
            </tr>
            <tr>
                <td>HLTr</td>
                <td>relayer</td>
                <td>9581</td>
            </tr>
            <tr>
                <td>HLTs</td>
                <td>storage</td>
                <td>9582</td>
            </tr>
        </table>
    </td>
  </tr>
</table>

## Build and run

Launching an anonymous network is primarily the launch of an anonymizing HLS service. There are three ways to run HLS: through `docker`, through `source code`, and through the `release version`. It is recommended to run applications with the available release version, [tag](https://github.com/number571/go-peer/tags).

### 1. Running from docker

```bash
$ git clone -b <tag-name> --depth=1 https://github.com/number571/go-peer.git
$ cd go-peer/cmd/hidden_lake/service
$ make docker-default
...
> [INFO] 2023/12/29 16:27:40 HLS is running...
> [INFO] 2023/12/29 16:27:45 service=HLS type=BRDCS hash=FF59D8F9...1D2EAF8D addr=23FDFF8A...FE95F5E0 proof=0000160292 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 16:27:50 service=HLS type=BRDCS hash=ACE23689...CA39CB6D addr=23FDFF8A...FE95F5E0 proof=0001491994 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 16:27:55 service=HLS type=BRDCS hash=0C6EC76D...FB83729D addr=23FDFF8A...FE95F5E0 proof=0001762328 size=8192B conn=127.0.0.1:
...
```

### 2. Running from source code

```bash
$ git clone -b <tag-name> --depth=1 https://github.com/number571/go-peer.git
$ cd go-peer/cmd/hidden_lake/service
$ make default
...
> [INFO] 2023/12/29 23:29:43 HLS is running...
> [INFO] 2023/12/29 23:29:48 service=HLS type=BRDCS hash=8B77F546...3CE1421C addr=E04D2DC8...61D4FE2A proof=0001379020 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 23:29:53 service=HLS type=BRDCS hash=3EA3189F...DC793A4E addr=E04D2DC8...61D4FE2A proof=0000076242 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 23:29:58 service=HLS type=BRDCS hash=475A8886...0F77621F addr=E04D2DC8...61D4FE2A proof=0001964664 size=8192B conn=127.0.0.1:
...
```

### 3. Running from release version

When starting from the release version, you must specify the processor architecture and platform used. Available architectures: `amd64`, `arm64`. Available platforms: `windows`, `darwin`, `linux`.

```bash
$ wget https://github.com/number571/go-peer/releases/download/<tag-name>/hls_<arch-name>_<platform-name>
$ chmod +x hls_<arch-name>_<platform-name>
$ ./hls_<arch-name>_<platform-name>
...
> [INFO] 2023/12/29 23:31:43 HLS is running...
> [INFO] 2023/12/29 23:31:48 service=HLS type=BRDCS hash=E8CDB448...FF23639E addr=E04D2DC8...61D4FE2A proof=0001277744 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 23:31:53 service=HLS type=BRDCS hash=C6B5C47F...AB63128A addr=E04D2DC8...61D4FE2A proof=0001062655 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 23:31:58 service=HLS type=BRDCS hash=5789D462...B81C3A5F addr=E04D2DC8...61D4FE2A proof=0000517841 size=8192B conn=127.0.0.1
...
```

### Production running

The HLS node is easy to connect to a production environment. To do this, it is sufficient to specify two parameters: `network_key` and `connections`. The network_key parameter is used to separate networks from each other, preventing them from merging. The connections parameter is used for direct network connection to HLS and HLT nodes.

```bash
$ cd cmd/hidden_lake/service

# [From docker]: $ yes | cp ../_configs/prod/1/hls.yml ./_mounted
$ yes | cp ../_configs/prod/1/hls.yml . # [From source | release]

# [From docker]:  $ make docker-run
# [From source]:  $ make run
$ ./hls_<arch-name>_<platform-name> # [From release]
...
> [INFO] 2023/12/29 23:49:26 HLS is running...
> [INFO] 2023/12/29 23:49:27 service=HLS type=UNDEC hash=1079200E...FFCD5871 addr=00000000...00000000 proof=0000165513 size=8192B conn=94.103.91.81:9581
> [INFO] 2023/12/29 23:49:31 service=HLS type=BRDCS hash=1DE4BC0F...AC611F44 addr=E04D2DC8...61D4FE2A proof=0000265462 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 23:49:32 service=HLS type=UNDEC hash=EECEF795...1B042618 addr=00000000...00000000 proof=0002571939 size=8192B conn=94.103.91.81:9581
> [INFO] 2023/12/29 23:49:36 service=HLS type=BRDCS hash=98AC78E7...AAA7F8F1 addr=E04D2DC8...61D4FE2A proof=0001741261 size=8192B conn=127.0.0.1:
> [INFO] 2023/12/29 23:49:37 service=HLS type=UNDEC hash=7CB7B53A...1FE35530 addr=00000000...00000000 proof=0000199886 size=8192B conn=94.103.91.81:9581
> [INFO] 2023/12/29 23:49:41 service=HLS type=BRDCS hash=5D609534...9CC17DAE addr=E04D2DC8...61D4FE2A proof=0001091209 size=8192B conn=127.0.0.1:
...
```

> There are also examples of running HL applications in a production environment. For more information, follow the links: [echo_service](https://github.com/number571/go-peer/tree/master/examples/echo_service/prod_test), [anon_messenger](https://github.com/number571/go-peer/tree/master/examples/anon_messenger/prod_test).
