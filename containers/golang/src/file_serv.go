/*
** Made by LUBERT Arnaud, Epitech Student (Promo 2023)
** 06/02/2021 arnaud.lubert@epitech.eu
**
*/
package main

import (
    "compress/gzip"
//    "archive/zip"
    "io/ioutil"
    "net/http"
    "os/exec"
    "strconv"
    "strings"
    "bufio"
    "path"
//    "time"
    "io"
    "os"
)

const (
    ALL   int = 0
    TEXT  int = 1
    FONT  int = 2
    IMAGE int = 3
    VIDEO int = 4
    AUDIO int = 5
    DOC   int = 6
    MODEL int = 7
)

// Set the correct header for every request
// return true if enconding is available
/*
** cache levels:
** 0: no cache
** 1: cache until Last-Modified date is not modified
** 2: cache forever
*/
func setHeader(w *http.ResponseWriter, req *http.Request, cache uint8, lastModDate *string, contentType string) bool {
    Header := (*w).Header()
    Header.Set("Content-type", contentType + "; charset=utf-8") // can be utf-16
    Header.Set("Server", "La Citadelle du web (1.0)")
    Header.Set("Connection", "Keep-Alive")
    Header.Set("Keep-Alive", "timeout=3, max=100")
    Header.Set("Vary", "Accept-Encoding")
    Header.Set("Link", "<" + BaseUrl + req.URL.Path + ">; rel=\"canonical\"")

    if contentType == "application/javascript" {
        Header.Set("Content-Language", "en-US")
    } else {
        Header.Set("Content-Language", "fr-FR")
    }
    if cache != 0 && (strings.Contains(req.Header.Get("Cache-control"), "no-cache") || strings.Contains(req.Header.Get("Pragma"), "no-cache")) {
        cache = 0
    }

    /*switch cache {
    case 0:*/
        Header.Set("Pragma", "no-cache")
        Header.Set("Cache-control", "no-store")
    /*case 1:
        Header.Set("Last-Modified", *lastModDate)
        Header.Set("Cache-control", "max-age=3600") // 60 minutes of cache before asking the ressource
        //Header.Set("Cache-control", "no-cache") // always ask for the ressource
    case 2:
        Header.Set("Cache-control", "public, max-age=15552000, immutable")
    }*/

    if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
        Header.Set("Content-Encoding", "gzip")
        return true
    }
    return false
}

func setSessionCookies(rw *http.ResponseWriter, token string, expire string) {
    Header := (*rw).Header()

    if strings.Contains(BaseUrl, "https://") {
        Header.Set("Set-Cookie", "session=" + token + "; Secure; SameSite=Lax; Path=/; Expires=" + expire)
    } else {
        Header.Set("Set-Cookie", "session=" + token + "; SameSite=Lax; Path=/; Expires=" + expire)
    }
    Header.Set("Server", "La Citadelle du web (1.0)")
    Header.Set("Connection", "Keep-Alive")
    Header.Set("Keep-Alive", "timeout=3, max=100")
    Header.Set("Vary", "Accept-Encoding")
}

// return an io.Writer wich can be a gzip.Writer or the default ResponseWriter
// also return a pointer to gzip.Writer (if compressed), to close when writing is done
func compressedWritter(rw *http.ResponseWriter, compressed bool) (*io.Writer, *gzip.Writer) {
    if compressed {
        zip := gzip.NewWriter(*rw)
        cast := new(io.Writer)
        *cast = ((io.Writer)(zip))

        return cast, zip
    } else {
        cast := new(io.Writer)
        *cast = ((io.Writer)(*rw))

        return cast, nil
    }
}

// clear writers and close gzip.Writer if exists
func freeCompression(w **io.Writer, zip *gzip.Writer) {
    if zip != nil {
        zip.Close()
    }
    *w = nil
}

// write 401 error page on ResponseWriter
func unauthorizedResponse(rw *http.ResponseWriter) {
    (*rw).WriteHeader(http.StatusUnauthorized)
}

