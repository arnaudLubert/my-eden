/*
** Made by LUBERT Arnaud, Epitech Student (Promo 2023)
** 06/02/2021 arnaud.lubert@epitech.eu
**
*/
package main

import (
    "golang.org/x/text/unicode/norm"
    "golang.org/x/text/transform"
    "math/rand"
    "net/http"
    "strings"
    "strconv"
    "unicode"
    "reflect"
    "errors"
    "time"
//    "log"
)

// append c char at the end of str string
func strAppend(str string, c byte) string {
    var stringBuilder strings.Builder

    stringBuilder.WriteString(str)
    stringBuilder.WriteByte(c)

    return stringBuilder.String()
}

func concatByteArray(a []byte, b []byte) []byte {
    array := make([]byte, len(a) + len(b))
    i := len(b) - 1
    j := len(b) + len(a) - 1

    for i != -1 {
        array[j] = b[i]
        i--
        j--
    }

    i = len(a) - 1
    for i != -1 {
        array[j] = a[i]
        i--
        j--
    }

    return array
}

// add line break '\n' when character limit is reached on the line
func strMaxLength(rawStr string, limit int) string {
    var string strings.Builder
    var buff strings.Builder // store characters until a space/ln is reached
    var space bool // temp space, break mark
    str := []rune(rawStr)
    len := 0

    for i := range(str) {
        if str[i] == ' ' {
            if space {
                string.WriteByte(' ')
            }
            string.WriteString(buff.String())
            buff.Reset()
            space = true
            len++
        } else if str[i] == '\n' {
            if space {
                string.WriteByte(' ')
            }
            string.WriteString(buff.String())
            string.WriteByte('\n')
            buff.Reset()
            space = false
            len = 0
        } else if len >= limit {
            string.WriteByte('\n')
            buff.WriteRune(str[i])
            space = false
            len = 0
        } else {
            buff.WriteRune(str[i])
            len++
        }
    }
    if space {
        string.WriteByte(' ')
    }
    string.WriteString(buff.String())

    return string.String()
}

// convert a float to price string (1.2 -> 1,20)
func formatPrice(value float64) string {
    str := strconv.FormatFloat(float64(value), 'f', 2, 32)
    str = strings.Replace(str, ".", ",", 1)

    if strings.IndexByte(str, '.') == len(str) - 2 {
        str = strAppend(str, '0')
    }

    return (str + " €")
}

// convert a size (bytes) to string (20000 -> 20 Ko)
func formatSize(sz int64) string {
    size := float64(sz)

    if size < 1000 {
        return strconv.FormatFloat(size, 'f', 0, 64) + " o"
    } else if size < 1000000 {
        return strconv.FormatFloat(size / 1000, 'f', 0, 64) + " Ko"
    } else if size < 100000000 {
        return strconv.FormatFloat(size / 1000000, 'f', 0, 64) + " Mo"
    }

    return strconv.FormatFloat(size / 100000000, 'f', 1, 64) + " Go"
}

func formatDate(date time.Time) string {
    return date.Format("02/01/2006 15:04:05")
}

// "Tue, 10 Nov 2009 23:00:00 GMT" --> "2009-11-10"
func sitemapDate(formated string) string{
    date, err := time.Parse(http.TimeFormat, formated)

    if err != nil {
        logging(err.Error())
        return ""
    }
	return date.Format("2006-01-02")
}

func getWeight(value float32) string {
    str := strconv.FormatFloat(float64(value / 1000), 'f', 2, 32)
    str = strings.Replace(str, ".", ",", 1)

    if strings.IndexByte(str, '.') == len(str) - 2 {
        str = strAppend(str, '0')
    }

    return str
}

