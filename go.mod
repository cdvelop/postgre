module github.com/cdvelop/postgre

go 1.20

require (
	github.com/cdvelop/model v0.0.86
	github.com/cdvelop/objectdb v0.0.92
	github.com/cdvelop/timeserver v0.0.27
	github.com/cdvelop/unixid v0.0.28
	github.com/lib/pq v1.10.9
)

require (
	github.com/cdvelop/dbtools v0.0.68 // indirect
	github.com/cdvelop/input v0.0.63 // indirect
	github.com/cdvelop/maps v0.0.8 // indirect
	github.com/cdvelop/strings v0.0.9 // indirect
	github.com/cdvelop/timetools v0.0.28 // indirect
)

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/maps => ../maps

replace github.com/cdvelop/unixid => ../unixid

replace github.com/cdvelop/timetools => ../timetools

replace github.com/cdvelop/timeserver => ../timeserver

replace github.com/cdvelop/input => ../input

replace github.com/cdvelop/objectdb => ../objectdb

replace github.com/cdvelop/dbtools => ../dbtools