func servErr(rw *http.ResponseWriter, code int) {

    (*rw).WriteHeader(code)
/*
    switch(code) {
    case http.StatusBadRequest:
        servError(rw, code, "Il semblerait que le serveur n'ai pas reçu les bons paramètres.")
    case http.StatusUnauthorized:
        servError(rw, code, "Contenu à accès restreint.")
    case http.StatusNotFound:
        serv404(rw)
    case http.StatusInternalServerError:
        servError(rw, code, "Il semblerait qu'il y ait un problème lié au serveur.")
    default:
        servError(rw, code, "")
    }*/
}

// write custom error page on ResponseWriter
func servError(rw *http.ResponseWriter, code int, m string) {
    (*rw).WriteHeader(code)
}

func serv404(rw *http.ResponseWriter) {
    (*rw).WriteHeader(http.StatusNotFound)
}

/*
// check browser version
func browserCompatibility(userAgent string) bool {
    var safariSet bool

    if stringContains(userAgent, &[]string{"Windows NT", "bot", "slurp", "crawl", "facebook", "qwant", "Iron", "GSA"}) { // Servers and crawlers
        return true
    } else if strings.Contains(userAgent, "MSIE") { // Internet Explorer
        return false
    }

    if len(userAgent) > 6 && userAgent[:6] == "Safari" {
        safariSet = true
    }
	parts := strings.Split(userAgent, " ")
	var separator, end, version int
	var verionStr string
	var err error

	for i := range(parts) {
		separator = strings.IndexByte(parts[i], '/')

		if separator != -1 {
			if end = strings.IndexByte(parts[i], '.'); end == -1 {
			    end = len(parts[i])
			}
            verionStr = parts[i][separator + 1:end]

            if verionStr[0] < '0' || verionStr[0] > '9' { // not a verions number
                continue
            } else if version, err = strconv.Atoi(verionStr); err != nil {
        		logging(err.Error(), userAgent)
				continue
			}

			switch (parts[i][:separator]) {   // 2014-2016
			case "Chrome":  if version < 42 { return false } else { return true }
			case "Edge":    if version < 14 { return false } else { return true }
            case "Firefox": if version < 39 { return false } else { return true }
            case "Trident", "IE": return false // MSIE IE
			case "Opera":   if version < 29 { return false } else { return true }
            case "Version":  if version < 10 { return false } else { return true } // Version/10.0.3 Safari/602.4.8
            case "Safari":
                if safariSet && version < 602 {
                    logging("Obsolete safari detected:")
                    return false
                }
            }
		}
	}
	return true
}
*/
// Handles /img & /pictures request
func images(rw http.ResponseWriter, req *http.Request) {
    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, getContentType(req.URL.Path, IMAGE))
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

func imgArticles(rw http.ResponseWriter, req *http.Request) {

    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, getContentType(req.URL.Path, IMAGE))
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

func imgSupports(rw http.ResponseWriter, req *http.Request) {

    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, getContentType(req.URL.Path, IMAGE))
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

func imgSlideshow(rw http.ResponseWriter, req *http.Request) {

    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, getContentType(req.URL.Path, IMAGE))
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

func imgBlog(rw http.ResponseWriter, req *http.Request) {

    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, getContentType(req.URL.Path, IMAGE))
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

// Handles /js/ request
func js(rw http.ResponseWriter, req *http.Request) {

    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, "application/javascript")
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

// Handles /css/ request
func css(rw http.ResponseWriter, req *http.Request) {
    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, "text/css")
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

// Handles /models/ request
func models(rw http.ResponseWriter, req *http.Request) {
    file, err := os.Open(req.URL.Path[1:])

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        return
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        servErr(&rw, http.StatusNotFound)
        file.Close()
        return
    }
    file.Close()

    compressed := setHeader(&rw, req, 1, nil, getContentType(req.URL.Path, MODEL))
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
}

func servAny(rw http.ResponseWriter, req *http.Request) {

    if !servFile(&rw, req, req.URL.Path[1:]) {
        servErr(&rw, http.StatusNotFound)
    }
}

