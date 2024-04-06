package commonsfake

func EstructuraValidarTarjeta() (tableDriverTestPeyment []TableDriverBuildTarjeta) {

	date1 := TableDriverBuildTarjeta{
		TituloPrueba: "Validar numero de tarjeta Naranja Mary",
		WantTable:    true,
		Tarjeta:      "5895621687899586",
	}
	date2 := TableDriverBuildTarjeta{
		TituloPrueba: "Validar numero de tarjeta Naranja Julio",
		WantTable:    true,
		Tarjeta:      "5895628851930788",
	}
	date3 := TableDriverBuildTarjeta{
		TituloPrueba: "Validar numero de tarjeta Naranja Lucas",
		WantTable:    true,
		Tarjeta:      "5895625134713196",
	}
	date4 := TableDriverBuildTarjeta{
		TituloPrueba: "Validar numero de tarjeta Naranja X",
		WantTable:    true,
		Tarjeta:      "4029184698061398",
	}
	date5 := TableDriverBuildTarjeta{
		TituloPrueba: "Validar numero de tarjeta Naranja visa",
		WantTable:    true,
		Tarjeta:      "4029182269476342",
	}
	// var tableDriverTestPeyment1 []TableDriverBuildTarjeta
	tableDriverTestPeyment = append(tableDriverTestPeyment, date1, date2, date3, date4, date5)
	// tableDriverTestPeyment = append(tableDriverTestPeyment, TableBuildtarjeta{
	// 	Table: tableDriverTestPeyment1,
	// })
	return
}

// func EstructuraDiferenceuint() (tableDriverTestDiff TableDriverBuildDiferenceUint) {

// 	// date1 := TableDriverBuildTarjeta{
// 	// 	TituloPrueba: "Validar numero de tarjeta Naranja Mary",
// 	// 	WantTable:    true,
// 	// 	:      [1,2,3,4,6,],
// 	// }
// 	// var tableDriverTestPeyment1 []TableDriverBuildTarjeta
// 	tableDriverTestDiff = TableDriverBuildDiferenceUint{
// 		TituloPrueba: "Prueba",
// 		WantTable: false,
// 		String1:   [1,2,3,4,5],
// 		String2:   [1],
// 	}
// 	// tableDriverTestPeyment = append(tableDriverTestPeyment, TableBuildtarjeta{
// 	// 	Table: tableDriverTestPeyment1,
// 	// })
// 	return
// }
