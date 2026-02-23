package usecase

func (masterdata masterdataUsecase) Health() (string, error) {
	err := masterdata.postgreRepository.Health()
	if err != nil {
		return "", err
	}

	return "Service is healthy", nil
}
