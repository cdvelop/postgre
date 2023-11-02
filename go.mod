module github.com/cdvelop/postgre

go 1.20

require (
	github.com/cdvelop/model v0.0.64
	github.com/cdvelop/objectdb v0.0.66
	github.com/cdvelop/timeserver v0.0.1
	github.com/cdvelop/unixid v0.0.2
	github.com/lib/pq v1.10.9
)

require (
	github.com/cdvelop/dbtools v0.0.41 // indirect
	github.com/cdvelop/gotools v0.0.43 // indirect
	github.com/cdvelop/input v0.0.34 // indirect
	github.com/cdvelop/timetools v0.0.4 // indirect
	golang.org/x/text v0.13.0 // indirect
)

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/unixid => ../unixid

replace github.com/cdvelop/timetools => ../timetools

replace github.com/cdvelop/timeserver => ../timeserver

replace github.com/cdvelop/input => ../input

replace github.com/cdvelop/objectdb => ../objectdb

replace github.com/cdvelop/dbtools => ../dbtools
