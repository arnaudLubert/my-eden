/*
** Made by LUBERT Arnaud, Epitech Student (Promo 2023)
** 06/02/2021 arnaud.lubert@epitech.eu
**
*/
package main

import (
    _ "github.com/go-sql-driver/mysql"
    "net/http"
    "strings"
    "html"
)

type ApiFunction func(*http.ResponseWriter, *http.Request)

func api(rw http.ResponseWriter, req *http.Request) {
    var function ApiFunction // (void *)
    path := html.EscapeString(req.URL.Path)[4:]
    index := strings.IndexByte(path[1:], '/')

    if index != -1 {
        path = path[:1 + index]
    }

    switch req.Method{
    case "GET":
        switch path{

        case "/ping": function = apiPing
        default: function = apiError
        }
    /*case "POST":
        switch path{
        case "/login": function = apiLogin
        default: function = apiError
        }
    case "PUT":
        switch path{
        case "/reload-ressource": function = apiReloadRessource
        default: function = apiError
        }
    case "DELETE":
        switch path{
        case "/file": function = apiDeleteFile
        default: function = apiError
        }*/
    default: function = apiError
    }
    function(&rw, req)
}

func apiError(w *http.ResponseWriter, _ *http.Request) {
    http.Error(*w, "{\n  \"Repsonse\": \"Unknown Request\"\n}", http.StatusNotFound)
}



/*-----------------------------  SPECIAL CASES  ------------------------------*/

// return client address:port
func apiPing(rw *http.ResponseWriter, req *http.Request) {
    (*rw).Write([]byte(getAddress(req.RemoteAddr)))
}
