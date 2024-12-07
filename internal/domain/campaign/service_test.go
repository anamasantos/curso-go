package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internal-errors"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test Y",
		Content:   "body HI!",
		Emails:    []string{"teste1@test.com"},
		CreatedBy: "teste@test.com.br",
	}
	service = campaign.ServiceImp{}
)

func Test_Create_Campaign(t *testing.T) {
	assert := assert.New(t)
	// service := Service{}
	fmt.Println(newCampaign)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(nil)
	service.Repository = repositoryMock
	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)

}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)
	newCampaign.Name = ""
	_, err := service.Create(newCampaign)
	assert.False(errors.Is(internalerrors.ErrInternal, err))

}

func Test_Create_SaveCampaign(t *testing.T) {
	// assert := assert.New(t)

	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.MatchedBy(
		func(campaign *campaign.Campaign) bool {
			if campaign.Name != newCampaign.Name ||
				campaign.Content != newCampaign.Content ||
				len(campaign.Contacts) != len(newCampaign.Emails) {
				return false
			}
			return true
		})).Return(nil)
	service.Repository = repositoryMock
	service.Create(newCampaign)
	repositoryMock.AssertExpectations(t)

}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))
	service.Repository = repositoryMock

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))

}

func Test_GetById_returnCampaign(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaign.ID
	})).Return(campaign, nil)
	service.Repository = repositoryMock
	campaignReturned, _ := service.GetBy(campaign.ID)
	fmt.Println(campaignReturned)
	assert.Equal(campaign.ID, campaignReturned.ID)
	assert.Equal(campaign.Name, campaignReturned.Name)
	assert.Equal(campaign.Status, campaignReturned.Status)
	assert.Equal(campaign.Content, campaignReturned.Content)
	assert.Equal(campaign.CreatedBy, campaignReturned.CreatedBy)

}

func Test_GetById_returnErrorWhenSomethingWrongExist(t *testing.T) {
	assert := assert.New(t)
	campaign, _ := campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	service.Repository = repositoryMock
	fmt.Println(campaign.ID)
	_, err := service.GetBy(campaign.ID)
	fmt.Println(err)
	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_returnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	campaignIDInvalid := "invalid"
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	service.Repository = repositoryMock
	err := service.Delete(campaignIDInvalid)
	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())

}

func Test_Delete_returnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
	service.Repository = repositoryMock
	err := service.Delete(campaign.ID)
	assert.Equal("Campaign status invalid", err.Error())

}

func Test_Delete_returnInternalError_when_delete_has_problem(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("test 1", "Body !!", []string{"emai1@gmail.com"}, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(errors.New("error to delete campaign"))
	service.Repository = repositoryMock
	err := service.Delete(campaignFound.ID)
	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Delete_returnINil_when_delete_has_success(t *testing.T) {
	assert := assert.New(t)
	campaignFound, _ := campaign.NewCampaign("test 1", "Body !!", []string{"emai1@gmail.com"}, newCampaign.CreatedBy)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignFound, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignFound == campaign
	})).Return(nil)
	service.Repository = repositoryMock
	err := service.Delete(campaignFound.ID)
	assert.Nil(err)

}

func Test_Start_returnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	assert := assert.New(t)
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	campaignIDInvalid := "invalid"
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	service.Repository = repositoryMock

	err := service.Start(campaignIDInvalid)
	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())

}

func Test_Start_returnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	assert := assert.New(t)
	campaign := &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
	service.Repository = repositoryMock
	err := service.Start(campaign.ID)
	assert.Equal("Campaign status invalid", err.Error())

}

func Test_Start_should_send_mail(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Peding}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	service.Repository = repositoryMock
	emailWasSend := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignSaved.ID {
			emailWasSend = true
		}

		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignSaved.ID)
	assert.True(emailWasSend)

}

func Test_Start_ReturnError_when_func_SencMail_fail(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Peding}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	service.Repository = repositoryMock
	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(campaignSaved.ID)
	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Start_ReturnNil_when_updated_to_done(t *testing.T) {
	assert := assert.New(t)
	campaignSaved := &campaign.Campaign{ID: "1", Status: campaign.Peding}
	repositoryMock := new(internalmock.CampaignRepositoryMock)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignSaved, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignSaved.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	service.Repository = repositoryMock
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignSaved.ID)
	assert.Equal(campaign.Done, campaignSaved.Status)

}
