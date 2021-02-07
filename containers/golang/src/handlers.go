/*
** Made by LUBERT Arnaud, Epitech Student (Promo 2023)
** 06/02/2021 arnaud.lubert@epitech.eu
**
*/
package main

import (
    "encoding/xml"
    "io/ioutil"
    "net/http"
    "bufio"
    "time"
    "os"

/*    _ "github.com/go-sql-driver/mysql"
    "html/template"
    "encoding/json"
    "database/sql"
    "encoding/hex"
    "crypto/md5"
    "strconv"
    "strings"
    "net/url"
    "math"
    */
)

// Handles root request (welcome page)
func home(rw http.ResponseWriter, req *http.Request) {

    if (req.URL.Path == "" || req.URL.Path == "/") {
        file, err := os.Open("html/index.html")

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

        compressed := setHeader(&rw, req, 2, nil, "text/html")
        w, zip := compressedWritter(&rw, compressed)
        (*w).Write(content)
        file.Close()
        freeCompression(&w, zip)
    } else {
        servErr(&rw, http.StatusNotFound)
    }
}

func favicon(rw http.ResponseWriter, req *http.Request) {
    file, err := os.Open("img/favicon/favicon.ico")

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

    compressed := setHeader(&rw, req, 2, nil, "image/vnd.microsoft.icon")
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    file.Close()
    freeCompression(&w, zip)
}

func robots(rw http.ResponseWriter, req *http.Request) {
    file, err := os.Open("robots.txt")

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

    compressed := setHeader(&rw, req, 0, nil, "text/plain")
    w, zip := compressedWritter(&rw, compressed)
    (*w).Write(content)
    file.Close()
    freeCompression(&w, zip)
}

func sitemap(rw http.ResponseWriter, req *http.Request) {
    compressed := setHeader(&rw, req, 0, nil, "application/xml")
    date := time.Now().Format(http.TimeFormat)
    w, zip := compressedWritter(&rw, compressed)
    token := xml.ProcInst{ "xml", []byte("version=\"1.0\" encoding=\"UTF-8\"") }
    nm := xml.Name{"", "url"}
    data := Sitemap{
        UrlsetAttr: "http://www.sitemaps.org/schemas/sitemap/0.9",
        Urls: []SitemapUrl{
            SitemapUrl{nm, BaseUrl, sitemapDate(date), "weekly", 1},
            SitemapUrl{nm, BaseUrl + "/supports", sitemapDate(date), "monthly", 1},
            SitemapUrl{nm, BaseUrl + "/articles", sitemapDate(date), "monthly", 0.9},
            SitemapUrl{nm, BaseUrl + "/blog", sitemapDate(date), "weekly", 0.8},
            SitemapUrl{nm, BaseUrl + "/qui-sommes-nous", sitemapDate(date), "monthly", 0.75},
            SitemapUrl{nm, BaseUrl + "/service-apres-vente", sitemapDate(date), "yearly", 0.72},
            SitemapUrl{nm, BaseUrl + "/infos-contact", sitemapDate(date), "yearly", 0.7},
            SitemapUrl{nm, BaseUrl + "/foire-aux-questions", sitemapDate(date), "weekly", 0.5},
            SitemapUrl{nm, BaseUrl + "/mentions-legales", sitemapDate(date), "yearly", 0.4},
            SitemapUrl{nm, BaseUrl + "/conditions-generales-de-vente", sitemapDate(date), "yearly", 0.4},
            SitemapUrl{nm, BaseUrl + "/connexion-au-compte", sitemapDate(date), "yearly", 0.2},
            SitemapUrl{nm, BaseUrl + "/plan-du-site", sitemapDate(date), "yearly", 0.2}}}

    enc := xml.NewEncoder(*w)
    enc.EncodeToken(token)
    enc.Indent("  ", "    ")

    if err := enc.Encode(&data); err != nil {
    	logging(err.Error())
        rw.WriteHeader(http.StatusInternalServerError) //  :$
        freeCompression(&w, zip)
        return
    }
    freeCompression(&w, zip)
}
