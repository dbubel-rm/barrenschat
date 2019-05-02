package handlers

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dbubel/bchat-api/internal/platform/web/internal/platform/db"
	"github.com/dbubel/bchat-api/internal/platform/web"
	"github.com/dbubel/bchat-api/internal/platform/web/internal/traces"
	"github.com/stretchr/testify/assert"
)

var a *web.App
var d *db.SQLite
var l *log.Logger
var XMLTestFile string
var dsn = "../tests/bf.db"

var xmlStr = `<?xml version="1.0" encoding="UTF-8"?>
	<!DOCTYPE nmaprun>
	<?xml-stylesheet href="file:///usr/local/bin/../share/nmap/nmap.xsl" type="text/xsl"?>
	<!-- Nmap 7.70 scan initiated Fri Apr 20 12:36:44 2018 as: nmap -Pn -p80,443,8443,5000,8080 -iL ips.txt -oA nmap.results -vvvvv -->
	<nmaprun scanner="nmap" args="nmap -Pn -p80,443,8443,5000,8080 -iL ips.txt -oA nmap.results -vvvvv" start="1524242204" startstr="Fri Apr 20 12:36:44 2018" version="7.70" xmloutputversion="1.04">
	<scaninfo type="connect" protocol="tcp" numservices="5" services="80,443,5000,8080,8443"/>
	<verbose level="5"/>
	<debugging level="0"/>
	<taskbegin task="Parallel DNS resolution of 40 hosts." time="1524242204"/>
	<taskend task="Parallel DNS resolution of 40 hosts." time="1524242214"/>
	<taskbegin task="Connect Scan" time="1524242214"/>
	<taskend task="Connect Scan" time="1524242218" extrainfo="200 total ports"/>
	<host starttime="1524242214" endtime="1524242216"><status state="up" reason="user-set" reason_ttl="0"/>
	<address addr="81.107.115.203" addrtype="ipv4"/>
	<hostnames>
	<hostname name="cpc123026-glen5-2-0-cust970.2-1.cable.virginm.net" type="PTR"/>
	</hostnames>
	<ports><port protocol="tcp" portid="80"><state state="open" reason="syn-ack" reason_ttl="0"/><service name="http" method="table" conf="3"/></port>
	<port protocol="tcp" portid="443"><state state="open" reason="syn-ack" reason_ttl="0"/><service name="https" method="table" conf="3"/></port>
	<port protocol="tcp" portid="5000"><state state="filtered" reason="no-response" reason_ttl="0"/><service name="upnp" method="table" conf="3"/></port>
	<port protocol="tcp" portid="8080"><state state="filtered" reason="no-response" reason_ttl="0"/><service name="http-proxy" method="table" conf="3"/></port>
	<port protocol="tcp" portid="8443"><state state="filtered" reason="no-response" reason_ttl="0"/><service name="https-alt" method="table" conf="3"/></port>
	</ports>
	<times srtt="178708" rttvar="108359" to="612144"/>
	</host>
	<runstats><finished time="1524242218" timestr="Fri Apr 20 12:36:58 2018" elapsed="13.95" summary="Nmap done at Fri Apr 20 12:36:58 2018; 40 IP addresses (40 hosts up) scanned in 13.95 seconds" exit="success"/><hosts up="40" down="0" total="40"/>
	</runstats>
	</nmaprun>`

// Setup the DB for the test run
func init() {
	l = log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	os.Remove("../tests/bf.db")
	d, _ = db.New(dsn)
	fixtureFiles, _ := filepath.Glob("../../../*.sql")
	sql, _ := ioutil.ReadFile(fixtureFiles[0])
	d.Database.Exec(string(sql))
	dat, _ := ioutil.ReadFile("../tests/nmap-sample.xml")
	XMLTestFile = string(dat)
}

func resetTestDB() {
	os.Remove("../tests/bf.db")
	d, _ = db.New(dsn)
	fixtureFiles, _ := filepath.Glob("../../../*.sql")
	sql, _ := ioutil.ReadFile(fixtureFiles[0])
	d.Database.Exec(string(sql))
}

