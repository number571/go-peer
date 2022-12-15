# HLM

> Hidden Lake Messenger

## Description

The HLM is a messenger based on the core of an anonymous network with theoretically provable anonymity of HLS. A feature of this messenger is the provision of anonymity of the fact of transactions (sending, receiving). 

## Cryptographic algorithms and functions

1. AES-256-CTR (Data encryption)
2. RSA-4096-OAEP (Key encryption)
3. RSA-4096-PSS (Hash signing)
4. SHA-256 (Data hashing)
5. HMAC-SHA-256 (Network hashing)
6. PoW-20 (Hash proof)

## Signup page

Sign up login/password and additional private key. If field with private key is null than private key generated.

<img src="examples/images/v2/signup.png" alt="signup.png"/>

## Signin page

Sign in with login/password. Authorized client identified by a private key and can push messages into database.

<img src="examples/images/v2/signin.png" alt="signin.png"/>

## About page

Base information about projects HLM and HLS with links to source.

<img src="examples/images/v2/about.png" alt="about.png"/>

## Settings page

Information about public key and connections. Connections can be appended and deleted.

<img src="examples/images/v2/settings.png" alt="settings.png"/>

## Friends page

Information about friends. Friends can be appended and deleted.

<img src="examples/images/v2/friends.png" alt="friends.png"/>

## Chat page

Chat with friend. The chat is based on web sockets, so it can update messages in real time. Messages can be sent.

<img src="examples/images/v2/chat.png" alt="chat.png"/>
