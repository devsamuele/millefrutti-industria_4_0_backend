package model

import "regexp"

func CheckEmail(email string) (bool, error) {
	matched, err := regexp.MatchString(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, email)
	if err != nil {
		return false, err
	}

	return matched, nil
}
