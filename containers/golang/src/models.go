/*
** Made by LUBERT Arnaud, Epitech Student (Promo 2023)
** 06/02/2021 arnaud.lubert@epitech.eu
**
*/
package main

import (
    "encoding/xml"
)
/*
type ShippingContext struct {
    Id        int         `json:"-"`
    Position  int         `json:"position"`
    Visible   bool        `json:"visible"`
    Ref       string      `json:"ref"`
    Title     string      `json:"title"`
    Desc      string      `json:"description"`
    Prices    []PriceData `json:"prices"`
}
*/

type Sitemap struct {
    XMLName   xml.Name `xml:"urlset"`
    UrlsetAttr  string `xml:"xmlns,attr"`
    Urls  []SitemapUrl `xml:""`
}

type SitemapUrl struct {
    XMLName    xml.Name `xml:"url"`
    Loc        string   `xml:"loc"`
    Lastmod    string   `xml:"lastmod"`
    Changefreq string   `xml:"changefreq"`
    Priority   float32  `xml:"priority"`
}
