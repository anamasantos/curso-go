package campaign

import (
	"emailn/internal/contract"
	internalerrors "emailn/internal/internal-errors"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Service interface {
	Create(newCampaign contract.NewCampaign) (string, error)
	GetBy(id string) (*contract.CampaignResponse, error)
	Cancel(id string) error
	Delete(id string) error
}

type ServiceImp struct {
	Repository Repository
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {

	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails)
	if err != nil {
		return "", err
	}
	err = s.Repository.Create(campaign)
	if err != nil {
		return "", internalerrors.ErrInternal
	}
	return campaign.ID, nil
}

func (s *ServiceImp) GetBy(id string) (*contract.CampaignResponse, error) {
	fmt.Println("SERVICE:" + id)
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internalerrors.ErrInternal
		}
		return nil, internalerrors.ProcessErrorToReturn(err)
	}
	// if campaign == nil {
	// 	return nil, nil
	// }
	return &contract.CampaignResponse{
		ID:                   campaign.ID,
		Content:              campaign.Content,
		Name:                 campaign.Name,
		Status:               campaign.Status,
		AmountOfEmailsToSend: len(campaign.Contacts),
	}, nil

}

func (s *ServiceImp) Cancel(id string) error {
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign.Status != Peding {
		return errors.New("Campaign status invalid")
	}
	campaign.Cancel()
	err = s.Repository.Update(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}

	return nil

}

func (s *ServiceImp) Delete(id string) error {
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign.Status != Peding {
		return errors.New("Campaign status invalid")
	}
	campaign.Delete()
	err = s.Repository.Delete(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}

	return nil

}
