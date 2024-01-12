package postgre

import (
	"fmt"

	"log"
	"os"
	"os/exec"
	"runtime"
)

func (d *PG) BackupDataBase() {

	switch runtime.GOOS {
	case "windows":
		d.backupWindowsDB()
	default:
		log.Println("error backup postgres db en sistema operativo: " + runtime.GOOS + " no implementado")
	}

}

func (d *PG) backupWindowsDB() {
	//pg_dump.exe -Fc postgresql://pa100t:73D-s137v4sDB@127.0.0.1:5432/pa100t > D:\postgres_backup\pa100t.backup
	const e = "backupWindowsDB error "
	pg_dump := fmt.Sprintf("%v/bin/pg_dump.exe", installation_directory)

	name_file_backup, err := d.idUnix.GetNewID()
	if err != "" {
		log.Println(e + err)
		return
	}
	name_file_backup += "-db.backup"

	destination_directory := d.BackupDirectory + "/" + name_file_backup
	// cmd := exec.Command("D:/Program Files/FreeFileSync/FreeFileSync.exe", "D:/cesar/SyncWin/SyncSettings.ffs_batch")

	out, er := exec.Command(pg_dump, "-Fc", d.ConnectionString()).Output()
	if er != nil {
		log.Printf("¡error en el comando de respaldo! %s\n", er)
		return
	}

	er = os.WriteFile(destination_directory, out, 0666)
	if er != nil {
		log.Printf(e+"¡al guardar archivo respaldo! %v\n", er)
		return
	}

	// log.Printf(">>> backup base de datos %v\n", name_file_backup)
	d.maintenanceBackupDirectory()

}

func (d *PG) RestoreDataBase(file_name string) bool {

	pg_restore := fmt.Sprintf("%v/bin/pg_restore.exe", installation_directory)
	//pg_restore [Connection-option…] [option…] [filename]

	dbname := fmt.Sprintf("--dbname=postgres://postgres:%v@127.0.0.1:5432/postgres?sslmode=disable", d.passwordDB)

	// command := fmt.Sprintf(`"%v" %v --create %v`, pg_restore, dbname, file_name)
	// fmt.Printf("CADENA A CONECTAR:\n%v", command)

	out, err := exec.Command(pg_restore, dbname, "--create", file_name).Output()

	if err != nil {
		log.Printf("ERROR COMANDO DE RESTAURAR DB %v [%v]", err, out)
		return false
	}

	fmt.Printf(">>> BASE DE DATOS %v RESTAURADA OK\n", file_name)

	return true
}

func (d *PG) maintenanceBackupDirectory() {

	const e = "maintenanceBackupDirectory error "

	oldest_file, total_files, err := findOldestFile(d.BackupDirectory)
	if err != "" {
		log.Println(e + err)
		return
	}

	if total_files > d.TotalBackupsMaintain {
		er := os.Remove(d.BackupDirectory + "/" + oldest_file.Name())
		if er != nil {
			log.Println(e + "¡al eliminar fichero: " + oldest_file.Name() + " respaldo anterior! " + er.Error())
		}
	}
}
