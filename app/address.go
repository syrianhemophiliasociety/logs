package app

import "shs/app/models"

func (a *App) CreateAddress(address models.Address) (models.Address, error) {
	return a.repo.CreateAddress(address)
}

func (a *App) GetAllAddresses() ([]models.Address, error) {
	return a.repo.GetAllAddresses()
}

func (a *App) GetAllAddressesALike(searchAddress models.Address) ([]models.Address, error) {
	return a.repo.GetAllAddressesALike(searchAddress)
}

func (a *App) DeleteAddress(id uint) error {
	return a.repo.DeleteAddress(id)
}
