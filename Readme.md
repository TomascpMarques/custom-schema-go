# Custom Schema Go

> Módulo do projeto robin backend - PAP Tomás Marques

Através de um ficheiro schema fornecido ao executável como parâmetro, o programa lê e cria a schema desse ficheiro, em Golang structs equivalentes, que podem ser importadas por outros ficheiros.

> De momento só funciona na pasta do projeto e em contentores docker que clonem o projeto, de resto, se chamar o .exe, o output sai todo comido.

## Install:
``` go get github.com/tomascpmarques/custom-schema-go ```

## Execução:
``` go run github.com/tomascpmarques/custom-schema-go/ ~caminho do ficheiro schema~  ```

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
