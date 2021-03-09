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
func SetUpFilesAndDirs() (files []*os.File, err error) {
	// Cria a dir onde o ficheiro final estará
	// Permissões: Owner: rwx; Others: r-x
	dirErr := os.MkdirAll("./resolvedschema", 0755)
	if dirErr != nil {
		fmt.Println(dirErr)
		return nil, dirErr
	}
	// Cria o ficheiro com as structs
	fileStructs, err := os.Create("resolvedschema/schema_structs.go")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	files = append(files, fileStructs)

	// Cria o ficheiro com as structs
	fileFuncs, err := os.Create("resolvedschema/schema_gen_funcs.go")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	files = append(files, fileFuncs)

	return files, nil
}

// WritePackageNameOnFiles Escre o nome do package em todos os ficheiros dentro do vetor files
func WritePackageNameOnFiles(files []*os.File) {
	nomePacote := "package resolvedschema\n"
	for _, file := range files {
		// Inssere o nome do package no ficheiro
		written, err := file.WriteString(nomePacote)

		if err != nil || written != len(nomePacote) {
			fmt.Println("Wrriten - ERROR")
			fmt.Println("ERROR: ", err)
			panic("Fo impossivel atribuir um pacote aos novos ficheiros")
		}

		if file.Name() == "resolvedschema/schema_gen_funcs.go" {
			written, err := file.WriteString("\nimport(\"encoding/json\")\n")

			if err != nil || written != len("\nimport(\"encoding/json\")\n") {
				fmt.Println("Wrriten - ERROR")
				fmt.Println("ERROR: ", err)
				panic(err)
			}
		}

		fmt.Println("Wrriten - OK")
	}
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
	// Se a lenght do nome for menor ou igual a 3, capitaliza as letras, de acordo com os estilo da
	if len(v.StructHeader["body"]) <= 3 {
		v.StructHeader["body"] = strings.ToUpper(v.StructHeader["body"])
	}

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

/*
	// FromMapToStruct -
	func FromMapToStruct(param1 *map[string]interface{}) interface{} {
		var exampStruct Examplestruct
		temp, _ := json.Marshal(param1)
		_ = json.Unmarshal(temp, &exampStruct)

		return exampStruct
	}
*/

// GenerateConvertFunc Cria a função que converte de um map[string]interface{} para a struct adequada
func GenerateConvertFunc(v datastructs.Estrutura, file *os.File) error {
	funcDef := fmt.Sprintf(
		"\nfunc %sParaStruct(param1 *map[string]interface{}) %s {\n",
		v.StructHeader["body"], v.StructHeader["body"])

	// Escreve no ficheiro a defenição da struct
	written, err := WriteBuffer(funcDef, file)
	if err != nil || written < len(funcDef) {
		fmt.Println("Erro na defenição da : ", err)
		return err
	}

	funcBody := []string{
		"\tvar returnStruct " + v.StructHeader["body"] + "\n",
		"\ttemp, err := json.Marshal(param1)\n",
		"\tif err != nil {\n\t\treturn ", v.StructHeader["body"], "{} \n\t}\n",
		"\terr = json.Unmarshal(temp, &returnStruct)\n",
		"\tif err != nil {\n\t\treturn ", v.StructHeader["body"], "{}\n\t}\n",
		"\treturn returnStruct\n",
		"}",
	}

	for _, value := range funcBody {
		// Escreve no ficheiro a linha respetiva do corpo da função
		written, err := WriteBuffer(value, file)
		if err != nil || written < len(value) {
			fmt.Println("Erro na defenição da : ", err)
			return err
		}
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