func TestXMLMarshal(t *testing.T) {
	var host traces.Nmaprun
	data := []byte(xmlStr)
	err := xml.Unmarshal(data, &host)

	assert.Equal(t, 1, len(host.Host), "Number of hosts scanned should be 1")
	assert.Equal(t, "cpc123026-glen5-2-0-cust970.2-1.cable.virginm.net", host.Host[0].Hostnames.Hostname.Name)
	assert.NoError(t, err, "Unmarshal should have no error")
	assert.Equal(t, 5, len(host.Host[0].Ports.Port), "Should have 5 port elements")
}

func TestXMLMarshal40Hosts(t *testing.T) {
	var host traces.Nmaprun
	data := []byte(XMLTestFile)
	err := xml.Unmarshal(data, &host)

	assert.Equal(t, 40, len(host.Host), "Number of hosts scanned should be 40")
	assert.Equal(t, "cpc123026-glen5-2-0-cust970.2-1.cable.virginm.net", host.Host[0].Hostnames.Hostname.Name)
	assert.NoError(t, err, "Unmarshal should have no error")
	assert.Equal(t, 5, len(host.Host[0].Ports.Port), "Should have 5 port elements")
	// TODO: Verify more things
}

func TestXMLEndpoint(t *testing.T) {
	resetTestDB()
	a := API(l, d).(*web.App)
	r := httptest.NewRequest("POST", "/v1/upload.xml", strings.NewReader(xmlStr))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	assert.Equal(t, http.StatusAccepted, w.Code, "Response code should be 202")

	var queryReturn traces.Host
	err := d.Database.QueryRowx("SELECT * FROM hosts WHERE addr=?", "81.107.115.203").StructScan(&queryReturn)

	assert.NoError(t, err, "Query to DB shold not return an error")
	assert.Equal(t, "81.107.115.203", queryReturn.Addr, "Verify IP is in the DB")
	assert.Equal(t, "cpc123026-glen5-2-0-cust970.2-1.cable.virginm.net", queryReturn.Hostname, "Verify hostname is in DB")
	// TODO: Verify more things...
}

func TestXMLEndpointBadPayload(t *testing.T) {
	resetTestDB()
	a := API(l, d).(*web.App)
	badPayload := xmlStr
	badPayload = strings.Replace(badPayload, "hostnames", "bad_xml_tag", 1) // replace one host tag with an unsupported tag

	r := httptest.NewRequest("POST", "/v1/upload.xml", strings.NewReader(badPayload))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code, "Response code should be 400")

	var queryReturn traces.Host

	err := d.Database.QueryRowx("SELECT * FROM hosts WHERE addr=?", "81.107.115.203").StructScan(&queryReturn)
	assert.Error(t, err, "Database should return error")
	assert.Equal(t, sql.ErrNoRows, err, "Error should be no rows returned")
	// TODO: Verify more things...
}

func TestGetHosts(t *testing.T) {
	resetTestDB()
	a := API(l, d).(*web.App)
	r := httptest.NewRequest("POST", "/v1/upload.xml", strings.NewReader(xmlStr))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	assert.Equal(t, http.StatusAccepted, w.Code, "Response code should be 202")

	var queryReturn traces.Host
	err := d.Database.QueryRowx("SELECT * FROM hosts WHERE addr=?", "81.107.115.203").StructScan(&queryReturn)
	assert.NoError(t, err, "Database should not eturn error")

	r = httptest.NewRequest("GET", "/v1/gethosts", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")

	var res []traces.Host
	err = json.NewDecoder(w.Body).Decode(&res)

	assert.NoError(t, err, "Decoding the payload should not have an error")
}

func TestGetPorts(t *testing.T) {
	resetTestDB()
	a := API(l, d).(*web.App)
	r := httptest.NewRequest("POST", "/v1/upload.xml", strings.NewReader(xmlStr))
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)

	assert.Equal(t, http.StatusAccepted, w.Code, "Response code should be 202")

	var queryReturn traces.Host
	err := d.Database.QueryRowx("SELECT * FROM hosts WHERE addr=?", "81.107.115.203").StructScan(&queryReturn)
	assert.NoError(t, err, "Database should not eturn error")

	r = httptest.NewRequest("GET", "/v1/getports/1", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")

	var res []traces.Port
	err = json.NewDecoder(w.Body).Decode(&res)

	assert.NoError(t, err, "Decoding the payload should not have an error")

	assert.Equal(t, 5, len(res), "Should have 5 port records")
}
