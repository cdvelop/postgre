package postgre

// SwapConnect intercambia a la conexión por defecto del motor de la base de datos
func (d *PG) SwapConnect() (backupDNS PG) {
	backupDNS = *d
	d.dataBaseName = d.DataBasEngine()
	d.userDatabase = d.DataBasEngine()
	return
}
