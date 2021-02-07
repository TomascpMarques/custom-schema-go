package resolvedschema

// InfoComputador -
type InfoComputador struct {
	ID      string
	Nome    string
	Sala    string
	Formato string
	Ativo   bool
	Area    string
}

// ComponentesComputador -
type ComponentesComputador struct {
	CPU    []CPU
	GPU    []GPU
	RAM    []RAM
	MBoard Mboard
}

// Mboard -
type Mboard struct {
	ID            string
	Nome          string
	RAMSlots      int
	MaxRAM        int
	SlotsExpanso  map[string]int
	Armazenamento map[string]int
	Lan           string
	USBports      map[string]int
	PortasIO      map[string]int
	Formato       string
	Notas         string
}

// CPU -
type CPU struct {
	ID         string
	Nome       string
	Nucleos    int
	Threads    int
	Frequencia string
	MaxRAM     int
	TDP        string
	Socket     string
	Marca      string
}

// GPU -
type GPU struct {
	ID            string
	Nome          string
	VelocidadeMem string
	VRAM          int
	Dimensoes     string
	HDIM          int
	VGA           int
	DisplayPort   int
}

// RAM -
type RAM struct {
	Nome       string
	Velocidade string
	Memoria    string
	Tipo       string
}

// Computador -
type Computador struct {
	ID   string
	Info InfoComputador
}
