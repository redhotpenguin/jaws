package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/hoisie/web"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
)

type DataConnection struct {
	xml.Name `xml:"person"`
	DbType   string `xml:"dbtype"`
	Host     string `xml:"host"`
	User     string `xml:"user"`
	Database string `xml:"dbname"`
	Password string `xml:"password"`
}

type DsaQuery struct {
	Id             int
	Query          string
	SourceFile     string
	SourceFunction string
	SourceLine     int
	IssueTicket    string
	ApiCall        string
}

func list(val string) string {

	// grab db config options
	var configFile = "conf/database.xml"
	connectFile, error := os.Open(configFile)
	if error != nil {
		exception := fmt.Sprintf("Failed to open config file %s, err %v", configFile, error)
		glog.Error(exception)
		return exception
	}

	connectDetails, error := ioutil.ReadAll(connectFile)
	if error != nil {
		exception := fmt.Sprintf("Failed to read file %s, err %v", configFile, error)
		glog.Error(exception)
		return exception
	}

	// read the xml
	dbConfig := new(DataConnection)
	xmlError := xml.Unmarshal(connectDetails, &dbConfig)
	if xmlError != nil {
		exception := fmt.Sprintf("Failed to unmarshal %s: %v", configFile, xmlError)
		glog.Error(exception)
		return exception
	}

	// open up the database connection
	connectString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Database,
		dbConfig.Host)

	db, err := sql.Open("postgres", connectString)
	if err != nil {
		exception := fmt.Sprintf("failed db open: %s", err)
		glog.Error(exception)
		return exception
	} else {
		glog.Info("Connected to ", connectString)
	}
	defer db.Close()

	rows, qerr := db.Query("SELECT * FROM gm3_queries")
	if qerr != nil {
		exception := fmt.Sprintf("dbQuery failed %s", qerr)
		glog.Error(exception)
		return exception
	}

	var dsaResults []DsaQuery
	for rows.Next() {

		var Row DsaQuery
		if err := rows.Scan(
			&Row.Id,
			&Row.Query,
			&Row.SourceFile,
			&Row.SourceFunction,
			&Row.SourceLine,
			&Row.IssueTicket,
			&Row.ApiCall,
		); err != nil {

			exception := fmt.Sprintf("query failure: %s", err)
			glog.Error(exception)
			return exception
		}
		// map the query result into a struct for json translation
		dsaResult := DsaQuery{
			Id:             Row.Id,
			SourceFile:     Row.SourceFile,
			SourceFunction: Row.SourceFunction,
			SourceLine:     Row.SourceLine,
			IssueTicket:    Row.IssueTicket,
			ApiCall:        Row.ApiCall,
			Query:          Row.Query,
		}

		// glog.Info(fmt.Sprintf("Result is %s"), dsaResult)

		// add this result onto the slice (array)
		dsaResults = append(dsaResults, dsaResult)
	}

	if err := rows.Err(); err != nil {
		glog.Error(err)
	}

	json, err := json.Marshal(dsaResults)
	return fmt.Sprintf("%s", json)
}

func main() {

	// parse command line arguments
	flag.Parse()
	web.Get("/(list)", list)
	//web.Get("/(.*)", list)
	//    web.SetLogger(glog)
	web.Run("0.0.0.0:9999")
}
