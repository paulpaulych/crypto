## elgamal messaging

```shell
# run bob
crypto elgamal-msg recv -P 30803 -G 2 -port 12346 -o file

# run alice
crypto elgamal-msg send -bob-pub bob_elgamal.key -P 30803 -G 2 -i file localhost:12346 .gitignore
```

## rsa messaging

```shell
# run bob
crypto rsa-msg recv -P 30803 -Q 1297 -o file -port 12346

# run alice
crypto rsa-msg send -bob-pub bob_rsa.key -i file localhost:12346 .gitignore
```

## rsa digital signature

```shell
# generate public and secret keys. P and Q are large prime numbers
crypto rsa-ds key-gen -P 30803 -Q 1297 

# sign message with private key
crypto rsa-ds sign -secret rsa.key 'some message'

# validate signed message with public key
crypto rsa-ds validate -pub rsa_pub.key signed.txt
```