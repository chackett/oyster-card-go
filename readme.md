[![Go Report Card](https://goreportcard.com/badge/github.com/chackett/oyster-test)](https://goreportcard.com/report/github.com/chackett/oyster-test)

# Tech test Oyster card

Conor Hackett - December 2019

Before running the application, please download required dependencies ```$ go mod download```

To run the application simply execute 
``$ go run cmd/main.go`` from the project (where this readme is located).

Expected output, which demonstrates test case specified in test:
```Card balance: £23.70
  Travel Summary:
  (£2.50) Holborn -> Earl's Court
  (£1.80) Earls Court Bus Stop
  (£2.00) Earl's Court -> Hammersmith
```

Go version `go version`  
``go version go1.13.4 darwin/amd64``

Notes:  
* ~270 lines of code (excluding tests)
* Reasonable test coverage to cover the more "involved" functions for sanity.
