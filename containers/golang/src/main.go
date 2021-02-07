/*
** Made by LUBERT Arnaud, Epitech Student (Promo 2023)
** 06/02/2021 arnaud.lubert@epitech.eu
**
*/
package main

import (
    "golang.org/x/crypto/acme/autocert"
    _ "github.com/go-sql-driver/mysql"
//    "golang.org/x/net/websocket"
//    "html/template"
    "database/sql"
    "crypto/tls"
//    "io/ioutil"
    "net/http"
    "strconv"
    "runtime"
    "strings"
    "os/exec"
//    "bufio"
    "time"
    "log"
    "os"
)

var logger *log.Logger
var server *http.Server
var db *sql.DB
var cert tls.Certificate
var Domain, BaseUrl, MailFooter, logFilePath, SmtpHost, SmtpAddr, MondayApiKey, MondayAppKey string
var MailSender, MailPass, MailSavPass, MailReceiver_1, MailReceiver_2, MailSav, MailRC, MailComm string
var GoogleClientId, GoogleSecret string
var Debug bool

func main() {
    var err error
// get env
    MailReceiver_1 = os.Getenv("MAIL_RECEIVER_1")
    MailReceiver_2 = os.Getenv("MAIL_RECEIVER_2")
    MailComm = os.Getenv("MAIL_RECEIVER_COMM")
    MailSav = os.Getenv("MAIL_SAV")
    MailSavPass = decodeBase64(os.Getenv("MAIL_SAV_PASS"))
    MailRC = os.Getenv("MAIL_RECEIVER_RC")
    MailSender = os.Getenv("MAIL_SENDER")
    MailPass = decodeBase64(os.Getenv("MAIL_PASS"))
    MondayApiKey = os.Getenv("MONDAY_API_KEY")
    MondayAppKey = os.Getenv("MONDAY_APP_KEY")
    GoogleClientId = os.Getenv("GOOGLE_CLI_ID")
    GoogleSecret = os.Getenv("GOOGLE_SECRET")
    Domain = os.Getenv("DOMAIN_NAME")
    BaseUrl = os.Getenv("URL_SCHEME")
    SmtpHost = os.Getenv("SMTP_HOST")
    SmtpAddr = os.Getenv("SMTP_ADDR")

    if os.Getenv("DEBUG") != "" {
        Debug = true
    }
    if strings.Contains(BaseUrl, "http://") {
        Domain = "localhost"
    }
    BaseUrl += Domain

    if BaseUrl == "" {
        panic("Missing env variable URL_SCHEME or DEBUG")
    }

// init loggers
    if Debug {
        logFilePath = "../containers_logs/golang/" + time.Now().String()[:19] + ".txt" // local execution
    } else {
        logFilePath = "/logs/" + time.Now().String()[:19] + ".txt" // Docker
    }
    logFile, err := os.Create(logFilePath)

    if err != nil {
        panic(err.Error())
    }
    logger = log.New(logFile, "", log.LstdFlags)

    //connectToDB()
    initHandlers()

    now := time.Now() // sync clearLeaks on 4 am
    Y, M, d := now.Date()
    H, _, _ := now.Clock()
    target := time.Date(Y, M, d, 4, 0, 0, 0, now.Location())

    if H >= 4 {
        target = target.AddDate(0, 0, 1)
    }
    time.AfterFunc(target.Sub(now), clearLeaks) // first call at 4 am
    defer db.Close()

    if (strings.Contains(BaseUrl, "https://")) {

        // Subdomains won't work with autocert, instead use an env var to load the cert key pair
        certManager := autocert.Manager{
    		Prompt: autocert.AcceptTOS,
    		Cache:  autocert.DirCache("/certificates"),
            HostPolicy: autocert.HostWhitelist(Domain, "www." + Domain),
            Email: os.Getenv("MAINTENER_EMAIL"),
        //    RenewBefore: time.Hour * 24 * 30 = default
    	}

        cfg := &tls.Config{GetCertificate: certManager.GetCertificate, NextProtos: []string{"http/1.1"}, ServerName: "signatix" }
        server = &http.Server{
            Addr: ":https",
            Handler: nil, // default http mux
            TLSConfig: cfg,
        }
        go http.ListenAndServe(":80", certManager.HTTPHandler(http.HandlerFunc(httpHandler)))
        //go http.ListenAndServe(":80", http.HandlerFunc(httpHandler))
        logging("Server ready", BaseUrl)
    	log.Fatal(server.ListenAndServeTLS("", "")) // letsencrypt ssl certificate
    } else {
        server = &http.Server{ Addr: ":80" }
        logging("Server ready", BaseUrl)
        log.Fatal(server.ListenAndServe())
    }
}
/*
func GetCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
    return &cert, nil
}*/

