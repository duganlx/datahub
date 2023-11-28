package main

import "gopkg.in/oauth2.v3"

type Client struct {
	ID     string
	Secret string
	Domain string
	UserID string
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret client domain
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetUserID user id
func (c *Client) GetUserID() string {
	return c.UserID
}

type ClientStore struct{}

func (cs *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	cm := Client{
		ID:     ClientId,
		Secret: ClientSecret,
		Domain: "",
		UserID: "",
	}

	return &cm, nil
}

// token 存储
type TokenStore struct{}

// create and store the new token information
func (ts *TokenStore) Create(info oauth2.TokenInfo) error {
	return nil
}

// delete the authorization code
func (ts *TokenStore) RemoveByCode(code string) error {
	return nil
}

// use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) error {
	return nil
}

// use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) error {
	return nil
}

// use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (oauth2.TokenInfo, error) {
	return nil, nil
}

// use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	return nil, nil
}
