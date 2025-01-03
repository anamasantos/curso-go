package campaign_test

import (
	"emailn/internal/contract"
	"emailn/internal/domain/campaign"
	internalerrors "emailn/internal/internal-errors"
	internalmock "emailn/internal/test/internal-mock"
	"errors"
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
	campaignPedenting *campaign.Campaign
	campaignStarted   *campaign.Campaign
	repositoryMock    *internalmock.CampaignRepositoryMock
	service           = campaign.ServiceImp{}
)

func setUp() {
	campaignPedenting, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
}

func Test_Create_Campaign(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("Create", mock.Anything).Return(nil)
	id, err := service.Create(newCampaign)

	assert.NotNil(id)
	assert.Nil(err)

}

func Test_Create_ValidateDomainError(t *testing.T) {
	assert := assert.New(t)
	_, err := service.Create(contract.NewCampaign{})
	assert.False(errors.Is(internalerrors.ErrInternal, err))

}

func Test_Create_SaveCampaign(t *testing.T) {
	setUp()
	repositoryMock.On("Create", mock.MatchedBy(
		func(campaign *campaign.Campaign) bool {
			if campaign.Name != newCampaign.Name ||
				campaign.Content != newCampaign.Content ||
				len(campaign.Contacts) != len(newCampaign.Emails) {
				return false
			}
			return true
		})).Return(nil)
	service.Create(newCampaign)
	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("Create", mock.Anything).Return(errors.New("error to save on database"))

	_, err := service.Create(newCampaign)

	assert.True(errors.Is(internalerrors.ErrInternal, err))

}

func Test_GetById_returnCampaign(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPedenting.ID
	})).Return(campaignPedenting, nil)
	campaignReturned, _ := service.GetBy(campaignPedenting.ID)
	assert.Equal(campaignPedenting.ID, campaignReturned.ID)
	assert.Equal(campaignPedenting.Name, campaignReturned.Name)
	assert.Equal(campaignPedenting.Status, campaignReturned.Status)
	assert.Equal(campaignPedenting.Content, campaignReturned.Content)
	assert.Equal(campaignPedenting.CreatedBy, campaignReturned.CreatedBy)

}

func Test_GetById_returnErrorWhenSomethingWrongExist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))
	_, err := service.GetBy(campaignPedenting.ID)
	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Delete_returnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaignIDInvalid := "invalid"
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	err := service.Delete(campaignIDInvalid)
	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())

}

func Test_Delete_returnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)
	err := service.Delete(campaignStarted.ID)
	assert.Equal("Campaign status invalid", err.Error())

}

func Test_Delete_returnInternalError_when_delete_has_problem(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPedenting, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))
	service.Repository = repositoryMock
	err := service.Delete(campaignPedenting.ID)
	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Delete_returnINil_when_delete_has_success(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPedenting, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPedenting == campaign
	})).Return(nil)
	err := service.Delete(campaignPedenting.ID)
	assert.Nil(err)

}

func Test_Start_returnRecordNotFound_when_campaign_does_not_exist(t *testing.T) {
	setUp()
	assert := assert.New(t)
	campaignIDInvalid := "invalid"
	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
	err := service.Start(campaignIDInvalid)
	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())

}

func Test_Start_returnStatusInvalid_when_campaign_has_status_not_equals_pending(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)
	err := service.Start(campaignStarted.ID)
	assert.Equal("Campaign status invalid", err.Error())

}

func Test_Start_should_send_mail(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPedenting, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)
	emailWasSend := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignPedenting.ID {
			emailWasSend = true
		}

		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPedenting.ID)
	assert.True(emailWasSend)

}

func Test_Start_ReturnError_when_func_SencMail_fail(t *testing.T) {
	setUp()
	assert := assert.New(t)
	repositoryMock.On("GetBy", mock.Anything).Return(campaignPedenting, nil)
	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(campaignPedenting.ID)
	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Start_ReturnNil_when_updated_to_done(t *testing.T) {
	setUp()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignPedenting, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPedenting.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPedenting.ID)
	assert.Equal(campaign.Done, campaignPedenting.Status)

}
