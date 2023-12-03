module github.com/cdvelop/postgre

go 1.20

require (
	github.com/cdvelop/model v0.0.75
	github.com/cdvelop/objectdb v0.0.88
	github.com/cdvelop/timeserver v0.0.23
	github.com/cdvelop/unixid v0.0.24
	github.com/lib/pq v1.10.9
)

require (
	github.com/cdvelop/dbtools v0.0.66 // indirect
	github.com/cdvelop/input v0.0.59 // indirect
	github.com/cdvelop/maps v0.0.7 // indirect
	github.com/cdvelop/strings v0.0.7 // indirect
	github.com/cdvelop/timetools v0.0.24 // indirect
)

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/maps => ../maps

replace github.com/cdvelop/unixid => ../unixid

replace github.com/cdvelop/timetools => ../timetools

replace github.com/cdvelop/timeserver => ../timeserver

replace github.com/cdvelop/input => ../input

replace github.com/cdvelop/objectdb => ../objectdb

replace github.com/cdvelop/dbtools => ../dbtools
