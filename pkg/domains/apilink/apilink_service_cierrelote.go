package apilink

/*
Crea un nuevo debin
requerimentoId es un identificador único para la operación.
*/

func (s *aplinkService) PutApilinkCierrelote(listaDebinesId []string) (erro error) {
	return s.repository.PutApilinkCierrelote(listaDebinesId)
}
