module github.com/cdvelop/postgre

go 1.20

require (
	github.com/cdvelop/dbtools v0.0.24
	github.com/cdvelop/model v0.0.34
	github.com/cdvelop/objectdb v0.0.32
	github.com/lib/pq v1.10.9
)

require (
	github.com/cdvelop/gotools v0.0.16 // indirect
	github.com/cdvelop/input v0.0.15 // indirect
	golang.org/x/text v0.11.0 // indirect
)

// replace github.com/cdvelop/model => ../model

// replace github.com/cdvelop/input => ../input
replace github.com/cdvelop/objectdb => ../objectdb

replace github.com/cdvelop/dbtools => ../dbtools
