package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

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
	// Abre em formato read only
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

	// Configura e cria a directoria onde a schema traduzida para go, vai estar
	// Devolve o ficheiro criado
	ficheirosGO, err := genhelperfuncs.SetUpFilesAndDirs()
	if err != nil {
		fmt.Println("Erro: ", err)
	}
	defer ficheirosGO[0].Close()
	defer ficheirosGO[1].Close()

	genhelperfuncs.WritePackageNameOnFiles(ficheirosGO)

	ficheiroStructs := ficheirosGO[0]
	ficheiroFuncs := ficheirosGO[1]

	// Itera por todos os items insseridos em lista
	for _, v := range lista {
		// Verifica se o header deste item está correto
		if v.StructHeader["head"] == "tipo" {
			// Traduz o conteudo dentro do elemento da lista, para valores e elementos Go
			err := genhelperfuncs.ParseStructHeader(v, ficheiroStructs)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			err = genhelperfuncs.GenerateConvertFunc(v, ficheiroFuncs)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		} else {
			// Erro se o header não cumprir as regras de formatação
			fmt.Println("Erro: ", v.StructHeader["head"], "-> é irreconhecivél.")
			return
		}

		// Itera pelo conteúdo do corpo do elemento de cada item contido lista
		for _, v := range v.StructBody {
			currentType := v["type"]
			// Verifica o tipo de campo corrente e parsa para o equivalente em Go
			if len(regexp.MustCompile(`^lista\s+\w+$`).FindAllString(currentType, -1)) != 0 {
				err := genhelperfuncs.ParseStructBodyArray(v, currentType, ficheiroStructs)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				continue
			}
			// Verifica o tipo de campo corrente e parsa para o equivalente em Go
			if len(regexp.MustCompile(`^[A-z]+\s>\s[A-z]+$`).FindAllString(currentType, -1)) != 0 {
				err := genhelperfuncs.ParseStructBodyMap(v, ficheiroStructs)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
				continue
			}
			// Verifica o tipo de campo corrente e parsa para o equivalente em Go
			err := genhelperfuncs.ParseStructDefaultField(v, ficheiroStructs)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}
		// Encerra a struct com '}\n' para se poder passar á próxima
		written, err := genhelperfuncs.WriteBuffer("}\n", ficheiroStructs)
		if err == nil && written > 0 {
			fmt.Println("Wrriten - OK")
		} else {
			fmt.Println("Erro: ", err, written)
		}
	}
}
