package endpoints

import (
	"emailn/internal/contract"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	internalmock "emailn/internal/test/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CampaignGetById_should_return_campaign(t *testing.T) {
	assert := assert.New(t)
	campaign := contract.CampaignResponse{
		ID:      "343",
		Name:    "Test",
		Content: "Hi everyone",
		Status:  "Pending",
	}
	fmt.Println(campaign)
	service := new(internalmock.CampaignServiceMock)
	service.On("GetBy", mock.Anything).Return(&campaign, nil)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	response, status, _ := handler.CampaignGetById(res, req)
	fmt.Println(campaign.ID)
	assert.Equal(200, status)
	assert.Equal(campaign.ID, response.(*contract.CampaignResponse).ID)
	assert.Equal(campaign.Name, response.(*contract.CampaignResponse).Name)

}

func Test_CampaignGetById_should_return_error_when_something_wrong(t *testing.T) {
	assert := assert.New(t)
	service := new(internalmock.CampaignServiceMock)
	errExpected := errors.New("something whong")
	service.On("GetBy", mock.Anything).Return(nil, errExpected)
	handler := Handler{CampaignService: service}
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	_, _, errReturned := handler.CampaignGetById(res, req)
	assert.Equal(errExpected.Error(), errReturned.Error())

}
