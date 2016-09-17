package passwd

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type credentials struct {
	user     string
	password []byte
}

var pass = New("")
var usertable = map[string][]byte{} // map[user] = password
//var Options options
//var Credentials credentials

func SetLocation(loc string) {
	pass.location = loc
}
func init() {
	if pass.location == "" {
		pass.location = "default.passwd.txt"
	}
}

func createIfNotExists(loc string) {
	_, err := os.Open(loc)
	if err != nil { // Some kind of error
		// Doesn't exist
		if err.Error() == "open "+loc+": no such file or directory" {
			_, err = os.Create(loc)
			if err != nil {

			}
			// Just was created
			return
		}

	}
	// Exists
	return
}
func open() ([]byte, error) {
	createIfNotExists(pass.location)
	file, err := os.Open(pass.location)
	if err != nil {

	}
	fileinfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return nil, err

	}
	var b = make([]byte, fileinfo.Size()*2)
	file.Read(b)
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

// passwd2hash returns a bcrypt hash
func passwd2hash(password []byte) []byte {
	key, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil
	}
	return key
}

// compareHash login
func compareHash(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err == nil {
		return true
	}
	return false

}
