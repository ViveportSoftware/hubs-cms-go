package service

import (
	"context"
	"hubs-cms-go/config"
	"hubs-cms-go/logger"

	"github.com/machinebox/graphql"
)

//
// It uses graphQL lib not using resty
// , and force refresh token everytime
//
func SendDirectusGraphQLCmd(mutations string) (result map[string]interface{}, err error) {

	directusAccessToken := ""
	directusAccessToken, err = GetDirectusAccessToken(true)
	if err != nil {
		logger.Error.Printf("[SendDirectusGraphQLCmd] unable to get directus access token error: %v\n", err)
		return
	}

	// create graphql client
	client := graphql.NewClient(config.GetDirectusGraphQLURI())
	client.Log = func(msg string) {
		logger.Debug.Printf("[SendDirectusGraphQLCmd] %+v", msg)
	}

	// make a request
	req := graphql.NewRequest(mutations)
	req.Header.Set("Authorization", directusAccessToken)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	// result map[string]interface{} or map[string]map[string]float64
	if err = client.Run(ctx, req, &result); err != nil {
		logger.Error.Printf("[SendDirectusGraphQLCmd] %+v", err)
	}
	return
}