// convert 2020-12-31 24:60:60 to 31/12/2020 24:60
// convert 2020-12-31 to 31/12/2020
// else: return empty string
func getDateFromSql(date string) string {
    var chars []byte
    len := len(date)

    if len == 19 {
        chars = make([]byte, 16)
    } else if len > 9 {
        chars = make([]byte, len)
    } else {
        return ""
    }
    chars[0] = date[8]
    chars[1] = date[9]
    chars[2] = '/'
    chars[3] = date[5]
    chars[4] = date[6]
    chars[5] = '/'
    chars[6] = date[0]
    chars[7] = date[1]
    chars[8] = date[2]
    chars[9] = date[3]

    if len == 19 {
        chars[10] = ' '
        chars[11] = date[11]
        chars[12] = date[12]
        chars[13] = ':'
        chars[14] = date[14]
        chars[15] = date[15]
    }

    return string(chars)
}

// 0612345678 -> 06.12.34.56.78
func formatPhoneNumber(number string) string {
    if len(number) == 10 {
        return number[:2] + "." + number[2:4] + "." + number[4:6] + "." + number[6:8] + "." + number[8:]
    }
    return number
}

func getMonthFr(index int) string {
    switch(index) {
    case 0: return "Janvier"
    case 1: return "Février"
    case 2: return "Mars"
    case 3: return "Avril"
    case 4: return "Mai"
    case 5: return "Juin"
    case 6: return "Juillet"
    case 7: return "Août"
    case 8: return "Septembre"
    case 9: return "Octobre"
    case 10: return "Novembre"
    case 11: return "Décembre"
    }

    return strconv.Itoa(index + 1)
}

func getDayFr(index int) string {
    switch(index) {
    case 0: return "Dimanche"
    case 1: return "Lundi"
    case 2: return "Mardi"
    case 3: return "Mercredi"
    case 4: return "Jeudi"
    case 5: return "Vendredi"
    case 6: return "Samedi"
    }

    return "-"
}

func addWorkDays(date *time.Time, add int) {
    day := int(date.Weekday())
    count := 0

    for add != 0 {
        count++
        if (count + day) % 6 == 0 {
            count += 2
        }
        add--
    }
    *date = (*date).Add(time.Duration(count) * time.Hour * 24)
}

// generate a random string of size length
// only return readable characters
func randomString(length int) string {
    var stringBuilder strings.Builder
    var char byte
    r := rand.New(rand.NewSource(time.Now().Unix()))

    for i := 0; i != length; i++ {
        char = byte(r.Int31n(91) + 35)

        if char == ';' || char == '=' || char == '`' || char == '_' || char == '|' || char == '\\' {
            stringBuilder.WriteByte('#')
        } else {
            stringBuilder.WriteByte(char)
        }

    }

    return stringBuilder.String()
}

func safeRandomString(length int) string {
    var stringBuilder strings.Builder
    var char byte
    r := rand.New(rand.NewSource(time.Now().Unix()))

    for i := 0; i != length; i++ {
        char = byte(r.Int31n(74) + 48)

        if char > '9' && char < 'A' {
            stringBuilder.WriteByte('A')
        } else if char > 'Z' && char < 'a' {
            stringBuilder.WriteByte('a')
        } else {
            stringBuilder.WriteByte(char)
        }
    }

    return stringBuilder.String()
}

// remove the port from ip address
// 192.168.1.23:47716  -->  192.168.1.23
func getAddress(remoteAddr string) string {
    index := strings.LastIndexByte(remoteAddr, ':')

    if index != -1 {
        return remoteAddr[:index]
    }
    return remoteAddr
}

func strExistsIn(array *[]string, compared string) bool {

    for i := len(*array) - 1; i != -1; i-- {
        if (*array)[i] == compared {
            return true
        }
    }
    return false
}

// append tags to base (result is stored in base)
func concatTags(base *[]string, rawTags *string) {
    var found bool

    i := 0
    tags := strings.Split(*rawTags, ",")
    newTags := make([]string, len(tags) + len(*base))

    for j := range(tags) {
        found = false
        for k := range(*base) {
            if tags[j] == (*base)[k] {
                found = true
                break
            }
        }
        if !found && tags[i] != "" {
            newTags[i] = tags[j]
            i++
        }
    }
    newTags = newTags[:i]
    *base = append(*base, newTags...)
}

