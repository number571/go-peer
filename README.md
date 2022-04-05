# go-peer

> Framework for create secure and anonymity decentralized networks. Version: 1.4

## Research Article
* Theory of the structure of hidden systems: [hidden_systems.pdf](https://github.com/Number571/go-peer/blob/master/hidden_systems.pdf "TSHS");

## Framework based applications
* Hidden Lake Service: [github.com/number571/go-peer/tree/master/cmd/hls](https://github.com/number571/go-peer/tree/master/cmd/hls "HLS");
* Hidden Message Service: [github.com/number571/go-peer/tree/master/cmd/hms](https://github.com/number571/go-peer/tree/master/cmd/hms "HMS");

## Deprecated framework based applications
* Hidden Lake: [github.com/number571/hidden-lake](https://github.com/number571/hidden-lake "HL");
* Hidden Email Service: [github.com/number571/hes](https://github.com/number571/hes "HES");

## Need TODO

At the moment, the framework is able to recreate the five stage of anonymity, but is not suitable for the seven. This is due to the two pitfalls of the seven stage of anonymity that need to be corrected. The list is as follows:
1. Request time. You need to implement a simulation of packet generation time, either on a request-based or routing-response basis.
2. The period of states. This problem should be solved dynamically by the user and based on the framework it is quite impossible to fix the vulnerability.

## Specifications of go-peer

1. Prefix 's'/'S' - structure type
2. Prefix 'i'/'I' - interface type
3. Prefix 'f'/'F' - field of structure
4. Prefix 'c'/'C' - constant
5. Prrfix 'g'/'G' - global variable
6. Prefix 't'/'T' - test constant/variable/structure
