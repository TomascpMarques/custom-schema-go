package datastructs

import "os"

// Estrutura -
type Estrutura struct {
	StructHeader map[string]string
	StructBody   []map[string]string
}

// FileStructs -
type FileStructs []Estrutura

// SchemaGOFile -
type SchemaGOFile os.File
