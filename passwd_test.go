package passwd_test

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aerth/passwd"
)

var a []int
var c chan int
var cb chan bool
var c1, c2, c3, c4 chan int
var i1, i2 int
var b bool
var err error
var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func dummy() {
	user := strconv.FormatUint(uint64(r.Uint32()), 10)
	user2 := strconv.FormatUint(uint64(r.Uint32()), 10)
	err = passwd.Insert(user+user2, []byte("password"))
	if err != nil {
		fmt.Println("Dummy Insert Error: ", err)
		log.Println("Created Dummy:", user+user2)
	} else {
		e := passwd.Write()
		if e != nil {
			fmt.Println("Didn't create dummy:", e)
		} else {
			log.Println("Created Dummy:", user+user2)
		}
	}

}
func Example() {

	passwd.Parse()

	dummy()
	dummy()
	dummy()

	b = passwd.Match("root", []byte("pasksword"))
	fmt.Println("Match:", b)
	b = passwd.Match("root", []byte("password"))
	fmt.Println("Match:", b)

	err = passwd.Update("root", []byte("pasksword"))
	if err != nil {
		fmt.Println("Insert Error: ", err)
	} else {
		passwd.Write()
	}
	b = passwd.Match("root", []byte("pasksword"))
	fmt.Println("Match:", b)
	b = passwd.Match("root", []byte("password"))
	fmt.Println("Match:", b)

	// err = passwd.Delete("root")
	// if err != nil {
	// 	fmt.Println("Delete Error: ", err)
	// } else {
	// 	passwd.Write()
	// }
	// b = passwd.Match("root", []byte("pasksword"))
	// fmt.Println("Match:", b)
	// b = passwd.Match("root", []byte("password"))
	// fmt.Println("Match:", b)

	for {
	}

}

func sleep(n int64) {
	time.Sleep(time.Second * 1)
}