// serve any type of file (folders are prohibited)
func servFile(rw *http.ResponseWriter, req *http.Request, path string) bool {
    file, err := os.Open(path)

    if (err != nil) {
        return false
    }
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)

    if (err != nil) {
        file.Close()
        return false
    }
    file.Close()

    compressed := setHeader(rw, req, 0, nil, getContentType(path, ALL))
    w, zip := compressedWritter(rw, compressed)
    (*w).Write(content)
    freeCompression(&w, zip)
    return true
}

// Return Content-Type string from filename
func getContentType(path string, target int) string {
    idx := strings.LastIndexByte(path, '.')

    if (idx == -1) {
        return "application/octet-stream"
    }
    extension := path[idx + 1:]

    if (target == ALL || target == TEXT) {
        switch (extension) {
        case "html":
            return "text/html"
        case "css":
            return "text/css"
        case "js":
            return "application/javascript"
        case "json":
            return "application/json"
        case "pdf":
            return "application/pdf"
        case "php":
            return "application/php"
        case "txt":
            return "text/plain"
        case "xml":
            return "application/xml"
        case "eot":
            return "application/vnd.ms-fontobject"
        case "htm":
            return "text/html"
        case "bin":
            return "application/octet-stream"
        }
    }
    if (target == ALL || target == FONT) {
        switch (extension) {
        case "ttf":
            return "font/ttf"
        case "rtf":
            return "application/rtf"
        case "woff":
            return "font/woff"
        case "woff2":
            return "font/woff2"
        case "otf":
            return "font/otf"
        }
    }
    if (target == ALL || target == IMAGE) {
        switch (extension) {
        case "jpg":
            return "image/jpeg"
        case "jpeg":
            return "image/jpeg"
        case "png":
            return "image/png"
        case "ico":
            return "image/vnd.microsoft.icon"
        case "gif":
            return "image/gif"
        case "heic":
            return "image/heic"
        case "heif":
            return "image/heif"
        case "webp":
            return "image/webp"
        case "svg":
            return "image/svg+xml"
        case "bmp":
            return "image/bmp"
        case "tif":
            return "image/tiff"
        case "tiff":
            return "image/tiff"
        }
    }
    if (target == ALL || target == VIDEO) {
        switch (extension) {
        case "mpeg":
            return "video/mpeg"
        case "mp4":
            return "video/mp4"
        case "webm":
            return "video/webm"
        case "avi":
            return "video/x-msvideo"
        case "ogv":
            return "video/ogg"
        case "ogx":
            return "application/ogg"
        case "swf":
            return "application/x-shockwave-flash"
        }
    }
    if (target == ALL || target == AUDIO) {
        switch (extension) {

        case "mp3":
            return "audio/mpeg"
        case "wav":
            return "audio/wav"
        case "oga":
            return "audio/ogg"
        case "weba":
            return "audio/webm"
        case "opus":
            return "audio/opus"
        case "aac":
            return "audio/aac"
        }
    }
    if (target == ALL || target == DOC) {
        switch (extension) {
        case "rar":
            return "application/vnd.rar"
        case "zip":
            return "application/zip"
        case "7z":
            return "application/x-7z-compressed"
        case "tar":
            return "application/x-tar"
        case "gz":
            return "application/gzip"
        case "bz":
            return "application/x-bzip"
        case "bz2":
            return "application/x-bzip2"
        case "csv":
            return "text/csv"
        case "ics":
            return "text/calendar"
        case "doc":
            return "application/msword"
        case "docx":
            return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        case "obp":
            return "application/vnd.oasis.opendocument.presentation"
        case "ods":
            return "application/vnd.oasis.opendocument.spreadsheet"
        case "odt":
            return "application/vnd.oasis.opendocument.text"
        case "ppt":
            return "application/vnd.ms-powerpoint"
        case "pptx":
            return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
        }

        if (target == ALL || target == MODEL) {
            switch (extension) {

            case "gltf":
                return "application/json"
            case "obj":
                return "application/object"
            case "dae":
                return "application/collada+xml"
            case "zae":
                return "application/collada-zipped"
            }
        }
    }

    return "application/octet-stream"
}

