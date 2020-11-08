package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"

	"learn.auth.billing/model"
)

func TestClaimDecoder(t *testing.T) {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJWemcyZ0FrRTBUaGFpTFpBakVZU0FMLWczWW5TNVdmUGkzVDZVcXE4MlNNIn0.eyJleHAiOjE2MDM1NjYwODQsImlhdCI6MTYwMzU2NTc4NCwiYXV0aF90aW1lIjoxNjAzNTY1Nzc3LCJqdGkiOiJmNjZmZGVkMi02YjIwLTQ0NDItOWI0My02MjFkN2RlMjI3MDMiLCJpc3MiOiJodHRwOi8vMTkyLjE2OC4yLjEwOjgwODAvYXV0aC9yZWFsbXMvbGVhcm5pbmdBcHAiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiYTZkN2M4NjEtMGE2NS00MDEzLTk4OTQtNmJhMjEwMjAzZmE4IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiYmlsbGluZ0FwcCIsInNlc3Npb25fc3RhdGUiOiJhNjI5ZDMyYS1mZTM1LTQyMDMtOWMyNy0zODYxMzUzZjI0ODIiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly9sb2NhbGhvc3Q6ODA4MSJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsib2ZmbGluZV9hY2Nlc3MiLCJ1bWFfYXV0aG9yaXphdGlvbiJdfSwicmVzb3VyY2VfYWNjZXNzIjp7ImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwibmFtZSI6ImJvYiBib2IiLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib2IiLCJnaXZlbl9uYW1lIjoiYm9iIiwiZmFtaWx5X25hbWUiOiJib2IiLCJlbWFpbCI6ImJvYkBib2IuY29tIn0.QAmN9TOHbgJKV5IaMq_BGFhJF63-313UQHQ5AQ2aFv51cmTIUDwQBfsup-ZujfnkEHslDjBviHZywLbobiYk3s6lGYPgEAR8pgEbb6K20WJHIZ4H3s53POTtI4p7nk4dUIG2RjS_VDnTJ99oi0NToMhVTfCYKwsVyBMSyX2rHMYVcIsSkGY94fU-F4Dg7abHd6gMRAujCNEjyEN5l32OiDab8itztGKbo_7Tw7sJGBvwu06Ok3EoZj6i-6IUI-asH8SHuEUjyh0M9pG1dncN1jVj4SSmE0hCHPoosBR1YZFL2WAOhI8SSjCrpjJVVffPK3SsMelbrnOhbLPG9yomSw"

	tokenParts := strings.Split(tokenString, ".")
	claims, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	log.Println("tokenParts[1]\n", tokenParts[1])
	if err != nil {
		t.Error(err)
	}
	log.Println("Claim: ", string(claims))
}

func TestAudString(t *testing.T) {
	fmt.Println("")
	token := `
	{
		"aud": "billingServiceV2"
	  }
	`
	claim := &model.Tokenclaim{}
	json.Unmarshal([]byte(token), claim)
	values := claim.AudAsSlice()
	length := len(values)
	if length != 1 {
		t.Errorf("Expected 1 element in slice. Got: %v", length)
	}

	if values[0] != "billingServiceV2" {
		t.Errorf("Expected value billingServiceV2. Got: %v", values[0])
	}
}

func TestAudSlice(t *testing.T) {
	fmt.Println("")
	token := `
	{
		"aud": ["billingService", "billingServiceV2"]
	  }
	`
	claim := &model.Tokenclaim{}
	json.Unmarshal([]byte(token), claim)
	values := claim.AudAsSlice()
	length := len(values)
	if length != 2 {
		t.Errorf("Expected 1 element in slice. Got: %v", length)
	}

	if values[0] != "billingService" {
		t.Errorf("Expected value billingService. Got: %v", values[0])
	}

	if values[1] != "billingServiceV2" {
		t.Errorf("Expected value billingServiceV2. Got: %v", values[1])
	}
}