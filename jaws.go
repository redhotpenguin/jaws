package main

import (
//    "fmt"
    "io/ioutil"
    "os"
//    "net/http"
//    "strings"
//    "database/sql"
    _ "github.com/lib/pq"
    "encoding/xml"
    glog "github.com/golang/glog"
)

type DataConnection struct {
    xml.Name `xml:"person"`
    DbType   string `xml:"dbtype"`
    Host     string `xml:"host"`
    User     string `xml:"user"`
    Database string `xml:"dbname"`
    Password string `xml:"password"`
}

func main() {

    var configFile = "conf/database.xml"

    connectFile, error := os.Open(configFile)
    if error != nil {
        glog.Error("Failed to open config file %s, err %v", configFile, error)
        return 
    }

    connectDetails, error := ioutil.ReadAll(connectFile)
    if error != nil {
        glog.Error("Failed to read file %s, err %v", configFile, error)
        return 
    }

    v := new(DataConnection)

    xmlError := xml.Unmarshal(connectDetails, &v)
    if xmlError != nil {
        glog.Error("Failed to unmarshal %s: %v", configFile, xmlError)
    }


    glog.Error("Oh hai here's some database %s", v)
    
//    pgConnection = fmt.Sprintf("user=%s password=%s host=%s dbname=%s", connection.User, connection.Password, connection.Host, connection.Database)
    return 
}


