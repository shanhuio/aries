package oauth

// SingleAdmin returns a user check function for a single admin.
func SingleAdmin(admin string) func(user string) (interface{}, int, error) {
	return func(user string) (interface{}, int, error) {
		if user == admin {
			return user, 10, nil
		}
		return "", 0, nil
	}
}
