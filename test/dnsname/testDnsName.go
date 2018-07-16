package main

/*
https://github.com/hashicorp/go-immutable-radix

other choices:
https://github.com/tchap/go-patricia
https://github.com/armon/go-radix

 */
import(
	"fmt"
	"RsDnsTools/util"
)

func main() {
	/*
	r := iradix.New()

	r, _, _ = r.Insert([]byte(reverseString("sina.com")), 1)

	r, _, _ = r.Insert([]byte(reverseString("163.com")), 2)

	r, _, _ = r.Insert([]byte(reverseString("www.sina.com")), 3)


	rkey, rvalue, rresult := r.Root().LongestPrefix([]byte(reverseString("sports.www.sina.com")))

	fmt.Printf("match string: %s, return value is %d, return result is %v\r\n", reverseString(string(rkey)), rvalue, rresult)
*/
	//s := reverseString("www.sina.com")
	//fmt.Println(s)
	dnamelist := util.NewDnameList()
	dnamelist.Insert("sina.com")
	dnamelist.Insert("163.com")
	dnamelist.Insert("www.sina.com")

	testnames := []string{"sports.sina.com", "ea.com", "www.taobao.com"}
	for _, testname := range testnames{
		fmt.Printf("test name %s match result is %v\r\n", testname, dnamelist.Match(testname))
	}


}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}