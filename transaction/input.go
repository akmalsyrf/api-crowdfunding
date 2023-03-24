package transaction

import "api-crowdfunding/user"

type GetTransactionByCampaignInput struct {
	ID   int `uri:"id" binding:"required"`
	User user.User
}
