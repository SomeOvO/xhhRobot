package xhh

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 混淆函数开始
func vm(num int) int {
	if num&128 != 0 {
		return int(255 & ((uint16(num) << 1) ^ 27))
	}
	return num << 1
}
func qm(num int) int {
	return vm(num) ^ num
}
func _m(num int) int {
	return qm(vm(num))
}
func ym(num int) int {
	return _m(qm(vm(num)))
}
func gm(num int) int {
	return ym(num) ^ _m(num) ^ qm(num)
}

func mixed(e []int) [6]int {
	t := [6]int{}
	t[0] = gm(e[0]) ^ ym(e[1]) ^ _m(e[2]) ^ qm(e[3])
	t[1] = qm(e[0]) ^ gm(e[1]) ^ ym(e[2]) ^ _m(e[3])
	t[2] = _m(e[0]) ^ qm(e[1]) ^ gm(e[2]) ^ ym(e[3])
	t[3] = ym(e[0]) ^ _m(e[1]) ^ qm(e[2]) ^ gm(e[3])
	t[4] = e[4]
	t[5] = e[5]
	return t
}

// 混淆函数结束
func getNonce(Time int64) (string, error) {

	random, err := rand.Int(rand.Reader, big.NewInt(time.Now().UnixMilli()))
	if err != nil {
		return "", err
	}
	str := strconv.Itoa(int(Time)) + strconv.Itoa(int(random.Int64()))
	_md5 := md5.Sum([]byte(str))
	return strings.ToUpper(hex.EncodeToString(_md5[:])), nil
}

func GetKeys(reqpath string) (hkey string, nonce string, Rtime int, err error) {
	_time := time.Now().Unix()
	nonce, err = getNonce(_time)
	if err != nil {
		return
	}
	r := "AB45STUVWZEFGJ6CH01D237IXYPQRKLMN89"
	str1 := av(strconv.Itoa(int(_time)), r, -2)
	str2 := sv(reqpath, r)
	str3 := sv(nonce, r)
	var strArr = [3]string{str1, str2, str3}
	sort.Slice(strArr[:], func(i, j int) bool {
		return len(strArr[i]) < len(strArr[j])
	})
	NewString := newStr(strArr[:])
	_md5 := md5.Sum([]byte(NewString)[0:20])
	Strmd5 := hex.EncodeToString(_md5[:])
	lastsix := Strmd5[len(Strmd5)-6:]
	var lastsixArr [6]int
	for i, v := range lastsix {
		lastsixArr[i] = int(v)
	}
	mix := mixed(lastsixArr[:])
	var count int
	for _, v := range mix {
		count += v
	}
	a := fmt.Sprintf("%02d", count%100)
	s := av(Strmd5[0:5], r, -4)
	return "" + s + a, nonce, int(_time), nil
}

func newStr(arr []string) string {
	var str strings.Builder
	for i := range arr[2] {
		if len(arr[0]) > i {
			str.WriteString(string(arr[0][i]))
		}
		if len(arr[1]) > i {
			str.WriteString(string(arr[1][i]))
		}
		if len(arr[2]) > i {
			str.WriteString(string(arr[2][i]))
		}
	}
	return str.String()
}

func av(str string, key string, n int) string {
	var r strings.Builder
	i := key[0 : len(key)+n]
	for _, v := range str {
		p := i[int(v)%len(i)]
		r.WriteString(string(p))
	}
	return r.String()
}
func sv(str, key string) string {
	var n strings.Builder
	for _, v := range str {
		p := key[int(v)%len(key)]
		n.WriteString(string(p))
	}
	return n.String()
}
