package commonsfake

type TableBuildtarjeta struct {
	Table TableDriverBuildTarjeta
}

type TableDriverBuildTarjeta struct {
	TituloPrueba string
	WantTable    bool
	Tarjeta      string
}

// const ERROR_CALCULO_COMISION = "error de validación: no se pudo obtener calculo de comisiones"
type TableDriverBuildDiferenceUint struct {
	TituloPrueba string
	WantTable    bool
	String1      []uint64
	String2      []uint64
}
