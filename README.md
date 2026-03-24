# TP0: Docker + Comunicaciones + Concurrencia

Alumno: Matias Kaled Dib  
Padrón: 111078  

A continuación se describe cómo ejecutar (o lo que se realizó en) cada uno de los ejercicios del trabajo práctico.

### Ejercicio N°1:

Ejecutando el comando `bash ./generar-compose.sh <nombre-archivo-de-salida> <cantidad-clientes>` en la raíz del proyecto se ejecuta un script el cual permite configurar dinámicamente la cantidad de clientes que se desean levantar al ejecutar el sistema, ya que genera/modifica el archivo `docker-compose-dev.yaml` (en caso de ser ese el archivo de salida indicado), el cual contiene la configuración utilizada por Docker Compose para inicializar los contenedores.

### Ejercicio N°2:

En este ejercicio se añadieron `volumes` en el archivo de configuración generado al ejecutar el comando del ejercicio1 con el objetivo de poder modificar la configuración del cliente y del servidor sin necesidad de reconstruir las imágenes de Docker cada vez que se realizan cambios en los archivos de configuración.

### Ejercicio N°3:

Para verificar el correcto funcionamiento del servidor se implementó el script `validar-echo-server.sh`.

Al ejecutar `bash ./validar-echo-server.sh`, el script envía un mensaje de prueba al servidor utilizando `nc` desde un contenedor BusyBox conectado a la misma red de Docker.

Luego se compara la respuesta recibida con el mensaje enviado. Si ambos coinciden, el script imprime `action: test_echo_server | result: success`; de lo contrario, indica `result: fail`.

### Ejercicio N°4:

En este ejercicio se modificaron tanto el cliente como el servidor para manejar correctamente la señal `SIGTERM` y finalizar su ejecución de forma graceful.

Cuando el proceso recibe la señal `SIGTERM` (por ejemplo al ejecutar `docker compose down`), se ejecuta un handler (tanto en el cliente como en el servidor) que se encarga de cerrar correctamente los recursos abiertos antes de finalizar la aplicación.

### Ejercicio N°5:

En este ejercicio se implementó la lógica de negocio para simular un sistema de apuestas donde cada cliente puede enviar su apuesta al servidor.

Para enviar las apuestas, el protocolo las serializa de la siguiente manera:
`client_id;nombre;apellido;dni;nacimiento;numero\n`
Por ejemplo:
`1;Santiago;Lorca;30904465;1999-03-17;7574`

A su vez se tuvo en cuenta el manejo de short write y se loggea todo lo que ocurre como fue pedido.

### Ejercicio N°6:

En este ejercicio se extiende la solución anterior incorporando procesamiento por **batches (chunks)**, permitiendo a los clientes enviar múltiples apuestas en una sola comunicación con el servidor, acortando tiempos de transmisión y procesamiento.

El cliente se encarga de leer un archivo con apuestas `.data/agency-{N}.csv` (donde `N` es el número del cliente). A medida que recorre el archivo, va construyendo batches y enviándolos al servidor.  
Luego de cada envío, espera la confirmación del servidor y reintenta el envío en caso de error.

Para construir los batches se tienen en cuenta dos restricciones:
- Un límite máximo de apuestas configurable (`batch.maxAmount`)
- Un límite de tamaño de mensaje de **8KB**

Esto se refleja en la condición: `if len(batch) >= maxAmount || currentSize+betSize > MaxBatchBytes`

Por otro lado, en este ejercicio se modificó el protocolo de comunicación agregando un header al inicio de cada mensaje para evitar problemas de lectura parcial (short read).
El formato es: `LEN:<tamaño>\n`

Luego del header se envían las apuestas serializadas de la misma forma que en el ejercicio anterior.

Este enfoque permite al receptor conocer exactamente cuántos bytes debe leer.

### Ejercicio N°7:

En este ejercicio se extiende la solución anterior incorporando coordinación entre clientes y servidor para realizar el sorteo y consultar ganadores.

Luego de enviar todos los batches, cada cliente:
1. Notifica al servidor que finalizó el envío de apuestas (`DONE`)
2. Realiza una consulta por los ganadores de su agencia (`ASK_WINNERS`)
3. Espera la respuesta y loguea la cantidad de ganadores obtenidos

Finalmente el cliente imprime: `action: consulta_ganadores | result: success | cant_ganadores: N`

Además, en este ejercicio se extiende el protocolo incorporando tipos de mensaje, el header ahora esta compuesto por: `LEN:<tamaño>;TYPE:<tipo>\n`

Del lado del servidor, se incorpora lógica de sincronización para garantizar que el sorteo se realice únicamente cuando todas las agencias hayan terminado de enviar sus apuestas.
Cada vez que un cliente envía DONE, se ejecuta: `self.finished_clients.append(client_sock)`. Entonces cuando la cantidad de clientes finalizados alcanza el total esperado, se sortean las apuestas y se logea: `action: sorteo | result: success`.

A su vez si un cliente consulta ganadores antes de que el sorteo haya finalizado, su request queda en espera (una vez realizado el sorteo se envían los ganadores a todos aquellos que quedaron en espera).

Los ganadores que se envían a cada agencia son los de su correspondiente agencia, cumpliendo así con el enunciado "No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.".

### Ejercicio N°8:

En este ejercicio se modifica el servidor para permitir aceptar conexiones y procesar mensajes de forma concurrente mediante el uso de **multithreading**.

Por cada nueva conexión entrante, el servidor crea un nuevo thread encargado de manejar la comunicación con ese cliente. De esta forma, múltiples clientes pueden interactuar con el servidor de manera concurrente, sin bloquear la aceptación de nuevas conexiones.

Dado que múltiples threads acceden a estructuras compartidas, se incorporan mecanismos de sincronización para evitar race-conditions (condiciones de carrera).
Se utilizan dos herramientas principales:
* Lock: para proteger secciones críticas (por ejemplo cuando se hace `store_bets(..)`).
* Condition: para coordinar eventos entre threads (por ejemplo para coordinar el momento en el que se realiza el sorteo).

A su vez se separaron dos locks dependiendo de su uso, por un lado `file_lock` el cual se encarga de cubrir los registros de las apuestas y por el otro `lock` el cual era usado por la `Condition` mencionada previamente.
Esta `Condition` también se utilizó al momento de consultar los ganadores del sorteo, de forma tal que se prevenga el spurious wakeups.  
