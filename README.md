# Cripto CLI

  - [Summary](#summary)
  - [Run from sources](#run-from-sources)
  - [Features](#features)
  - [Sending files](#sending-files)
    - [Elgamal cipher](#elgamal-cipher)
    - [RSA cipher](#rsa-cipher)
    - [Shamir cipher](#shamir-cipher)
  - [Digital signature](#digital-signature)
    - [RSA based](#rsa-based)
  - [Implementation](#implementation)


## Summary

    Educational project. CLI utilities implemented with basic cryptography algorithms

## Run from sources

```shell 
go run */**.go 
```

## Features

*Each command supports -h/--help

## Sending files 

You can send files using several encryption algorithms

### Elgamal cipher

```shell
# this starts receiver server and saves public key to file
crypto elgamal-msg recv -P 30803 -G 2 -port 12346 -o file 

# this sends given file to server  
crypto elgamal-msg send -bob-pub bob_elgamal.key -P 30803 -G 2 -i file localhost:12346 some-file.txt
```

### RSA cipher

```shell
# this starts receiver server and saves public key to file
crypto rsa-msg recv -P 30803 -Q 1297 -o file -port 12346

# this sends given file to server  
crypto rsa-msg send -bob-pub bob_rsa.key -i file localhost:12346 some-file.txt
```

### Shamir cipher

```shell
# this starts receiver server
crypto rsa-msg recv -P 30803 -Q 1297 -o file -port 12346

# this sends given file to server  
crypto rsa-msg send -bob-pub bob_rsa.key -i file localhost:12346 some-file.txt
```

## Digital signature

### RSA based

```shell
# generate public and secret keys. P and Q are large prime numbers
crypto rsa-ds key-gen -P 30803 -Q 1297 

# sign file with private key
crypto rsa-ds sign -secret rsa.key some-file.txt

# validate signed file with public key
crypto rsa-ds validate -pub rsa_pub.key signed.txt
```

## Implementation

    Despite of the fact this application is not really complex, I tried to support ideas of clean architecture and explicitly manage dependencies between codebase parts.

    Wherever as possible I use interfaces and functional arguments rather than theirs implementations to separate responsibilities.

    Also such approach allows to write highly testable code.

    Project structure is pretty standart:
    `cmd` - deals with conponents of user interface such as flag         parsing, command structure specification etc. 
    `internal/core` - contains pure math and algorithms implementation
    `internal/app` - use cases and infrastrucure interactions, e.g. implementations of TCP-based file transfer protocols