func httpHandler(rw http.ResponseWriter, req *http.Request) {

    if len(req.URL.Path) > 4 && req.URL.Path[:4] == "/api" {
        api(rw, req)
//    } else if len(req.URL.Path) > 10 && req.URL.Path[:10] == "/fichiers/" {
//        sendCompressedFiles(rw, req)
    } else {
        http.Redirect(rw, req, BaseUrl + req.URL.Path, http.StatusPermanentRedirect)
    }
}

func connectToDB() {
    var err error

    for db == nil || err != nil {
        if Debug {
            db, err = sql.Open("mysql", os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("LOCAL_IP") + ":" + os.Getenv("MYSQL_PORT") + ")/" + os.Getenv("MYSQL_DATABASE")) // local execution
        } else {
            db, err = sql.Open("mysql", os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(mysql:3306)/" + os.Getenv("MYSQL_DATABASE")) // Docker
        }

        if err != nil {
            logging(err.Error())
            continue
        }

        if err = db.Ping(); err != nil {
            time.Sleep(5 * time.Second)
            err = db.Ping()
        }
        try := 0

        for err != nil && try != 5 {
            logging("Waiting for database... (try: " + strconv.Itoa(try + 1) + ")")
            time.Sleep(6 * time.Second) // 6 sec
            err = db.Ping()
            try++
        }

        if err != nil {
            logging("Could not connect to the database :(")
            db.Close()
        }
    }
    db.SetMaxOpenConns(235) // max 250
    db.SetMaxIdleConns(240)
    db.SetConnMaxLifetime(10 * time.Second)
    db.SetConnMaxIdleTime(5 * time.Second)
}

// warning, no DDOS protection, if wanted use go channels with limit of 100
// if simultanioous reqNumber > 130, server will reboot
func checkLeaks() {
    var threads int
    var name string

    rows, err := db.Query("SHOW STATUS WHERE `variable_name` = 'Threads_connected'")

    if err != nil {
        logging(err.Error())
    //    reboot() this can be useful someday if leaks come to fast
    } else {
        if rows.Next() {

            if err = rows.Scan(&name, &threads); err != nil {
                logging(err.Error())
            }
        } else {
            logging("Next()")
        }
        rows.Close()
    }

    if threads > 130 { // default max: 150
        logging("too much open connections to DB, maybe leaks, rebooting...")
    //    reboot()
    } else {
        time.AfterFunc(time.Minute * 15, checkLeaks)
    }
}

// clear leaks once a week
func clearLeaks() {
    /*
    backupDatabase()

    updatePasswordMap = make(map[string]int) // force garbage collection
    newDevicesMap = make(map[string]DeviceTrusted)
    requestPaymentMap = make(map[string]int)

    if _, err := db.Exec("DELETE FROM sessions WHERE timestamp < ?", time.Now().Unix()); err != nil {
        logging(err.Error())
    }
    socketFactory.ClearLeaks() // clear expired sessions
    time.AfterFunc(time.Minute * 15, checkLeaks) // check nbr of simult. connections to DB

// clear empty orders
    var id, userId int
    var creditAmount float64

    date := time.Now().Add(-time.Hour * 4).String()[:19]
    rows, err := db.Query("SELECT id, member_id, credit FROM orders WHERE complete = 0 AND date < ?", date)

    if err != nil {
        logging(err.Error())
        return
    }

    for rows.Next() {
        if err = rows.Scan(&id, &userId, &creditAmount); err != nil {
            logging(err.Error())
            rows.Close()
            return
        }

        if creditAmount > 0.0 { // create credit back
            var credit CreditData

            rowsSec, err := db.Query("SELECT refClient FROM users WHERE id = ?", userId)

            if err != nil {
                rows.Close()
                logging(err.Error())
                logging("Avoir perdu:", formatPrice(creditAmount), "client:", strconv.Itoa(userId))
            } else {
                if rowsSec.Next() {
                    logging("Avoir perdu:", formatPrice(creditAmount), "client:", strconv.Itoa(userId))
                } else {
                    if err = rows.Scan(&credit.RefUser); err != nil {
                        logging(err.Error())
                        logging("Avoir perdu:", formatPrice(creditAmount), "client:", strconv.Itoa(userId))
                    } else {
                        credit.UserId = userId
                        credit.RawAmount = creditAmount
                        credit.Context = "Suppression commande inachevée #" + strconv.Itoa(id)
                        ret := createCredit(&credit)

                        if ret < 200 || ret > 299 {
                            logging("Un avoir n'a pas pu être recré lors de la suppresion d'une commande oubliée. Montant:", formatPrice(creditAmount), "UserId:", strconv.Itoa(userId))
                            go sendSimpleMail("Erreur avoir client", "Un problème est survenu, une commande oubliée a été supprimer mais l'avoir utilisé n'a pas pû être rétablit.<br>Un avoir du montant: " + formatPrice(creditAmount) + " doit être crée manuellement pour le client #" + strconv.Itoa(userId), []string{MailReceiver_1, MailReceiver_2}, MAIL_DEFAULT)
                        }
                    }
                }
                rowsSec.Close()
            }
        }
        deleteOrder(id)
    }
    rows.Close()
    time.AfterFunc(time.Hour * 24, clearLeaks) // callback loop
    */
}

func backupDatabase() {
    date := time.Now().String()[:10]

    if Debug { // outside container
        cmd := exec.Command("mysqldump", "-u", os.Getenv("MYSQL_USER"), "-p" + os.Getenv("MYSQL_PASSWORD"), "-P", os.Getenv("MYSQL_PORT"), "--no-tablespaces", "--protocol=TCP", os.Getenv("MYSQL_DATABASE"), ">", "../containers_backups/mysql+" + date + ".sql")
        err := cmd.Run()

        if err != nil {
            logging(err.Error())
        }
    } else {
        cmd := exec.Command("mysqldump", "-u", os.Getenv("MYSQL_USER"), "-p" + os.Getenv("MYSQL_PASSWORD"), "-h", "mysql", "-P", "3306", "--no-tablespaces", "--protocol=TCP", os.Getenv("MYSQL_DATABASE"))
        file, err := os.Create("/backups-sql/" + date + ".sql")

        if err != nil {
            logging(err.Error())
            return
        }
        cmd.Stdout = file

        if err = cmd.Start(); err != nil {
            logging(err.Error())
            file.Close()
            return
        }
        if err = cmd.Wait(); err != nil {
            logging(err.Error())
            file.Close()
            return
        }
        file.Close()
    }
    os.Remove("../containers_backups/mysql+" + time.Now().Add(-360 * time.Hour).String()[:10] + ".sql") // delete 15 days old backup
}

func logging(args ...interface{}) {
	av := args
    pc, file, line, success := runtime.Caller(1)

	if success {
        index := strings.LastIndexByte(file, '/')

        if index != -1 {
            file = file[index + 1:]
        }
    	log.Printf("%s() %s:%d %v", runtime.FuncForPC(pc).Name()[5:], file, line, av)
        logger.Printf("%s() %s:%d %v", runtime.FuncForPC(pc).Name()[5:], file, line, av)
	} else {
		log.Print(av)
        logger.Print(av)
	}
}

func initHandlers() {
    http.HandleFunc("/", home)
    http.HandleFunc("/favicon.ico", favicon)
    http.HandleFunc("/robots.txt", robots)
    http.HandleFunc("/sitemap.xml", sitemap)
    http.HandleFunc("/js/", js)
    http.HandleFunc("/css/", css)
    http.HandleFunc("/api/", api)
    http.HandleFunc("/img/", images)
    http.HandleFunc("/models/", models)
    //http.HandleFunc("/articles", articles)
}
