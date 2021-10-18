## elgamal messaging

crypto elgamal-msg recv -P 30803 -G 2 -port 12346 -o file crypto elgamal-msg send -bob-pub bob_elgamal.key -P 30803 -G 2
-i file localhost:12346 .gitignore

## rsa messaging

crypto rsa-msg recv -P 30803 -Q 1297 -o file -port 12346 crypto rsa-msg send -bob-pub bob_rsa.key -i file localhost:
12346 .gitignore
