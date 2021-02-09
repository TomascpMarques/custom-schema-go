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

// SetUpFilesAndDirs Cria a dir e o ficheiro onde a schema traduzida para go estará
func SetUpFilesAndDirs() (*os.File, error) {
	// Cria a dir onde o ficheiro final estará
	// Permissões: Owner: rwx; Others: r-x
	dirErr := os.MkdirAll("./resolvedschema", 0755)
	if dirErr != nil {
		fmt.Println(dirErr)
		return nil, dirErr
	}
	// Cria o ficheiro com as structs
	file, err := os.Create("resolvedschema/schema_structs.go")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return file, nil
}

// WriteBuffer Escreve no ficheiro defenido, o string passada, a string passada é esperada que seja já formatada
func WriteBuffer(writeValue string, file *os.File) (int, error) {
	// Escreve no ficheiro, devolve o nº de bytes escritos e um erro
	written, err := file.Write([]byte(writeValue))
	// Compára-se o erro e o número de bytes escritos é o esperado
	if err == nil && written == len([]byte(writeValue)) {
		return written, err
	}
	// Fecha o ficheiro, no fim de todas as operações,
	// Também fecha o ficheiro em caso de erros
	defer file.Close()
	return 0, err
}

// ParseStructHeader Cria o header em Go de uma Estrutura passada nos parametros
func ParseStructHeader(v datastructs.Estrutura, file *os.File) error {
	// Header em Go da struct, conteúdo extraído do parametro v
	line := fmt.Sprintf("// %s - \ntype %s struct {\n", v.StructHeader["body"], v.StructHeader["body"])

	// Escreve no ficheiro o struct header
	written, err := WriteBuffer(line, file)
	if err != nil || written < len(v.StructHeader["head"]) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

// ParseStructBodyArray Através do body da struct extraido pelo programa cria um go array
func ParseStructBodyArray(v map[string]string, currentType string, file *os.File) error {
	//currentType[6:] -> salta os chars na schema, esses são "lista "
	line := fmt.Sprintf("%s []%s `json:\"%s\"`\n", v["field"], currentType[6:], strings.ToLower(v["field"]))

	written, err := WriteBuffer(line, file)
	if err != nil || written < len(line) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

// ParseStructBodyMap Cria um Golang map com os dados fornecidos por v, o formato normalmente é - tipo > tipo -
func ParseStructBodyMap(v map[string]string, file *os.File) error {
	firstType := v["type"][:strings.Index(v["type"], " ")]
	secondType := v["type"][strings.Index(v["type"], " ")+3:]

	line := fmt.Sprintf("%s map[%s]%s `json:\"%s\"`\n", v["field"], firstType, secondType, strings.ToLower(v["field"]))

	written, err := WriteBuffer(line, file)
	if err != nil || written < len(line) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}

// ParseStructDefaultField Cria os campos das structs default, os que só são primitivos, ex:-> Nome string
func ParseStructDefaultField(v map[string]string, file *os.File) error {
	line := fmt.Sprintf("\t%s %s `json:\"%s\"`\n", v["field"], v["type"], strings.ToLower(v["field"]))
	written, err := WriteBuffer(line, file)
	if err != nil || written < len(line) {
		fmt.Println("Error: ", err)
		return err
	}
	return nil
}
