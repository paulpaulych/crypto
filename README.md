go run */*.go elgamal-msg recv -P 30803 -G 2 -port 12346 -o file

go run */*.go elgamal-msg send -bob-pub key.txt -P 30803 -G 2 -i file localhost:12346 .gitignore