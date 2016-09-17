package passwd

import (
	"bytes"
	"errors"
	"math/rand"
	"os"
	"sort"
	"time"
)

var r rand.Source

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type PasswdFile struct {
	location  string
	cipherkey []byte
}

func New(loc string, key ...[]byte) *PasswdFile {
	if len(key) == 0 {

		return &PasswdFile{location: loc}
	}

	return &PasswdFile{location: loc, cipherkey: key[0]}
}

// List returns a []string list of user identifiers
func List() []string {
	var list []string
	for user, _ := range usertable {
		list = append(list, user)
	}
	sort.Strings(list)

	return list
}

var lock bool

// Write actually writes the user names to the passwd file after truncating to 0
// Before calling Write(), the usertable is only in memory.
// Starting your code with a "defer passwd.Write()" will hopefully store any changes in the event of a panic
// Only call once at a time.
func Write() error {
	if lock == true {
		return errors.New("Error: Passwd file locked")

	}
	lock = true
	file, err := os.OpenFile(pass.location, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	for user, password := range usertable {

		_, err := file.WriteString(user + ":" + string(password) + "\n")
		if err != nil {
			return err
		}
		//	fmt.Printf("wrote %v bytes\n", i)
	}
	lock = false
	return nil

}

// Parse stores the usertable to memory to enable Match()
func Parse() {
	b, err := open()
	if err != nil {
		panic(err)
	}
	splits := bytes.Split(b, []byte("\n"))
	for i, line := range splits {
		i++
		lines := bytes.Split(line, []byte(":"))
		if len(lines) == 2 {

			usertable[string(lines[0])] = lines[1]
			//	log.Println("Parsed", lines[0], lines[1])
		} else if i != len(splits) { // last line
			//if len(lines) != 1 {
			//	log.Println("skipping line", i)
			//	}
		}
	}
}

// InsertOrUpdate returns no error.
func InsertOrUpdate(user string, password []byte) {
	usertable[user] = passwd2hash(password)
}

// Insert returns an error if the user exists already.
// Otherwise, it inserts the user and hashed password into the usertable.
// Write() must be called to make actual changes to the passwd file.
func Insert(user string, password []byte) error {
	if usertable[user] == nil {
		usertable[user] = passwd2hash(password)
	} else {
		return errors.New("User \"" + user + "\" exists")
	}
	return nil
}

// Delete removes a user from the usertable.
func Delete(user string) error {
	delete(usertable, user)
	return nil // dummy error for now?
}

// Update updates a user password.
func Update(user string, password []byte) error {

	if usertable == nil {
		return errors.New("No usertable.")
	}
	usertable[user] = passwd2hash(password)
	return nil
}

// Update updates a id/username from old to new, and if p != nil updates the password.
func UpdateID(old, new string, p []byte) error {

	if usertable == nil {
		return errors.New("UpdateID: No usertable.")
	}

	if new == "" {
		return errors.New("UpdateID: Can't change name to an empty string.")
	}

	if usertable[new] != nil {
		return errors.New("UpdateID: User \"" + new + "\" exists")
	}

	usertable[new] = usertable[old] // copy user
	delete(usertable, old)          // delete old
	if p != nil {
		usertable[new] = passwd2hash(p) // update password
	}

	return nil
}

// Match, for login sequences. Must be called after Parse(). Read-only.
func Match(user string, password []byte) bool {

	for realuser, realpass := range usertable {
		if string(realuser) == string(user) {
			if compareHash(realpass, password) {
				return true
			}
		}
	}
	return false
}
