1. Install (golang)[https://golang.org/]
1. Install `bash` and `gnu-make`.
1. Run  `scripts/devtools.sh` to install the required devtools
1. Run `make [local] vendor` to install the vendor dependencies
1. Run `make [local] test` to ensure you can run the existing tests
1. Write unit tests - (see existing examples). Changes will not be accepted without corresponding unit tests.
1. Make the code changes.
1. Run `make [local] test` to lint code and run tests
1. Profit! ok no profit, but raise a PR and get kudos :)

...or if you would rather not install `bash` or `make`

1. Install (golang)[https://golang.org/]
1. Run `go run mage.go EnsureMage`
1. Run `mage devtools` to install the required devtools
1. Run `mage [local] vendor` to install the vendor dependencies
1. Run `mage [local] test` to ensure you can run the existing tests
1. Write unit tests - (see existing examples). Changes will not be accepted without corresponding unit tests.
1. Make the code changes.
1. Run `mage [local] test` to lint code and run tests
1. Profit! ok no profit, but raise a PR and get kudos :)
