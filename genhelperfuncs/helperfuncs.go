package genhelperfuncs

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tomascpmarques/custom-schema-go/datastructs"
)

// MapearVetorEstruturas Extrai os valores presentes na string fornecida, expostos de uma maneira específica
// e mapea-os para um struct Estrutura, que será armazenada num vetor de Estrutura's
func MapearVetorEstruturas(fs *datastructs.FileStructs, linha string, counter *int) {
	// Retira os espaços e new-lines (prefixo,sufixo) da string de input
	lineTrim := strings.Trim(linha, " \n")

	// Se a linha não estiver vazia
	if len(lineTrim) > 1 {
		// Se a linha estivere a identificar um nova Estrutura com a keyword tipo
		if lineTrim[:4] == "tipo" {
			// Adicionar ao fim do vetor uma nova Estrutura
			*fs = append(*fs, datastructs.Estrutura{
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

// SetUpFilesAndDirs -
func SetUpFilesAndDirs() (*os.File, error) {
	//os.Chdir("..")
	dirErr := os.MkdirAll("resolvedschema", os.FileMode(os.O_APPEND))
	if dirErr != nil {
		fmt.Println(dirErr)
		return nil, dirErr
	}
	file, err := os.Create("resolvedschema/schema_structs.go")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return file, nil
}

// WriteBuffer -
func WriteBuffer(writeValue string, file *os.File) (written int, err error) {
	written, err = file.Write([]byte(writeValue))
	if err == nil && written == len([]byte(writeValue)) {
		return written, err
	}
	return 0, err
}

// ParseStructHeader -
func ParseStructHeader(v datastructs.Estrutura, file *os.File) error {
	line := fmt.Sprintf("// %s - \ntype %s struct {\n", v.StructHeader["body"], v.StructHeader["body"])

	written, err := WriteBuffer(line, file)
	if err != nil || written < len(v.StructHeader["head"]) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

// ParseStructBodyArray -
func ParseStructBodyArray(v map[string]string, currentType string, file *os.File) error {
	line := fmt.Sprintf("%s []%s\n", v["field"], currentType[6:])

	written, err := WriteBuffer(line, file)
	if err != nil || written < len(line) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

// ParseStructBodyMap -
func ParseStructBodyMap(v map[string]string, file *os.File) error {
	firstType := v["type"][:strings.Index(v["type"], " ")]
	secondType := v["type"][strings.Index(v["type"], " ")+3:]

	line := fmt.Sprintf("%s map[%s]%s\n", v["field"], firstType, secondType)

	written, err := WriteBuffer(line, file)
	if err != nil || written < len(line) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

// ParseStructDefaultField -
func ParseStructDefaultField(v map[string]string, file *os.File) error {
	line := fmt.Sprintf("\t%s %s\n", v["field"], v["type"])
	written, err := WriteBuffer(line, file)
	if err != nil || written < len(line) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}
