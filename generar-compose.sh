#!/bin/bash

output_file=$1
num_clients=$2

echo "Nombre del archivo de salida: $output_file"
echo "Cantidad de clientes: $num_clients"

go run tools/mi-generador.go $output_file $num_clients

echo "Archivo generado correctamente."