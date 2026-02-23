package masterdata

type Usecase interface {
	//health check usecase
	Health() (string, error)
}
