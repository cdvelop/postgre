package postgre

// SwapConnect intercambia a la conexión por defecto del motor de la base de datos
func (d *PG) SwapConnect() (backupDNS PG) {
	backupDNS = *d
	d.DBName = d.DataBasEngine()
	d.UserDB = d.DataBasEngine()
	return
}
