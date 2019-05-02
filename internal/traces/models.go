package traces

import "encoding/xml"

// Host Model for the hosts table
type Host struct {
	HostID    int    `json:"hostId" db:"host_id"`
	Hostname  string `json:"hostname" db:"hostname"`
	Addr      string `json:"addr" db:"addr"`
	AddrType  string `json:"addrType" db:"addrtype"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

// Port Model for the Port table
type Port struct {
	HostID    int    `json:"hostId" db:"host_id"`
	Protocol  string `json:"protocol" db:"protocol"`
	PortID    string `json:"portId" db:"port_id"`
	State     string `json:"state" db:"state"`
	Reason    string `json:"reason" db:"reason"`
	Name      string `json:"name" db:"name"`
	StartTime string `json:"startTime" db:"start_time"`
}

// Run Model for the xml nmap results
type Run struct {
	XMLName xml.Name `xml:"run"`
	Text    string   `xml:",chardata"`
	Host    []struct {
		Text      string `xml:",chardata"`
		Starttime string `xml:"starttime,attr"`
		Endtime   string `xml:"endtime,attr"`
		Status    struct {
			Text      string `xml:",chardata"`
			State     string `xml:"state,attr"`
			Reason    string `xml:"reason,attr"`
			ReasonTtl string `xml:"reason_ttl,attr"`
		} `xml:"status"`
		Address struct {
			Text     string `xml:",chardata"`
			Addr     string `xml:"addr,attr"`
			Addrtype string `xml:"addrtype,attr"`
		} `xml:"address"`
		Hostnames struct {
			Text     string `xml:",chardata"`
			Hostname struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				Type string `xml:"type,attr"`
			} `xml:"hostname"`
		} `xml:"hostnames"`
		Ports struct {
			Text string `xml:",chardata"`
			Port []struct {
				Text     string `xml:",chardata"`
				Protocol string `xml:"protocol,attr"`
				Portid   string `xml:"portid,attr"`
				State    struct {
					Text      string `xml:",chardata"`
					State     string `xml:"state,attr"`
					Reason    string `xml:"reason,attr"`
					ReasonTtl string `xml:"reason_ttl,attr"`
				} `xml:"state"`
				Service struct {
					Text   string `xml:",chardata"`
					Name   string `xml:"name,attr"`
					Method string `xml:"method,attr"`
					Conf   string `xml:"conf,attr"`
				} `xml:"service"`
			} `xml:"port"`
		} `xml:"ports"`
		Times struct {
			Text   string `xml:",chardata"`
			Srtt   string `xml:"srtt,attr"`
			Rttvar string `xml:"rttvar,attr"`
			To     string `xml:"to,attr"`
		} `xml:"times"`
	} `xml:"host"`
}
