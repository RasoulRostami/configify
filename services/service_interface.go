package services

type Service interface {
	Handel()
	IsRelatedMessage()
	update()
	delete()
	reload()
}
