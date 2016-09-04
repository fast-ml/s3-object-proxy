package lambda

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/apex/go-apex"
	"github.com/sheeley/s3-object-proxy/permission"
	"github.com/sheeley/s3-object-proxy/s3proxy"
)

type imageRequest struct {
	BodyJSON struct{} `json:"body-json"`
	Context  struct {
		Account_id                      string `json:"account-id"`
		API_id                          string `json:"api-id"`
		API_key                         string `json:"api-key"`
		Authorizer_principal_id         string `json:"authorizer-principal-id"`
		Caller                          string `json:"caller"`
		Cognito_authentication_provider string `json:"cognito-authentication-provider"`
		Cognito_authentication_type     string `json:"cognito-authentication-type"`
		Cognito_identity_id             string `json:"cognito-identity-id"`
		Cognito_identity_pool_id        string `json:"cognito-identity-pool-id"`
		HTTP_method                     string `json:"http-method"`
		Request_id                      string `json:"request-id"`
		Resource_id                     string `json:"resource-id"`
		Resource_path                   string `json:"resource-path"`
		Source_ip                       string `json:"source-ip"`
		Stage                           string `json:"stage"`
		User                            string `json:"user"`
		User_agent                      string `json:"user-agent"`
		User_arn                        string `json:"user-arn"`
	} `json:"context"`
	Params struct {
		Header      struct{} `json:"header"`
		Path        struct{} `json:"path"`
		Querystring struct {
			Key string `json:"o"`
		} `json:"querystring"`
	} `json:"params"`
	StageVariables struct{} `json:"stage-variables"`
}

func DefaultRedisHandleFunc(bucket, region, redisAddr string) {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var ir imageRequest
		if err := json.Unmarshal(event, &ir); err != nil {
			panic(err)
		}

		key := ir.Params.Querystring.Key
		if key == "" {
			panic(errors.New("key is empty"))
		}

		log.Println("handling", key)

		proxy, err := s3proxy.New(&s3proxy.Config{
			Bucket:           bucket,
			Region:           region,
			PermissionLookup: permission.NewRedisLookup(redisAddr),
		})

		if err != nil {
			panic(err)
		}
		return proxy.GetObject(key)
	})
}
