package campaign

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

var (
	name     = "CampanhaX"
	content  = "Bodyteste"
	contacts = []string{"email1@e.com", "email2@e.com"}
	createBy = "teste@teste.com.br"
	fake     = faker.New()
)

func Test_NewCampaign_CreateCampaign(t *testing.T) {

	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createBy)
	assert.Equal(campaign.Name, name)
	assert.Equal(campaign.Content, content)
	assert.Equal(len(campaign.Contacts), len(contacts))
	assert.Equal(createBy, campaign.CreatedBy)

}

func Test_NewCampaign_IDisNotNill(t *testing.T) {

	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createBy)
	assert.NotNil(campaign.ID)

}

func Test_NewCampaign_CreatedOnMustBeNow(t *testing.T) {

	assert := assert.New(t)
	now := time.Now().Add(-time.Minute)
	campaign, _ := NewCampaign(name, content, contacts, createBy)
	assert.Greater(campaign.CreatedOn, now)

}

func Test_NewCampaign_MustValidateNameMin(t *testing.T) {

	assert := assert.New(t)
	_, err := NewCampaign("", content, contacts, createBy)
	assert.Equal("name is required with min 5", err.Error())

}

func Test_NewCampaign_MustStatusStartWithPending(t *testing.T) {

	assert := assert.New(t)
	campaign, _ := NewCampaign(name, content, contacts, createBy)
	assert.Equal(Peding, campaign.Status)

}

func Test_NewCampaign_MustValidateNameMax(t *testing.T) {

	assert := assert.New(t)

	_, err := NewCampaign(fake.Lorem().Text(28), content, contacts, createBy)
	assert.Equal("name is required with max 24", err.Error())

}

func Test_NewCampaign_MustValidateContentMin(t *testing.T) {

	assert := assert.New(t)
	_, err := NewCampaign(name, "", contacts, createBy)
	assert.Equal("content is required with min 5", err.Error())

}

func Test_NewCampaign_MustValidateContentMax(t *testing.T) {

	assert := assert.New(t)
	_, err := NewCampaign(name, fake.Lorem().Text(1040), contacts, createBy)
	assert.Equal("content is required with max 1024", err.Error())

}

func Test_NewCampaign_MustValidateContactsMin(t *testing.T) {

	assert := assert.New(t)
	_, err := NewCampaign(name, content, nil, createBy)
	assert.Equal("contacts is required with min 1", err.Error())

}

func Test_NewCampaign_MustValidateContacts(t *testing.T) {

	assert := assert.New(t)
	_, err := NewCampaign(name, content, []string{"email_invalid"}, createBy)
	assert.Equal("email is invalid", err.Error())

}

func Test_NewCampaign_MustValidateCreatedBy(t *testing.T) {

	assert := assert.New(t)
	_, err := NewCampaign(name, content, contacts, "")
	assert.Equal("createdby is invalid", err.Error())

}
