package transaction

type Service interface {
	GetTransactionByCampaignID(input GetTransactionByCampaignInput) ([]Transaction, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionByCampaignID(input GetTransactionByCampaignInput) ([]Transaction, error) {
	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
