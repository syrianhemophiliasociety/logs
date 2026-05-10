package actions

import "shs/app/models"

type GetAddressesAlikeParams struct {
	ActionContext
	Address Address `json:"search_address"`
}

type GetAddressesAlikePayload struct {
	Data []Address `json:"data"`
}

func (a *Actions) GetAddressesAlike(params GetAddressesAlikeParams) (GetAddressesAlikePayload, error) {
	addresses, err := a.app.GetAllAddressesALike(models.Address{
		Governorate: params.Address.Governorate,
		Suburb:      params.Address.Suburb,
		Street:      params.Address.Street,
	})
	if err != nil {
		return GetAddressesAlikePayload{}, err
	}

	outAddresses := make([]Address, 0, len(addresses))
	for _, address := range addresses {
		outAddresses = append(outAddresses, Address{
			Governorate: address.Governorate,
			Suburb:      address.Suburb,
			Street:      address.Street,
		})
	}

	return GetAddressesAlikePayload{
		Data: outAddresses,
	}, nil
}
