package datastructs

// Estrutura - Representa a estrutura de uma golang struct
// dividída em head e body
type Estrutura struct {
	StructHeader map[string]string
	StructBody   []map[string]string
}

// FileStructs - A lista das structs possiveis de extrair do ficheiro schema
type FileStructs []Estrutura
