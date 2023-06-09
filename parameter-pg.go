package postgre

func (d *PG) SQLTableInfo() string {
	return "SELECT column_name FROM information_schema.columns WHERE table_name = '%v';"
}

func (d *PG) SQLColumName() string {
	return "column_name"
}

func (d *PG) SQLDropTable() string {
	return "DROP TABLE IF EXISTS %v CASCADE;" //sql de eliminación de tabla
}

func (d *PG) DBExists() string {
	return "SELECT 1 FROM pg_database WHERE datname='%v';"
}
