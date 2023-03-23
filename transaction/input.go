package transaction

type GetTransactionByCampaignInput struct {
	ID int `uri:"id" binding:"required"`
}
