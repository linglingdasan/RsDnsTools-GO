package util

import "github.com/hashicorp/go-immutable-radix"


type dnamelist struct{
	dnametree *iradix.Tree
}

func NewDnameList() *dnamelist{
	dnl := &dnamelist{}
	dnl.dnametree = iradix.New()
	return dnl
}

func (dnl *dnamelist)Insert(namestr string){
	dnl.dnametree, _, _ = dnl.dnametree.Insert([]byte(reverseString(namestr)), true)
}

func (dnl *dnamelist)Match(namestr string)bool{
	_, _, result := dnl.dnametree.Root().LongestPrefix([]byte(reverseString(namestr)))
	return result
}


func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}