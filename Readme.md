# Custom Schema Go

> Módulo do projeto robin backend - PAP Tomás Marques

Através de um ficheiro schema fornecido ao executável como paramêtro, o programa lê e cria a schema desse ficheiro, em Golang structs equivalentes, que podem ser importadas por outros ficheiros.

## Install:
``` go get github.com/tomascpmarques/custom-schema-go ```

## Execução:
``` go run github.com/tomascpmarques/custom-schema-go/init.go ~caminho do ficheiro~  ```

## Exemplo:
### Ficheiro: /Example_Schemas/schema.txt
```
. . .

tipo GPU
    ID string
    Nome string
    VelocidadeMem string
    VRAM int
    Dimensoes string
    HDIM int
    VGA int
    DisplayPort int

tipo RAM
    Nome string
    Velocidade string
    Memoria string
    Tipo string

tipo Computador 
    ID string
    Info InfoComputador
    Caracteristicas ComponentesComputador

. . .
```
