package app

import "shs/app/models"

func (a *App) CreateMedicine(medicine models.Medicine) (models.Medicine, error) {
	return a.repo.CreateMedicine(medicine)
}

func (a *App) DeleteMedicine(id uint) error {
	return a.repo.DeleteMedicine(id)
}

func (a *App) ListAllMedicines() ([]models.Medicine, error) {
	return a.repo.ListAllMedicines()
}

func (a *App) ListMedicinesByIds(ids []uint) ([]models.Medicine, error) {
	return a.repo.ListMedicinesByIds(ids)
}

func (a *App) UpdateMedicineAmount(id uint, newAmount int) error {
	return a.repo.UpdateMedicineAmount(id, newAmount)
}

func (a *App) DecrementMedicineAmount(id uint, amount int) error {
	return a.repo.DecrementMedicineAmount(id, amount)
}

func (a *App) GetMedicine(id uint) (models.Medicine, error) {
	return a.repo.GetMedicine(id)
}
