# Conversor de datos CSV a JSON

Este programa en Go lee un archivo CSV y lo convierte a formato JSON. Utiliza la biblioteca estándar `encoding/csv` para leer los datos del archivo CSV y `encoding/json` para convertirlos a JSON.

## Requisitos

- Go 1.15 o superior

## Uso

1. Asegúrate de tener un archivo CSV con los datos que deseas convertir. El archivo debe tener una estructura similar a la siguiente:

   ```csv
   organizacion,usuario,rol
   org1,jperez,admin
   org1,jperez,superadmin
   org1,asosa,writer
   org2,jperez,admin
   org2,rrodriguez,writer
   org2,rrodriguez,editor
