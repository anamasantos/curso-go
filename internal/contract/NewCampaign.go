package contract

type NewCampaign struct {
	Name      string
	Content   string
	Emails    []string
	Status    string
	CreatedBy string
}