// remove all carriage return & spaces in the byte array (read from a file)
func obfuscateBytes(bytes *[]byte) {
    j := 0
    len := len(*bytes)
    obfuscatedBytes := make([]byte, len)

    comment := 0 // 1== // 2 == /***/
    quote := 0 // 1=='\'' 2=='\"'

    for i := 0; i != len; i++ {

        if comment == 1 {
            if (*bytes)[i] == '\n' {
                comment = 0
            }
        } else if comment == 2 {
            if (*bytes)[i] == '/' && (*bytes)[i - 1] == '*' {
                comment = 0
            }
        } else {
            if quote == 0 {
                if (*bytes)[i] == '\'' && (*bytes)[i - 1] != '\\' {
                    quote = 1
                } else if (*bytes)[i] == '"' && (*bytes)[i - 1] != '\\' {
                    quote = 2
                } else if (*bytes)[i] == '/' && (*bytes)[i + 1] == '/' {
                    comment = 1
                } else if (*bytes)[i] == '/' && (*bytes)[i + 1] == '*' {
                    comment = 2
                }
            } else if (quote == 1 && (*bytes)[i] == '\'') && (*bytes)[i - 1] != '\\' ||
            (quote == 2 && (*bytes)[i] == '"') && (*bytes)[i - 1] != '\\' {
                quote = 0
            }

            if comment == 0 && (quote != 0 || ((*bytes)[i] != '\n' && (*bytes)[i] != '\t' &&
            ((*bytes)[i] != ' ' || ((*bytes)[i] == ' ' && (*bytes)[i - 1] != ' ') &&
            (*bytes)[i - 1] != ')' && (*bytes)[i - 1] != '}' && (*bytes)[i - 1] != ']' &&
            (*bytes)[i - 1] != ';' && (*bytes)[i - 1] != ',' && (*bytes)[i - 1] != '='))) {
                obfuscatedBytes[j] = (*bytes)[i]
                j++
            }
        }
    }
    *bytes = obfuscatedBytes[:j]
}

// use company or lname + fname
func createRefClient(company string, lname string, fname string) string {
    var str string

    if company != "" {
        str = company
    } else if fname != "" {
        str = lname + "_" + fname
    }

    return strings.ToLower(strings.Replace(str, " ", "_", -1))
}

// reduce max length of 700 characters XOR limit to 4 <br>
// also remove all html tags (except <br>)
func getBlogPreview(content string) string{
    var stringBuilder strings.Builder
    var tag bool
    var ln int

    chars := []byte(content)
    length := len(chars)

    for i := range(chars) {
        if !tag && chars[i] == '<' {
            if i < length - 2 && chars[i + 2] == 'r' { // <br>
                ln++

                if ln == 4 {
                    break
                }
                stringBuilder.WriteByte(chars[i])
            } else {
                tag = true
            }
        } else if tag {

            if chars[i] == '>' {
                tag = false
            }
        } else if i > 700 {
            break
        } else {
            stringBuilder.WriteByte(chars[i])
        }
    }

    return stringBuilder.String()
}

// create a link from a title (blog)
// "Comment envoyer un fichier ?" --> comment-envoyer-un-fichier
// regex [0-9a-z]-
func linkFromTitle(title string) string {
    var str []byte
    var stringBuilder strings.Builder
    st := strings.ToLower(title)

// remove diacritics
    t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
    result, _, err := transform.String(t, st)

    if err != nil {
        logging(err.Error())
        str = []byte(st)
    } else {
        str = []byte(result);
    }

// prevent segfault of: str[i + 1] != '?'
    if str[len(str) - 1] == ' ' {
        str = str[:len(str) - 1]
    }

// [0-9a-z]-
    for i := range(str) {
        if str[i] == ' ' && str[i + 1] != '?' {
            stringBuilder.WriteByte('-')
        } else if (str[i] >= '0' && str[i] <= '9') || (str[i] >= 'A' && str[i] <= 'Z') || (str[i] >= 'a' && str[i] <= 'z') || str[i] == '-' {
            stringBuilder.WriteByte(str[i])
        }
    }
    return stringBuilder.String()
}

