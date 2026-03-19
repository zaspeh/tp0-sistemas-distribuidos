package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {

	outputFile := os.Args[1]
	numClients, err := strconv.Atoi(os.Args[2])
	if err != nil || numClients < 1 { // valido cantidad de clientes inválida
		fmt.Println("Cantidad de clientes inválida")
		os.Exit(1)
	}

	file, err := os.Create(outputFile)
	if err != nil { // valido error de creación
		panic(err)
	}
	defer file.Close() // al final de la función cierro el archivo, garantiza cierre aunque haya un panic

	writeServer(file)
	writeClients(file, numClients)
	writeNetworks(file)
}

// escribo la definición del servidor
func writeServer(file *os.File) {
	fmt.Fprintln(file, `services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    networks:
      - testing_net
    volumes: 
      - ./server/config.ini:/config.ini:ro
  `) // al añadir 'volumes: xxx' aseguro que use el config.ini dentro de la carpeta del host y no dentro del de la imagen
}

// escribo la definición de los clientes
func writeClients(file *os.File, numClients int) {
	for i := 1; i <= numClients; i++ {
		fmt.Fprintf(file, `
  client%d:
    container_name: client%d
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=%d
      - NOMBRE=cliente
      - APELLIDO=%d
      - DOCUMENTO=%d
      - NACIMIENTO=1990-01-01
      - NUMERO=%d
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml:ro
    depends_on:
      - server
`, i, i, i, i, i, i)
	}
}

// agrego la sección de redes
func writeNetworks(file *os.File) {
	fmt.Fprintln(file, `
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24`)
}
