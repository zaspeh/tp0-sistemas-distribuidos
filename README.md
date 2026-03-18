# TP0: Docker + Comunicaciones + Concurrencia

Alumno: Matias Kaled Dib  
Padrón: 111078  

A continuación se describe cómo ejecutar (o lo que se realizó en) cada uno de los ejercicios del trabajo práctico.

### Ejercicio N°1:

Ejecutando el comando `bash ./generar-compose.sh docker-compose-dev.yaml 5` en la raíz del proyecto se ejecuta un script que recibe como parámetros el nombre del archivo de salida y la cantidad de clientes esperados.

Este script permite configurar dinámicamente la cantidad de clientes que se desean levantar al ejecutar el sistema, ya que genera/modifica el archivo `docker-compose-dev.yaml`, el cual contiene la configuración utilizada por Docker Compose para inicializar los contenedores.

### Ejercicio N°2:

En este ejercicio se añadieron `volumes` en el archivo de configuración `docker-compose-dev.yaml` con el objetivo de poder modificar la configuración del cliente y del servidor sin necesidad de reconstruir las imágenes de Docker cada vez que se realizan cambios en los archivos de configuración.

### Ejercicio N°3:

Para verificar el correcto funcionamiento del servidor se implementó el script `validar-echo-server.sh`.

Al ejecutar `bash ./validar-echo-server.sh`, el script envía un mensaje de prueba al servidor utilizando `nc` desde un contenedor BusyBox conectado a la misma red de Docker.

Luego se compara la respuesta recibida con el mensaje enviado. Si ambos coinciden, el script imprime `action: test_echo_server | result: success`; de lo contrario, indica `result: fail`.

### Ejercicio N°4:

En este ejercicio se modificaron tanto el cliente como el servidor para manejar correctamente la señal `SIGTERM` y finalizar su ejecución de forma graceful.

Cuando el proceso recibe la señal `SIGTERM` (por ejemplo al ejecutar `docker compose down`), se ejecuta un handler (tanto en el cliente como en el servidor) que se encarga de cerrar correctamente los recursos abiertos antes de finalizar la aplicación.
