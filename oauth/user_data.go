package oauth

import (
	"shanhu.io/aries"
)

const dataKey = "user"

// UserData fetches the user data in the context.
func UserData(c *aries.C) interface{} {
	v, ok := c.Data[dataKey]
	if !ok {
		return nil
	}
	return v
}

func setUserContext(c *aries.C, name string, u interface{}, lvl int) {
	c.User = name
	c.UserLevel = lvl
	if u != nil {
		c.Data[dataKey] = u
	}
}
