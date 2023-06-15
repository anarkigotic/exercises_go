# Resumen de Compras

Este programa en Go implementa un servidor HTTP que proporciona un endpoint para obtener un resumen de compras. Utiliza la biblioteca estándar `net/http` para manejar las solicitudes HTTP y `encoding/json` para serializar los datos en formato JSON.

## Requisitos

- Go 1.15 o superior

## Uso

1. Asegúrate de tener acceso a un servidor API que proporcione datos de compras en formato JSON. El programa utiliza la función `fetchData` para obtener los datos de cada día a través de una solicitud HTTP.

2. Abre el archivo `main.go` y modifica la URL base en la línea `urlBase := "https://apirecruit-gjvkhl2c6a-uc.a.run.app/compras/"` para que coincida con la URL de tu servidor API.

3. Ejecuta el programa utilizando el siguiente comando:

   ```bash
   go run main.go
