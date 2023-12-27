module github.com/cdvelop/postgre

go 1.20

require (
	github.com/cdvelop/model v0.0.107
	github.com/cdvelop/objectdb v0.0.110
	github.com/cdvelop/timeserver v0.0.31
	github.com/cdvelop/unixid v0.0.47
	github.com/lib/pq v1.10.9
)

require (
	github.com/cdvelop/dbtools v0.0.80 // indirect
	github.com/cdvelop/input v0.0.80 // indirect
	github.com/cdvelop/maps v0.0.8 // indirect
	github.com/cdvelop/strings v0.0.9 // indirect
	github.com/cdvelop/timetools v0.0.32 // indirect
)

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/objectdb => ../objectdb
