package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Estrutura -
type Estrutura struct {
	StructHeader map[string]string
	StructBody   []map[string]string
}

// FileStructs -
type FileStructs []Estrutura

// SchemaGOFile -
type SchemaGOFile os.File

// MapearVetorEstruturas Extrai os valores presentes na string fornecida, expostos de uma maneira específica
// e mapea-os para um struct Estrutura, que será armazenada num vetor de Estrutura's
func (fs *FileStructs) MapearVetorEstruturas(linha string, counter *int) {
	// Retira os espaços e new-lines (prefixo,sufixo) da string de input
	lineTrim := strings.Trim(linha, " \n")

	// Se a linha não estiver vazia
	if len(lineTrim) > 1 {
		// Se a linha estivere a identificar um nova Estrutura com a keyword tipo
		if lineTrim[:4] == "tipo" {
			// Adicionar ao fim do vetor uma nova Estrutura
			*fs = append(*fs, Estrutura{
				StructHeader: map[string]string{
					// O header da struct são sempre os primeiros 4 caracters
					"head": linha[:4],

					// O corpo da defenição é o nome específicado pelo utilizador
					// após o espaço, que segue a parte "tipo"
					"body": strings.Trim(linha[5:], " \n"),
				},
				StructBody: make([]map[string]string, 0), // Inicialisa o vetor que vai conter os campos defenidos pelo o utilizador
			})
			*counter++
			return
		}
		// Regex para filtrar o formato correto dos campos
		if len(regexp.MustCompile(`\w+\s[a-z]+|\w+\s\[\][a-z]+|\w+\s\[\][A-z]+|\w+\s[A-z]+`).FindAllStringSubmatch(lineTrim, -1)) != 0 {
			// procura o primeiro espaço e usa o como divisor entre nome_campo | tipo_campo
			indexSpace := strings.Index(lineTrim, " ")
			// Atribuição temporária do valor apontado por fs
			temp := *fs
			// Utilisa o contador passado nos parametros para adicionar o corpo(fields) da Estrutura corrente
			// o contador
			temp[*counter].StructBody = append(temp[*counter].StructBody, map[string]string{
				"field": lineTrim[:indexSpace],   // O nome do campo são todos os caracteres antes do espaço
				"type":  lineTrim[indexSpace+1:], // O tipo do campo é tudo depois do 1º espaço
			})
			return
		}
		fmt.Println("Error")
		return
	}
}

func main() {
	// Prepara a lisa de Estruturas a criar
	var lista = make(FileStructs, 0)

	// Verifica o número de argumentos nada mais do que 1 (o primeiro é o nº de ar)
	if len(os.Args) != 2 {
		fmt.Println("Erro: numero de parametros é incorreto")
		return
	}

	// Ficheiro com a schema para ler
	ficheiro, err := os.OpenFile(os.Args[1], os.O_RDONLY, os.FileMode(os.O_RDONLY))
	if err != nil {
		fmt.Println("Erro: ", err)
		return
	}
	defer ficheiro.Close()

	// Bufferd reader
	reader := bufio.NewReader(ficheiro)
	// O i é = a -1, porque se fosse 0 logo de início, tinhamos de decrementar na função que mapea os valores
	// Corre enquanto não houver erros (err == nil)
	for i := -1; err == nil; {
		// Lé o ficheiro de uma maneira bufferd, até o new-line char
		val, err := reader.ReadString('\n')
		// Verifica se o reader chegou ao fim do ficheiro
		if err == io.EOF {
			fmt.Println("End of file")
			break
		}
		// Verifica para quais queres outros erros
		if err != nil {
			fmt.Println(err)
			return
		}
		// Mapea os valores presentes na linha corrente, para a var lista
		lista.MapearVetorEstruturas(val, &i)
	}

	dirErr := os.MkdirAll("resolvedschema", os.FileMode(os.O_APPEND))
	if dirErr != nil {
		fmt.Println(dirErr)
		return
	}
	ficheiroGO, err := os.Create("resolvedschema/schema_structs.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ficheiroGO.Close()

	written, err := ficheiroGO.WriteString("package resolvedschema\n")
	if err == nil && written > 0 {
		fmt.Println("OK")
	}

	var buffer string
	for _, v := range lista {
		switch v.StructHeader["head"] {
		case "tipo":
			line := "// " + v.StructHeader["body"] + " -" + "\ntype " + v.StructHeader["body"] + " struct {\n"
			WriteBuffer(line, ficheiroGO)
			break
		default:
			fmt.Println("Erro: ", v.StructHeader["head"], "irreconhecivél.")
			break
		}
		for _, v := range v.StructBody {
			switch v["type"] {
			case "string > int":
				written, err := WriteBuffer(fmt.Sprintf("\t%s map[string]int\n", v["field"]), ficheiroGO)
				if err == nil && written > 0 {
					fmt.Println("OK")
				} else {
					fmt.Println("Erro: ", err, written)
				}
				break
			default:
				written, err := WriteBuffer(fmt.Sprintf("\t%s %s\n", v["field"], v["type"]), ficheiroGO)
				if err == nil && written > 0 {
					fmt.Println("OK")
				} else {
					fmt.Println("Erro: ", err, written)
				}
				break
			}
		}
		written, err := WriteBuffer("}\n", ficheiroGO)
		if err == nil && written > 0 {
			fmt.Println("OK")
		} else {
			fmt.Println("Erro: ", err, written)
		}
	}
	fmt.Println("Buffer string: ", buffer)

	// buffer := []byte("package resolved_schema\n import(\"fmt\")")

	// debug - stdout os valores dentro da var lista
	for k, v := range lista {
		fmt.Println(strings.Repeat("-", 55))
		fmt.Printf("Element %d:\n", k)
		fmt.Println("Conttents: ", v)
	}
}

// WriteBuffer -
func WriteBuffer(writeValue string, file *os.File) (written int, err error) {
	written, err = file.Write([]byte(writeValue))
	if err == nil && written == len([]byte(writeValue)) {
		return written, err
	}
	return written, err
}
