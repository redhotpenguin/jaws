package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    l4g "code.google.com/p/log4go"
    "database/sql"
    _ "github.com/lib/pq"
    "encoding/xml"
)

type ConnectInfo struct {
    XMLName    xml.Name    `xml."databases"`
    Connection []DataConnection `xml:"database"`
} 

var pgConnection string

func main() {

    // read application configuration
    l4g.LoadConfiguration("conf/log.xml")
    defer l4g.Close()

    // read the database config
    LoadDatabaseConfig("conf/database.xml")    

    // map the endpoints to codepaths
    
}

func LoadDatabaseConfig() {
    connectFile, error := os.Open(configFile)
    if error != nil {
        l4g.Error("Failed to open config file %s, err %v", configFile, error)
        return error
    }

    connectDetails, error := ioutil.ReadAll(configFile)
    if error != nil {
        l4g.Error("Failed to read file %s, err %v", configFile, error)
        return error
    }

    error := xml.Unmarshal(connectDetails, ConnectInfo);
    if err != nil {
        l4g.Error("Failed to unmarshal %s: %v", configFile, error)
    }

    
    pgConnection = fmt.Sprintf("user=%s password=%s host=%s dbname=%s", connection.User, connection.Password, connection.Host, connection.Database)
    return nil
}


