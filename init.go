package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/tomascpmarques/custom-schema-go/datastructs"
	"github.com/tomascpmarques/custom-schema-go/genhelperfuncs"
)

func main() {
	// Prepara a lisa de Estruturas a criar
	var lista = make(datastructs.FileStructs, 0)

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
		// Mapea os valores presentes na linha corrente, para a var lista
		genhelperfuncs.MapearVetorEstruturas(&lista, val, &i)
		// Verifica se o reader chegou ao fim do ficheiro
		if err == io.EOF {
			fmt.Println("Read to End of file")
			break
		}
		// Verifica para quais queres outros erros
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	ficheiroGO, err := genhelperfuncs.SetUpFilesAndDirs()
	if err != nil {
		fmt.Println("Erro: ", err)
	}
	defer ficheiroGO.Close()

	written, err := ficheiroGO.WriteString("package resolvedschema\n")
	if err == nil && written == len("package resolvedschema\n") {
		fmt.Println("Wrriten - OK")
	}

	for _, v := range lista {
		if v.StructHeader["head"] == "tipo" {
			line := "// " + v.StructHeader["body"] + " -" + "\ntype " + v.StructHeader["body"] + " struct {\n"
			written, err := genhelperfuncs.WriteBuffer(line, ficheiroGO)
			if err != nil || written < len(v.StructHeader["head"]) {
				fmt.Println("Error: ", err)
				return
			}
		} else {
			fmt.Println("Erro: ", v.StructHeader["head"], "irreconhecivél.")
			return
		}
		for _, v := range v.StructBody {
			currentType := v["type"]
			if len(regexp.MustCompile(`^lista\s+\w+$`).FindAllString(currentType, -1)) != 0 {
				line := fmt.Sprintf("%s []%s\n", v["field"], currentType[6:])
				written, err := genhelperfuncs.WriteBuffer(line, ficheiroGO)
				if err != nil || written < len(line) {
					fmt.Println("Error: ", err)
					return
				}
				continue
			}
			if len(regexp.MustCompile(`^[A-z]+\s>\s[A-z]+$`).FindAllString(currentType, -1)) != 0 {
				firstType := v["type"][:strings.Index(v["type"], " ")]
				secondType := v["type"][strings.Index(v["type"], " ")+3:]
				line := fmt.Sprintf("%s map[%s]%s\n", v["field"], firstType, secondType)
				written, err := genhelperfuncs.WriteBuffer(line, ficheiroGO)
				if err != nil || written < len(line) {
					fmt.Println("Error: ", err)
					return
				}
				continue
			}
			line := fmt.Sprintf("\t%s %s\n", v["field"], v["type"])
			written, err := genhelperfuncs.WriteBuffer(line, ficheiroGO)
			if err != nil || written < len(line) {
				fmt.Println("Error: ", err)
				return
			}
		}
		written, err := genhelperfuncs.WriteBuffer("}\n", ficheiroGO)
		if err == nil && written > 0 {
			fmt.Println("Wrriten - OK")
		} else {
			fmt.Println("Erro: ", err, written)
		}
	}
}
