package token

import (
	"context"
	"fmt"
	"testing"
)

func TestTokenGenerate(t *testing.T) {
	ctx := context.Background()
	tokenService := &TokenService{
		UserId: 1,
	}
	data := map[string]interface{}{
		"user_id": 1,
	}
	resp, err := tokenService.Generate(ctx, data)
	if err != nil {
		t.Error(err)
	} else {
		println(resp)
	}
}

func TestTokenVerify(t *testing.T) {
	ctx := context.Background()
	tokenService := &TokenService{}
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3OTE1NjIzMTcsInNpZ24iOiJiMTdkMmQyY2U0NDRmOTViNTgyODE3MTA1ZjNiYWYyNSIsInVzZXJfaWQiOjF9.TVyx1a3A9vqf-2nkt73voXBor314knElfeyRPAOuFhs"
	resp, err := tokenService.Verify(ctx, tokenStr)
	if err != nil {
		t.Error(err)
	} else {
		println(fmt.Sprintf("%+v", resp))
	}
}