// Return fileExtension from content-type
func getFileExtension(contentType string) string {

    switch (contentType) {
    case "image/jpeg", "image/pjpeg", "image/x-citrix-jpeg", "image/x-citrix-pjpeg":
        return ".jpg"
    case "image/png":
        return ".png"
    case "image/gif":
        return ".gif"
    case "image/heic":
        return ".heic"
    case "image/heif":
        return ".heif"
    case "image/webp":
        return ".webp"
    case "image/svg":
        return ".svg"
    case "image/bmp":
        return ".bmp"
    }

    return ""
}

func fileExists(filename string) bool {
    info, err := os.Stat(filename)

    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

// check magic number
func checkMagicNumber(bytes []byte) bool {
    magic := [][]byte{
        []byte{ 0xFF, 0xD8, 0xFF },       // jpeg
        []byte{ 0x89, 0x50, 0x4E, 0x47 }, // png
        []byte{ 0x42, 0x4D }, // bmp
        []byte{ 0x52, 0x49, 0x46, 0x46 }, // webp
        []byte{ 0x00, 0x00, 0x00, 0x18 }, // heic
        []byte{ 0x47, 0x49, 0x46, 0x38 }, // gif
        []byte{ 0x25, 0x50, 0x44, 0x46 }, // pdf, ai, weps
        []byte{ 0x50, 0x4B, 0x03, 0x04 }, // zip
        []byte{ 0x52, 0x61, 0x72, 0x21 }, // rar
        []byte{ '7', 'z', 0xBC, 0xAF }, // 7z
        []byte{ 0xC5, 0xD0, 0xD3, 0xC6 }, // eps (Adobe)
        []byte{ 0x25, 0x21, 0x50, 0x53 }, // eps (original)
    }
    length := len(magic)

    for i := length - 1; i != -1; i-- {
        for j := len(magic[i]) - 1; j != -1; j-- {
            if magic[i][j] != bytes[j] {
                break
            } else if j == 0 {
                return true
            }
        }
    }

    return false
}

// remove all files in directory dirPath
func clearDir(dirPath string) {
    dir, err := ioutil.ReadDir(dirPath)

    if err != nil {
        logging("path:", dirPath, err.Error())
    }
    for _, di := range dir {
        os.RemoveAll(path.Join([]string{dirPath, di.Name()}...))
    }
}

// return true if err != nil
func dirIsEmpty(path string) (bool) {
    dir, err := os.Open(path)

    if err != nil {
        return true
    }
    fileInfos, err := dir.Readdir(-1)

    if err != nil || len(fileInfos) == 0 {
        return true
    }
    dir.Close()
    emptySubDir := 0

    for i := range(fileInfos) {
        if dir, err = os.Open(path + "/" + fileInfos[i].Name()); err != nil {
            return true
        }

        if _, err = dir.Readdirnames(1); err == io.EOF {
            emptySubDir++
        }
        dir.Close()
    }

    if emptySubDir == len(fileInfos) {
        return true
    }

    return false
}

// upload a picture at a specified path, resize, and convert to jpeg
// if filename is empty, the name will be a number(auto inc)
// size is of type [2]int
func createPicture(folderPath string, contentType string, filename string, size []int, body *[]byte) (string, bool) {
    i := 0

    if !checkMagicNumber((*body)[:4]) {
        return "", false
    }
    extension := getFileExtension(contentType)

    if extension == "" || extension == ".svg" {
        return "", false
    }
    err := os.MkdirAll(folderPath, 0777)

    if err != nil {
        logging(folderPath, err.Error())
        return "", false
    }
    os.Chmod(folderPath, 0777)

// create temp file
    tempFilename := "~raw_" + strconv.Itoa(i)
    for fileExists(folderPath + "/" + tempFilename) {
        i++
        tempFilename = "~raw_" + strconv.Itoa(i)
    }
    file, err := os.Create(folderPath + "/" + tempFilename)

    if err != nil {
        logging("os.Create()", err.Error())
        return "", false
    } else if _, err = file.Write(*body); err != nil {
        file.Close()
        os.RemoveAll(folderPath + "/" + tempFilename)
        return "", false
    }

    if err = file.Close(); err != nil {
        logging("file.Close()", err.Error())
    } else if err = os.Chmod(folderPath + "/" + tempFilename, 0777); err != nil {
       logging(err.Error())
    }

    // create a fileName using integer as name (0.jpg -> 1.jpg -> 2.jpg)
    if filename == "" {
        i = 0
        filename = strconv.Itoa(i)

        for fileExists(folderPath + "/" + filename + ".jpg") {
            i++
            filename = strconv.Itoa(i)
        }
    }

// imagick
    if (extension == ".png" || extension == ".bmp" || extension == ".webp" || extension == ".heif") {
        cmd := exec.Command("convert", folderPath + "/" + tempFilename, folderPath + "/" + filename + ".jpg")
        err := cmd.Run()
        os.RemoveAll(folderPath + "/" + tempFilename)

        if err != nil {
            logging(err.Error())
            return "", false
        }
    } else if (extension == ".gif" || extension == ".heic") {
        cmd := exec.Command("convert", folderPath + "/" + tempFilename + "[0]", folderPath + "/" + filename + ".jpg")
        err := cmd.Run()
        os.RemoveAll(folderPath + "/" + tempFilename)

        if err != nil {
            logging(err.Error())
            return "", false
        }
    } else if (extension != ".jpg") {
        logging("Bad extension: " + extension)
        os.RemoveAll(folderPath + "/" + tempFilename)
        return "", false
    } else {
        if err = os.Rename(folderPath + "/" + tempFilename, folderPath + "/" + filename + ".jpg"); err != nil {
            logging(err.Error())
        }
    }

    if size != nil {
        return ("/" + folderPath + "/" + filename + ".jpg"), resizePicture(folderPath, filename, filename, size)
    }
    return ("/" + folderPath + "/" + filename + ".jpg"), true
}
//convert image.png image.jpg
//convert ah.jpg -resize 584x438^ -gravity center -extent 584x438 ah.jpg
// convert ah.jpg -resize 584x438^ -quality 80 ah.jpg

// imagick resize jpeg, center/crop/quality=80
func resizePicture(folderPath string, baseFile string, newFile string, size []int) bool {
    cmd := exec.Command("convert", folderPath + "/" + baseFile + ".jpg", "-resize", strconv.Itoa(size[0]) + "x" + strconv.Itoa(size[1]) + "^", "-gravity", "center", "-extent", strconv.Itoa(size[0]) + "x" + strconv.Itoa(size[1]), "-quality", "80", folderPath + "/" + newFile + ".jpg")
    err := cmd.Run()

    if err != nil {
        logging(err.Error())
        return false
    }
    return true

    return true
}

// swap files in a directory
func fileSwap(path string, ids *[]int) bool {
    length := len(*ids)
    var err error

    for i := 0; i != length; i++ {
        if i != (*ids)[i] {
            if err = os.Rename(path + strconv.Itoa(i) + ".jpg", path + "temp.jpg"); err != nil {
                logging(err.Error())
            } else if err = os.Rename(path + strconv.Itoa((*ids)[i]) + ".jpg", path + strconv.Itoa(i) + ".jpg"); err != nil {
                logging(err.Error())
            }
            if !fileSwapLoop(path, ids, i, (*ids)[i]) {
                return false
            }
            (*ids)[i] = i
        }
    }

    return true
}

// loop over the pattern (stops when temp is reached)
func fileSwapLoop(path string, ids *[]int, temp int, i int) bool {
    var t int
    var err error

    for i != (*ids)[i] {
        if (*ids)[i] == temp {
// logging("swap", temp, "->", i)
            if err = os.Rename(path + "temp.jpg", path + strconv.Itoa(i) + ".jpg"); err != nil {
                logging(err.Error())
            }
            (*ids)[i] = i
            return true
        } else {
// logging("swap", (*ids)[i], "->", i)
            if err = os.Rename(path + strconv.Itoa((*ids)[i]) + ".jpg", path + strconv.Itoa(i) + ".jpg"); err != nil {
                logging(err.Error())
            }
        }
        t = i
        i = (*ids)[i]
        (*ids)[t] = t
    }
    return false
}
