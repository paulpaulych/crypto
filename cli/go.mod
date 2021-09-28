module github.com/paulpaulych/crypto/shamir/cli

go 1.16

replace (
	github.com/paulpaulych/crypto/shamir => ./../shamir
	github.com/paulpaulych/crypto/commons => ./../commons
)

require github.com/paulpaulych/crypto/shamir v0.0.0-00010101000000-000000000000
