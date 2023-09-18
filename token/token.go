package token

import "terraform-provider-zoom/client"

func GenerateToken(c *client.Client, accountId, clientId, clientSecret string) error {
	return c.GenerateToken(accountId, clientId, clientSecret)
}
