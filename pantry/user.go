package pantry

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// RunAsUser sets syscall for the process
func RunAsUser(username string) error {
	u, err := user.Lookup(username)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	if err := syscall.Setuid(uid); err != nil {
		return err
	}

	return nil
}

// GetUserDetails returns user details
func GetUserDetails(username string) (*user.User, error) {
	var u *user.User
	u, err := user.Lookup(username)
	if err != nil {
		return u, err
	}

	return u, nil
}

// GetUIDAndGID returns the uid and gid of the passed user
func GetUIDAndGID(name string) (uint32, uint32, error) {
	var rawUID, rawGID string
	var uid, gid uint32

	if name == "self" {
		rawUID = os.Getenv("SUDO_UID")
		rawGID = os.Getenv("SUDO_UID")
	} else {
		var u *user.User
		u, err := user.Lookup(name)
		if err != nil {
			return uid, gid, err
		}

		rawUID = u.Uid
		rawGID = u.Gid
	}

	parsedUID, _ := strconv.ParseUint(rawUID, 10, 64)
	parsedGID, _ := strconv.ParseUint(rawGID, 10, 64)
	return uint32(parsedUID), uint32(parsedGID), nil
}
