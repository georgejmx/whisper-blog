package security

import (
	d "whisper-blog/controller"
)

/* Function to validate the provided hash against the chain law, determining
whether a lawful post can be made */
// TODO: Encode the proper chain law e.g after 6 days an extra hash is valid
func ValidateHash(dbo d.DatabaseObject, hash string) (bool, error) {
	storedHash, err := dbo.SelectHash()
	if err != nil {
		return false, err
	} else if storedHash != hash {
		return false, nil
	}

	// There's no error and hashes are equal
	return true, nil
}
