# Postgres Setup & Troubleshooting

Este documento detalla cómo solucionar problemas de conexión con la base de datos Postgres en este proyecto, basándose en el incidente resuelto el 02 de Marzo de 2026.

## Síntoma
Error al ejecutar `gotest`:
`Failed to connect: pq: password authentication failed for user "postgres"`

## Diagnóstico y Solución

### 1. Identificar el Entorno
Primero, verifica si Postgres está corriendo nativamente o en Docker:
```bash
docker ps | grep postgres
```
Si aparece un contenedor (ej. `pg14-dev`), la base de datos está aislada.

### 2. Sincronizar Contraseña

#### Si usas Docker:
Si el `.env` define `postgres` como contraseña, pero el contenedor tiene otra, ejecuta:
```bash
docker exec pg14-dev psql -U postgres -c "ALTER USER postgres WITH PASSWORD 'postgres';"
```

#### Si es una instalación nativa (Linux):
Usa `sudo` para entrar como el usuario del sistema `postgres`:
```bash
sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'postgres';"
```

> [!IMPORTANT]
> Si recibes un error de autenticación `peer` incluso con la contraseña correcta, revisa tu archivo `/etc/postgresql/XX/main/pg_hba.conf` y asegúrate de que el método para `local` esté en `md5` o `scram-sha-256`.

### 3. Verificar Fallbacks en Tests
Asegúrate de que los archivos de test no tengan contraseñas hardcodeadas que ignoren el `.env`.
Archivos críticos a revisar:
- `tests/ddl_test.go`
- `tests/adapter_test.go`

Los tests deben usar `os.Getenv("POSTGRES_DSN")` y, si está vacío, un fallback consistente con el `.env` del proyecto.

## Configuración Recomendada (.env)
```env
POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
```
