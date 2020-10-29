package util

import (
	"errors"
	"github.com/hashicorp/go-immutable-radix"
)


type dnamelist struct{
	dnametree *iradix.Tree
}

func NewDnameList() *dnamelist{
	dnl := &dnamelist{}
	dnl.dnametree = iradix.New()
	return dnl
}

func (dnl *dnamelist)Insert(namestr string, groupid int){
	dnl.dnametree, _, _ = dnl.dnametree.Insert([]byte(reverseString(namestr)), groupid)
}

func (dnl *dnamelist)Match(namestr string)bool{
	_, _, result := dnl.dnametree.Root().LongestPrefix([]byte(reverseString(namestr)))
	return result
}

func (dnl *dnamelist)GetId(namestr string)int{
	_, val, _ := dnl.dnametree.Root().LongestPrefix([]byte(reverseString(namestr)))

	var groupid int
	if val != nil{
		groupid = val.(int)
	}else {
		groupid = -1
	}

	return groupid
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

type dnametree struct{
	dnametree *iradix.Tree
}
func NewDnameTree() *dnametree{
	dnt := &dnametree{}
	dnt.dnametree = iradix.New()
	return dnt
}
func (dnt *dnametree)Insert(namestr string, pInterface interface{}){
	dnt.dnametree, _, _ = dnt.dnametree.Insert([]byte(reverseString(namestr)), pInterface)
}
func (dnt *dnametree)Match(namestr string)bool{
	_, _, result := dnt.dnametree.Root().LongestPrefix([]byte(reverseString(namestr)))
	return result
}
func (dnt *dnametree)GetPInterface(namestr string) (pInterface interface{}, err error){
	_, val, _ := dnt.dnametree.Root().LongestPrefix([]byte(reverseString(namestr)))
	if val != nil{
		return val, nil
	}else {
		return val, errors.New("no this dname")
	}
}