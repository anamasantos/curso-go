package database

import (
	"emailn/internal/domain/campaign"

	"gorm.io/gorm"
)

type CampaignRepository struct {
	Db *gorm.DB
	// campaigns []campaign.Campaign
}

func (c *CampaignRepository) Create(campaign *campaign.Campaign) error {
	// c.campaigns = append(c.campaigns, *campaign)
	tx := c.Db.Create(campaign)
	return tx.Error
}

func (c *CampaignRepository) Update(campaign *campaign.Campaign) error {
	// c.campaigns = append(c.campaigns, *campaign)
	tx := c.Db.Save(campaign)
	return tx.Error
}

func (c *CampaignRepository) Get() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign
	tx := c.Db.Find(&campaigns)
	return campaigns, tx.Error
}

func (c *CampaignRepository) GetBy(id string) (*campaign.Campaign, error) {
	var campaign campaign.Campaign
	tx := c.Db.Preload("Contacts").First(&campaign, "id = ?", id)
	// if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
	// 	return nil, nil
	// }
	return &campaign, tx.Error
}

func (c *CampaignRepository) Delete(campaign *campaign.Campaign) error {
	// c.campaigns = append(c.campaigns, *campaign)
	// for i, _ := range campaign.Contacts {
	// 	c.Db.Delete(campaign.Contacts[i])
	// }

	tx := c.Db.Select("Contacts").Delete(campaign)
	return tx.Error
}