// linkFromTitle() diacritics
func isMn(r rune) bool {
    return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func appendIfNew(array *[]string, new *string) bool {

    for i := range(*array) {
        if (*array)[i] == *new {
            return false
        }
    }
    *array = append(*array, *new)
    return true
}

// create a mysql request using an undefined map
func requestFromMap(data *map[string]interface{}, exceptions []string) (string, error) {
    var val, request string

    for k, e := range(*data) {
        if !strExistsIn(&exceptions, k) {
            switch v := e.(type) {
            case bool:
                if v {
                    val = "1"
                } else {
                    val = "0"
                }
            case int:
                val = strconv.Itoa(v)
            case float64:
                val = strconv.FormatFloat(v, 'f', -1, 64)
            case string:
                val = strings.Replace(v, "\"", "\\\"", -1)
            case []interface{}:
                if len(v) == 0 {
                    val = ""
                } else {
                    if reflect.TypeOf(v[0]).Name() == "string" {
                        stringSlice := make([]string, len(v))

                        for i := range(stringSlice) {
                            stringSlice[i] = v[i].(string)
                        }
                        val = strings.Join(stringSlice, ",")
                    } else {
                        return "", errors.New("Unknown type (" + reflect.TypeOf(e).String() + ": " + k + ")")
                    }
                }
            default:
                return "", errors.New("Unknown type (" + reflect.TypeOf(e).String() + ": " + k + ")")
            }

            if request == "" {
                request += "" + k + " = \"" + val + "\""
            } else {
                request += ", " + k + " = \"" + val + "\""
            }
        }
    }

    return request, nil
}

// remove all html tags to create plain text
func getPlainText(htmlStr string) string {
	var stringBuilder strings.Builder
	var inTag bool

	htmlStr = strings.Replace(htmlStr, "<br>", "\r\n", -1)
	length := len(htmlStr)
	i := 0

	for i != length {
		if htmlStr[i] == '<' {
			inTag = true
		} else if htmlStr[i] == '>' {
			inTag = false
		} else if !inTag {
			if htmlStr[i] == 195 && i + 1 < length {
				stringBuilder.WriteString(htmlStr[i:i + 2])
				i++
			} else {
				stringBuilder.WriteByte(htmlStr[i])
			}
		}
		i++
	}
	return stringBuilder.String()
}

// remove all html tags to create plain text
func getPlainTextFromBytes(htmlBytes []byte) string {
	var stringBuilder strings.Builder
	var inTag bool

	length := len(htmlBytes)
	i := 0

	for i != length {
        if i + 3 < length && htmlBytes[i] == '<' && htmlBytes[i + 1] == 'b' && htmlBytes[i + 2] == 'r' && htmlBytes[i + 3] == '>' {
            stringBuilder.WriteString("\r\n")
            i += 3
        } else if htmlBytes[i] == '<' {
			inTag = true
		} else if htmlBytes[i] == '>' {
			inTag = false
		} else if !inTag {
            if htmlBytes[i] == 195 && i + 1 < length {
				stringBuilder.WriteByte(htmlBytes[i])
                stringBuilder.WriteByte(htmlBytes[i + 1])
				i++
			} else {
				stringBuilder.WriteByte(htmlBytes[i])
			}
		}
		i++
	}
	return stringBuilder.String()
}

func stringContains(base string, array *[]string) bool {
	lBase := len(base)
	lArray := len(*array)
	i := 0
	j := 0
	k := 0

	for i != lBase {
		j = 0
		for j != lArray {
			k = len((*array)[j]) - 1

			if k > -1 && lBase - i > k {
				for k != -1 && (*array)[j][k] == base[i + k] {
					k--
				}

				if k == -1 {
					return true
				}
			}
			j++
		}
		i++
	}

	return false
}
