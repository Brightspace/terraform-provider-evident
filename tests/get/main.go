package main

import (
	"github.com/Brightspace/terraform-provider-evident/evident/api"
	"os"
	"log"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing ID for evident resource")
	}
	arg := os.Args[1]

	
	client := api.Evident{
		Credentials: api.Credentials{
			AccessKey: []byte(os.Getenv("EVIDENT_ACCESS_KEY")),
			SecretKey: []byte(os.Getenv("EVIDENT_SECRET_KEY")),
		},
		RetryMaximum: 5,
	}
	
	result, _ := client.Get(arg)
    fmt.Println("id:\n", result.ID)
	fmt.Println("name:\n", result.Attributes.Name)
	fmt.Println("arn:\n", result.Attributes.Arn)
	fmt.Println("external_id:\n", result.Attributes.ExternalID)
}