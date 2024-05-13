package Lendme

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// //////////////////////////////////////////////////////////////////////////////////////////////////
// Functions to generate OTP PINs and authentication
// //////////////////////////////////////////////////////////////////////////////////////////////////
const (
	letterAndDigitsBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 52 possibilities
	digitsBytes          = "0123456789"
	letterIdxBits        = 6                    // 6 bits to represent 64 possibilities / indexes
	letterIdxMask        = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func GenerateAPIAuthenticationKey(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = SecureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterAndDigitsBytes) {
			result[i] = letterAndDigitsBytes[idx]
			i++
		}
	}
	return string(result)
}

func GenerateOTP(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = SecureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(digitsBytes) {
			result[i] = digitsBytes[idx]
			i++
		}
	}
	return string(result)
}

// SecureRandomBytes returns the requested number of bytes using crypto/rand
func SecureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

func GenerateRandomString(length int) string {
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = SecureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterAndDigitsBytes) {
			result[i] = letterAndDigitsBytes[idx]
			i++
		}
	}
	return string(result)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////

func Ordinal(x int) string {
	suffix := "th"
	switch x % 10 {
	case 1:
		if x%100 != 11 {
			suffix = "st"
		}
	case 2:
		if x%100 != 12 {
			suffix = "nd"
		}
	case 3:
		if x%100 != 13 {
			suffix = "rd"
		}
	}
	return strconv.Itoa(x) + suffix
}

func GetTimeStamp() (YYYYMMDD string, HHMISS string) {
	CurrentDateTime := time.Now()
	yyyy, mm, dd := CurrentDateTime.Date()
	hr, mi, ss := CurrentDateTime.Clock()

	str_yyyy := strconv.Itoa(yyyy)
	str_mm := strconv.Itoa(int(mm))
	if mm < 10 {
		str_mm = "0" + str_mm
	}
	str_dd := strconv.Itoa(dd)
	if dd < 10 {
		str_dd = "0" + str_dd
	}
	YYYYMMDD = str_yyyy + str_mm + str_dd

	str_hr := strconv.Itoa(int(hr))
	if hr < 10 {
		str_hr = "0" + str_hr
	}
	str_mi := strconv.Itoa(int(mi))
	if mi < 10 {
		str_mi = "0" + str_mi
	}
	str_ss := strconv.Itoa(int(ss))
	if ss < 10 {
		str_ss = "0" + str_ss
	}
	HHMISS = str_hr + str_mi + str_ss
	return YYYYMMDD, HHMISS
}

func GetTimeParts(_Date time.Time) (YYYY string, MM string, MMM string, DD string, DD_Ordinal string, HH string, MI string) {
	yyyy, mm, dd := _Date.Date()
	hr, mi, _ := _Date.Clock()

	YYYY = strconv.Itoa(yyyy)
	MM = strconv.Itoa(int(mm))
	if mm < 10 {
		MM = "0" + MM
	}
	MMM = mm.String()[:3]

	DD = strconv.Itoa(dd)
	if dd < 10 {
		DD = "0" + DD
	}
	DD_Ordinal = Ordinal(dd)
	HH = strconv.Itoa(int(hr))
	if hr < 10 {
		HH = "0" + HH
	}
	MI = strconv.Itoa(int(mi))
	if mi < 10 {
		MI = "0" + MI
	}
	return YYYY, MM, MMM, DD, DD_Ordinal, HH, MI
}

type TimeParts_V2 struct {
	TargetDate time.Time
	YYYY       string
	MM         string
	MMM        string
	DD         string
	DD_Ordinal string
	HH         string
	MI         string
	SE         string
	Week       string
	Weekday    string
}

func GetTimeParts_V2(_Date time.Time) (timeparts TimeParts_V2) {
	timeparts.TargetDate = _Date
	yyyy, mm, dd := _Date.Date()
	hr, mi, se := _Date.Clock()

	timeparts.YYYY = strconv.Itoa(yyyy)
	timeparts.MM = strconv.Itoa(int(mm))
	if mm < 10 {
		timeparts.MM = "0" + timeparts.MM
	}
	timeparts.MMM = mm.String()[:3]

	timeparts.DD = strconv.Itoa(dd)
	if dd < 10 {
		timeparts.DD = "0" + timeparts.DD
	}
	timeparts.DD_Ordinal = Ordinal(dd)
	timeparts.HH = strconv.Itoa(int(hr))
	if hr < 10 {
		timeparts.HH = "0" + timeparts.HH
	}
	timeparts.MI = strconv.Itoa(int(mi))
	if mi < 10 {
		timeparts.MI = "0" + timeparts.MI
	}
	timeparts.SE = strconv.Itoa(int(se))
	if se < 10 {
		timeparts.SE = "0" + timeparts.SE
	}
	_, wk_int := _Date.ISOWeek()
	wk := strconv.Itoa(int(wk_int))
	if wk_int < 10 {
		timeparts.Week = "0" + wk
	} else {
		timeparts.Week = wk
	}
	timeparts.Weekday = strconv.Itoa(int(_Date.Weekday()))
	return
}

func GetRequestIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}
	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}

func GetDestinationPortFromRequest(r *http.Request) (port string, err error) {
	//Get IP from the X-Original-Host header -- Behind our API Gateway
	host := r.Header.Get("X-Original-Host")
	if host != "" {
		_, port, err = net.SplitHostPort(host)
		if err == nil {
			return
		} else {
			fmt.Println(err)
			return "", fmt.Errorf("no valid destination port")
		}
	}

	//Get IP from host -- Running locally
	_, port, err = net.SplitHostPort(r.Host)
	if err == nil {
		return
	} else {
		fmt.Println(err)
		return "", fmt.Errorf("no valid destination port")
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// //Hashing Functions///////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////////////
//var Hash_Secret_Key string = "A3k%a2l&$&&34Fo1~~2Fo|003j|`j%&*hlksalkdj7|~jb343!2e09df{12$$^1u)U(*@("

func Hashing(input string, secretKey string) (Output string) {
	md5 := md5.New()
	sha_256 := sha256.New()
	sha_512 := sha512.New()
	io.WriteString(md5, input)
	sha_256.Write([]byte(input))
	sha_512.Write([]byte(input))
	//sha_512_256 := sha512.Sum512_256([]byte(input))
	hmac512 := hmac.New(sha512.New, []byte(secretKey))
	hmac512.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(hmac512.Sum(nil))

	//fmt.Printf("md5:\t\t%x\n", md5.Sum(nil))
	//fmt.Printf("sha256:\t\t%x\n", sha_256.Sum(nil))
	//fmt.Printf("sha512:\t\t%x\n", sha_512.Sum(nil))
	//fmt.Printf("sha512_256:\t%x\n", sha_512_256)
	//fmt.Printf("hmac512:\t%s\n", base64.StdEncoding.EncodeToString(hmac512.Sum(nil)))
}

// format int64 number with comma separators
func FormatInt64(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

type EVC_Balance_Response struct {
	Data              float64 `bson:"Data" json:"Data"`
	Status            string  `bson:"Status" json:"Status"`
	StatusCode        int     `bson:"StatusCode" json:"StatusCode"`
	StatusDescription string  `bson:"StatusDescription" json:"StatusDescription"`
	ErrorDescription  string  `bson:"ErrorDescription" json:"ErrorDescription"`
}

func Get_EVC_Balance(MSISDN string) (balance float64, _rErr error) {
	url := "http://10.250.1.229:9903/HTTP_EVC_Dealer_Bal/?MSISDN=" + MSISDN
	method := "GET"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return balance, err
	}
	//client := &http.Client{}
	client := &http.Client{
		Timeout: 10 * time.Second, //is SMSC not reachable request will time out after "TimeOut" sec
	}
	resp, err := client.Do(req)
	if err != nil {
		return balance, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			err := errors.New("failed to read evc balance response body")
			return balance, err
		}
		var Balance_Response EVC_Balance_Response
		err = json.Unmarshal(body, &Balance_Response)
		if err != nil {
			return balance, err
		}
		return Balance_Response.Data, nil
	} else {
		err := errors.New("error getting evc balance (StatusCode: " + fmt.Sprint(resp.StatusCode) + ")")
		return balance, err
	}
}

// //////////////////////////////////////////////////
// Rounds like 12.3416 -> 12.35
func RoundUp(val float64, precision int) float64 {
	return math.Ceil(val*(math.Pow10(precision))) / math.Pow10(precision)
}

// Rounds like 12.3496 -> 12.34
func RoundDown(val float64, precision int) float64 {
	return math.Floor(val*(math.Pow10(precision))) / math.Pow10(precision)
}

// Rounds to nearest like 12.3456 -> 12.35
func RoundToNearest(val float64, precision int) float64 {
	return math.Round(val*(math.Pow10(precision))) / math.Pow10(precision)
}
