package helper

import (
	"benchmarks-gin/config"
	"benchmarks-gin/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/maps"
)

var ctx = context.Background()
var SessionExpirationTime time.Duration = 1800
var CachePrefix string = "CACHE_MANAGER_"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const intBytes = "0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RequestAppendHeader(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if c.Request.Method == "OPTIONS" {
		c.JSON(200, gin.H{"success": 1})
		panic("done")
	}
}
func CorsReply(c *gin.Context) {
	// time.Sleep(5 * time.Second)
	RequestAppendHeader(c)
}

func RandString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
func RandInt(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(intBytes) {
			b[i] = intBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
func GetUniqueSecret(key *string) (string, string) {
	keyCode := RandString(12)
	if key != nil {
		keyCode = *key
	}
	secret := fmt.Sprintf("%s.%s", os.Getenv("secret"), keyCode)
	return keyCode, secret
}
func ParseTemplate(templateFileName string, data interface{}) (*string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	body := buf.String()
	return &body, nil
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func ConvertPlaceholder(placeholder string) string {
	// placeholder = strings.ReplaceAll(placeholder, "\"", "")
	if !strings.Contains(placeholder, "##") {
		return placeholder
	}
	switch placeholder {
	case "##uuid":
		return gocql.TimeUUID().String()
	case "##email":
		return RandString(8) + "@" + RandString(4) + "." + RandString(3)
	case "##string":
		return RandString(12)
	case "##phone":
		return RandInt(12)
	case "##now":
		return time.Now().Format(time.RFC3339)
	case "##date":
		return time.Now().Format(time.DateOnly)
	}
	return placeholder
}
func GenerateSqlQuery(data model.Payload, result map[string]json.RawMessage, parentRecord map[string]json.RawMessage) (int, error) {
	records := data.Records
	recordsKeys := maps.Keys(records)
	incrementValue := 0
	type RecordsStruct struct {
		Key   string
		Value string
	}
	r, _ := regexp.Compile("@@([a-z]+).([a-z]+)")
	var recordsLoop []RecordsStruct
	//extract column that have a relation to other table
	for _, kx := range recordsKeys {
		vx := string(records[kx])
		vx = strings.ReplaceAll(vx, "\"", "")
		if r.MatchString(vx) {
			//relation found
			recordsLoop = append(recordsLoop, RecordsStruct{Key: kx, Value: vx})
		}
	}
	if len(recordsLoop) == 0 {
		recordsLoop = append(recordsLoop, RecordsStruct{Key: "default", Value: "default"})
	}
	sqlColumns := ""
	sqlValues := ""
	for a := 0; a < len(recordsLoop); a++ {
		recordStruct := recordsLoop[a]
		//no relation to other table
		if recordStruct.Key == "default" {
			//TODO: generate and insert single records
			for _, key := range recordsKeys {
				val := string(records[key])
				val = strings.ReplaceAll(val, "\"", "")
				loopRecordsKey(key, val, result, parentRecord, &sqlColumns, &sqlValues)
			}
			sqlColumns = fmt.Sprintf("%s)", sqlColumns)
			sqlValues = fmt.Sprintf("%s)", sqlValues)
			sql := fmt.Sprintf("insert into %s %s values %s", data.Table, sqlColumns, sqlValues)
			fmt.Println(sql)
			fmt.Println()
			err := config.SESSION.Query(sql).Exec()
			if err != nil {
				return 0, err
			}
			incrementValue = incrementValue + 1
		} else {
			//join to other table and do multiple insert based on the result
			var iter *gocql.Iter
			val := recordStruct.Value
			val = strings.ReplaceAll(val, "@@", "")
			val = strings.ReplaceAll(val, "\"", "")
			//get additional query conditions if any
			dbQuery := strings.SplitN(val, " ", 2)
			//separate table and column
			dbData := strings.Split(dbQuery[0], ".")
			appendQuery := ""

			if len(dbQuery) > 1 {
				appendQuery = dbQuery[1]
			}
			//generate query to fetch the linked data
			sql := fmt.Sprintf("select %s from %s %s", dbData[1], dbData[0], appendQuery)
			iter = config.SESSION.Query(sql).Iter()

			if iter != nil && iter.NumRows() != 0 {
				var dbVal string
				for iter.Scan(&dbVal) {
					for _, key := range recordsKeys {
						val := string(records[key])
						//ignore this key value and use the data from db
						if key == recordStruct.Key {
							val = dbVal
						} else {
							val = strings.ReplaceAll(val, "\"", "")
							//check if
							if r.MatchString(val) {
								//
							}
						}
						loopRecordsKey(key, val, result, parentRecord, &sqlColumns, &sqlValues)
					}
					sqlColumns = fmt.Sprintf("%s)", sqlColumns)
					sqlValues = fmt.Sprintf("%s)", sqlValues)
					// sqlValues = fmt.Sprintf(sqlValues, dbVal)
					sql := fmt.Sprintf("insert into %s %s values %s", data.Table, sqlColumns, sqlValues)
					fmt.Println("sql values:", sqlValues)
					fmt.Println("db val:", dbVal)
					fmt.Println(sql)
					fmt.Println()
					err := config.SESSION.Query(sql).Exec()
					if err != nil {
						return 0, err
					}
				}
				incrementValue = incrementValue + iter.NumRows()
			} else {
				//TODO: halt because of invalid data from user
			}
		}
	}
	return incrementValue, nil
}

func ValidateUuid(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
func ValidateInt(val string) bool {
	r := regexp.MustCompile("^[0-9]+$")
	return r.MatchString(val)
}

func ParseStringArray(input string) []string {
	// Remove square brackets from the input
	trimmedInput := strings.Trim(input, "[]")

	// Split the remaining string by commas
	stringArray := strings.Split(trimmedInput, ",")

	// Trim spaces from each element (optional)
	for i, value := range stringArray {
		stringArray[i] = strings.TrimSpace(value)
	}

	return stringArray
}
func GetRandomKey(arr []string) string {
	// Generate a random index within the bounds of the array
	randomIndex := rand.Intn(len(arr))

	// Return the string at the random index
	return arr[randomIndex]
}
func loopRecordsKey(key string, val string, result map[string]json.RawMessage, parentRecord map[string]json.RawMessage, sqlColumns *string, sqlValues *string) {
	rArray, _ := regexp.Compile(`^\[.*\]$`)
	if strings.Contains(val, "@@") {
		//data link to the parent request is found
		keyy := strings.ReplaceAll(val, "@@", "")
		keyy = strings.ReplaceAll(keyy, "\"", "")
		fmt.Println(keyy, ":", string(parentRecord[keyy]))
		val = string(parentRecord[keyy])
	} else if rArray.MatchString(val) {
		//is an enum to choose from
		texts := ParseStringArray(val)
		val = GetRandomKey(texts)
	} else {
		//check if there is a placeholder and convert it
		val = ConvertPlaceholder(val)
	}
	escapeValue := true
	// if ValidateInt(val) || (ValidateUuid(val) && (os.Getenv("database_type") == "cassandra" || os.Getenv("database_type") == "scylladb")) {
	if ValidateInt(val) && len(val) < 10 {
		//ignore escape int
		escapeValue = false
	}

	if len(*sqlColumns) == 0 {
		*sqlColumns = fmt.Sprintf("(%s", key)
		if !escapeValue {
			*sqlValues = fmt.Sprintf("(%s", val)
		} else {
			*sqlValues = fmt.Sprintf("('%s'", val)
		}
	} else {
		*sqlColumns = fmt.Sprintf("%s,%s", *sqlColumns, key)
		if !escapeValue {
			*sqlValues = fmt.Sprintf("%s, %s", *sqlValues, val)
		} else {
			*sqlValues = fmt.Sprintf("%s, '%s'", *sqlValues, val)
		}
	}
	result[key] = json.RawMessage(val)
}
