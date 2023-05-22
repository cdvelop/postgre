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
		d.backupDBWindows()
	default:
		log.Println("BACKUP POSTGRES LINUX NO IMPLEMENTADO")
	}

}

func (d *PG) backupDBWindows() {
	//pg_dump.exe -Fc postgresql://pa100t:73D-s137v4sDB@127.0.0.1:5432/pa100t > D:\postgres_backup\pa100t.backup

	pg_dump := fmt.Sprintf("%v/bin/pg_dump.exe", installation_directory)

	name_file_backup := d.GetNewID() + "-d.backup"

	destination_directory := d.backup_directory + "/" + name_file_backup
	// cmd := exec.Command("D:/Program Files/FreeFileSync/FreeFileSync.exe", "D:/cesar/SyncWin/SyncSettings.ffs_batch")

	out, err := exec.Command(pg_dump, "-Fc", d.ConnectionString()).Output()
	if err != nil {
		fmt.Printf("¡ERROR EN EL COMANDO RESPALDO BASE DE DATOS! %s\n", err)
	} else {
		err = os.WriteFile(destination_directory, out, 0666)
		if err != nil {
			log.Printf("¡ERROR AL GUARDAR ARCHIVO RESPALDO DB! %v\n", err)
		} else {
			log.Printf(">>> BACKUP BASE DE DATOS %v\n", name_file_backup)
			d.maintenanceBackupDirectory()
		}
	}
}

func (d *PG) RestoreDataBase(file_name string) bool {

	pg_restore := fmt.Sprintf("%v/bin/pg_restore.exe", installation_directory)
	//pg_restore [Connection-option…] [option…] [filename]

	dbname := fmt.Sprintf("--dbname=postgres://postgres:%v@127.0.0.1:5432/postgres?sslmode=disable", d.passwordDatabase)

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

	oldest_file, total_files, err := findOldestFile(d.backup_directory)
	if err != nil {
		log.Println("¡ERROR AL OBTENER FICHERO MAS ANTIGUO DEL DIRECTORIO DB BACKUP! " + err.Error())
	}

	if total_files > 5 { //solo mantener 5 respaldos
		e := os.Remove(d.backup_directory + "/" + oldest_file.Name())
		if e != nil {
			log.Println("¡ERROR AL ELIMINAR FICHERO: " + oldest_file.Name() + " RESPALDO ANTIGUO! " + e.Error())
		}
	}
}
