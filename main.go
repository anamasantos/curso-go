package main

import (
	"emailn/internal/domain/campaign"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func main() {
	campaign := campaign.Campaign{}
	validate := validator.New()
	err := validate.Struct(campaign)
	if err == nil {
		fmt.Println("nenhum erro")
	} else {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			fmt.Println(e.StructField() + "is invalid: " + e.Tag())
		}

	}

}
