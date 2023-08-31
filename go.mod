module github.com/cdvelop/postgre

go 1.20

require (
	github.com/cdvelop/dbtools v0.0.25
	github.com/cdvelop/model v0.0.43
	github.com/cdvelop/objectdb v0.0.37
	github.com/lib/pq v1.10.9
)

require (
	github.com/cdvelop/gotools v0.0.30 // indirect
	github.com/cdvelop/input v0.0.26 // indirect
)

replace github.com/cdvelop/model => ../model

replace github.com/cdvelop/input => ../input

replace github.com/cdvelop/objectdb => ../objectdb

replace github.com/cdvelop/dbtools => ../dbtools
