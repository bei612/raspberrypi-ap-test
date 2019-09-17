package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"io"
	syslog "log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type wifiCollector struct {
	Vendor string
	Link   *prometheus.Desc
	Level  *prometheus.Desc
	Noise  *prometheus.Desc
	Nwid   *prometheus.Desc
	Crypt  *prometheus.Desc
	Frag   *prometheus.Desc
	Retry  *prometheus.Desc
	Misc   *prometheus.Desc
	Rx_bytes *prometheus.Desc
	Tx_bytes *prometheus.Desc
	Rx_error *prometheus.Desc
	Tx_error *prometheus.Desc
	Rx_drop *prometheus.Desc
	Tx_drop *prometheus.Desc
}

func init() {
	syslog.SetPrefix("TRACE: ")
	syslog.SetFlags(syslog.Ldate | syslog.Lmicroseconds | syslog.Llongfile)
}

func sendCommond(command string) (string, string) {
	cmd := exec.Command("/bin/bash","-c",command)
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()
	outstring := out.String()
	errstring := err.String()
	return outstring, errstring
}

func getIfconfig() (rx_bytes int, tx_bytes int, rx_error int, tx_error int, rx_drop int, tx_drop int) {
	result,_ := sendCommond("ifconfig wlan0 | grep -E 'packets|error'")
	result = strings.Trim(result," ")
	resultslice :=strings.Split(result,"\n")
	rx_rate_rst := strings.Trim(resultslice[0]," ")
	rx_error_rst := strings.Trim(resultslice[1]," ")
	tx_rate_rst := strings.Trim(resultslice[2]," ")
	tx_error_rst := strings.Trim(resultslice[3]," ")
	rx_bytes,_ = strconv.Atoi(strings.Split(rx_rate_rst," ")[5])
	tx_bytes,_  = strconv.Atoi(strings.Split(tx_rate_rst," ")[5])
	rx_error,_  = strconv.Atoi(strings.Split(rx_error_rst," ")[2])
	tx_error,_  = strconv.Atoi(strings.Split(tx_error_rst," ")[2])
	rx_drop,_  = strconv.Atoi(strings.Split(rx_error_rst," ")[5])
	tx_drop,_ = strconv.Atoi(strings.Split(tx_error_rst," ")[5])
	return
}

func newWifiCollector(vendor string) *wifiCollector {
	return &wifiCollector{
		Vendor: vendor,
		Link:   prometheus.NewDesc("quality_of_reception", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Level:  prometheus.NewDesc("signal_strength", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Noise:  prometheus.NewDesc("silence_level", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Nwid:   prometheus.NewDesc("Discarded_packets_nwid", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Crypt:  prometheus.NewDesc("Discarded_packets_crypt", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Frag:   prometheus.NewDesc("Discarded_packets_frag", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Retry:  prometheus.NewDesc("Discarded_packets_retry", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Misc:   prometheus.NewDesc("Discarded_packets_msic", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Rx_bytes:  prometheus.NewDesc("received_Bytes", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Tx_bytes:  prometheus.NewDesc("sent_Bytes", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Rx_error:  prometheus.NewDesc("received_error", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Tx_error:  prometheus.NewDesc("sent_error", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Rx_drop:  prometheus.NewDesc("received_drop", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
		Tx_drop:  prometheus.NewDesc("sent_drop", "", []string{"ip"}, prometheus.Labels{"vendor": vendor}),
	}
}

func (c *wifiCollector) getIP() string {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					//fmt.Printf(ipnet.IP.String())
					if strings.HasPrefix(ipnet.IP.String(), "192.168.100.") {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return "192.168.100.X"
}

func (c *wifiCollector) getStatistics() (LinkByIP map[string]int, LevelByIP map[string]int, NoiseByIP map[string]int, NwidByIP map[string]int, CryptByIP map[string]int, FragByIP map[string]int, RetryByIP map[string]int, MiscByIP map[string]int, Rx_bytesByIP map[string]int, Tx_bytesByIP map[string]int, Rx_errorByIP map[string]int, Tx_errorByIP map[string]int, Rx_dropByIP map[string]int, Tx_dropByIP map[string]int) {
	selfip := c.getIP()
	cmd := exec.Command("cat", "/proc/net/wireless") //, "|", "tail", "-n1")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	outstring := out.String()
	length := strings.Split(outstring, "|")
	if strings.Contains(outstring,"wlan0"){
		outstring = length[len(length)-1]
		outsplit := strings.Split(outstring, "  ")
		linkbyip, _ := strconv.Atoi(strings.Replace(strings.Trim(outsplit[1], " "), ".", "", -1))
		levelbyip, _ := strconv.Atoi(strings.Replace(strings.Trim(outsplit[2], " "), ".", "", -1))
		noisebyip, _ := strconv.Atoi(strings.Trim(outsplit[3], " "))
		nwidbyip, _ := strconv.Atoi(strings.Trim(outsplit[4], " "))
		cryptbyip, _ := strconv.Atoi(strings.Trim(outsplit[5], " "))
		fragbyip, _ := strconv.Atoi(strings.Trim(outsplit[6], " "))
		retrybyip, _ := strconv.Atoi(strings.Trim(outsplit[7], " "))
		miscbyip, _ := strconv.Atoi(strings.Trim(outsplit[8], " "))
		LinkByIP = map[string]int{
			selfip: linkbyip,
		}
		LevelByIP = map[string]int{
			selfip: levelbyip,
		}
		NoiseByIP = map[string]int{
			selfip: noisebyip,
		}
		NwidByIP = map[string]int{
			selfip: nwidbyip,
		}
		CryptByIP = map[string]int{
			selfip: cryptbyip,
		}
		FragByIP = map[string]int{
			selfip: fragbyip,
		}
		RetryByIP = map[string]int{
			selfip: retrybyip,
		}
		MiscByIP = map[string]int{
			selfip: miscbyip,
		}
		rx_bytesyip,tx_bytesyip,rx_erroryip,tx_erroryip,rx_dropyip,tx_dropyip := getIfconfig()
		Rx_bytesByIP = map[string]int{
			selfip: rx_bytesyip,
		}
		Tx_bytesByIP = map[string]int{
			selfip: tx_bytesyip,
		}
		Rx_errorByIP = map[string]int{
			selfip: rx_erroryip,
		}
		Tx_errorByIP = map[string]int{
			selfip: tx_erroryip,
		}
		Rx_dropByIP = map[string]int{
			selfip: rx_dropyip,
		}
		Tx_dropByIP = map[string]int{
			selfip: tx_dropyip,
		}
	}
	return
}

func (c *wifiCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Link
	ch <- c.Level
	ch <- c.Noise
	ch <- c.Nwid
	ch <- c.Crypt
	ch <- c.Frag
	ch <- c.Retry
	ch <- c.Misc
	ch <- c.Rx_bytes
	ch <- c.Tx_bytes
	ch <- c.Rx_error
	ch <- c.Tx_error
	ch <- c.Rx_drop
	ch <- c.Tx_drop
}

func (c *wifiCollector) Collect(ch chan<- prometheus.Metric) {
	LinkByIP, LevelByIP, NoiseByIP, NwidByIP, CryptByIP, FragByIP, RetryByIP, MiscByIP, Rx_bytesByIP, Tx_bytesByIP, Rx_errorByIP, Tx_errorByIP, Rx_dropByIP,Tx_dropByIP := c.getStatistics()
	for ip, link := range LinkByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Link,
			prometheus.CounterValue,
			float64(link),
			ip,
		)
	}
	for ip, level := range LevelByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Level,
			prometheus.CounterValue,
			float64(level),
			ip,
		)
	}
	for ip, noise := range NoiseByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Noise,
			prometheus.CounterValue,
			float64(noise),
			ip,
		)
	}
	for ip, nwid := range NwidByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Nwid,
			prometheus.CounterValue,
			float64(nwid),
			ip,
		)
	}
	for ip, crypt := range CryptByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Crypt,
			prometheus.CounterValue,
			float64(crypt),
			ip,
		)
	}
	for ip, frag := range FragByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Frag,
			prometheus.CounterValue,
			float64(frag),
			ip,
		)
	}
	for ip, retry := range RetryByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Retry,
			prometheus.CounterValue,
			float64(retry),
			ip,
		)
	}
	for ip, misc := range MiscByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Misc,
			prometheus.CounterValue,
			float64(misc),
			ip,
		)
	}
	for ip, rx_bytes := range Rx_bytesByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Rx_bytes,
			prometheus.CounterValue,
			float64(rx_bytes),
			ip,
		)
	}
	for ip, tx_bytes := range Tx_bytesByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Tx_bytes,
			prometheus.CounterValue,
			float64(tx_bytes),
			ip,
		)
	}
	for ip, rx_error := range Rx_errorByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Rx_error,
			prometheus.CounterValue,
			float64(rx_error),
			ip,
		)
	}
	for ip, tx_error := range Tx_errorByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Tx_error,
			prometheus.CounterValue,
			float64(tx_error),
			ip,
		)
	}
	for ip, rx_drop := range Rx_dropByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Rx_drop,
			prometheus.CounterValue,
			float64(rx_drop),
			ip,
		)
	}
	for ip, tx_drop := range Tx_dropByIP {
		ch <- prometheus.MustNewConstMetric(
			c.Tx_drop,
			prometheus.CounterValue,
			float64(tx_drop),
			ip,
		)
	}
}

var ch1 = make(chan struct{})
var ch2 = make(chan struct{})

func changeWifi(vendor string) {
	if vendor != "idle" {
		passwdDict := make(map[string]string)
		passwdDict["cisco"] = "cisco008"
		passwdDict["arista"] = "arista08"
		passwdDict["aruba"] = "aruba008"
		commands := [14]string{
			"rm -rf /etc/wpa_supplicant/wpa_supplicant.conf",
			"echo 'ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev' > /etc/wpa_supplicant/wpa_supplicant.conf",
			"echo 'update_config=1' >> /etc/wpa_supplicant/wpa_supplicant.conf",
			fmt.Sprintf("wpa_passphrase %s-test %s >> /etc/wpa_supplicant/wpa_supplicant.conf", vendor, passwdDict[vendor]),
			"more /etc/wpa_supplicant/wpa_supplicant.conf",
			"killall wpa_supplicant",
			"rm -rf /var/run/wpa_supplicant/wlan0",
			"systemctl restart wpa_supplicant.service",
			"sleep 2",
			"wpa_supplicant -i wlan0 -c /etc/wpa_supplicant/wpa_supplicant.conf -B",
			"sleep 2",
			"wpa_cli reconfigure -i wlan0",
			"wpa_cli -i wlan0 list_network",
			"wpa_cli -i wlan0 select_network 0",
		}
		for _, cmd := range commands {
			s,_ := sendCommond(cmd)
			fmt.Println(s)
		}
	} else {
		s,_ := sendCommond("wpa_cli -i wlan0 disable_network 0")
		fmt.Println(s)
	}
}

func metrics(vendor string) {
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}
	s.Handler = mux
	if vendor != "idle" {
		fmt.Println(vendor)
		worker := newWifiCollector(vendor)
		reg := prometheus.NewPedanticRegistry()
		reg.MustRegister(worker)
		gatherers := prometheus.Gatherers{
			reg,
		}
		h := promhttp.HandlerFor(gatherers,
			promhttp.HandlerOpts{
				ErrorLog:      log.NewErrorLogger(),
				ErrorHandling: promhttp.ContinueOnError,
			})
		mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
		go s.ListenAndServe()
		fmt.Printf("prometheus exporter: http://localhost:8888/metrics\n")
	} else {
		mux.HandleFunc("/idel", func(w http.ResponseWriter, _ *http.Request) {
			fmt.Fprintf(w, "idle state \n")
		})
		go s.ListenAndServe()
		fmt.Printf("prometheus exporter is in idel state\n")
	}
	<-ch1
	s.Shutdown(context.Background())
	fmt.Printf("shutdown prometheus exporter")
	ch2 <- struct{}{}
}


func restartVendors(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	query := request.URL.Query()
	vendor := query.Get("name")
	changeWifi(vendor)
	ch1 <- struct{}{}
	<-ch2
	go metrics(vendor)
	io.WriteString(writer, vendor)
}

func index() {
	x := http.NewServeMux()
	server := &http.Server{
		Addr:  ":8080",
		Handler: x,
	}
	x.HandleFunc("/vendor", restartVendors)
	server.ListenAndServe()
}

func status(writer http.ResponseWriter, request *http.Request) {
	rst,_ := sendCommond("iw wlan0 link")
	if strings.Contains(rst, "cisco") {
		rst = "cisco"
	} else if strings.Contains(rst, "aruba"){
		rst = "aruba"
	} else if strings.Contains(rst, "arista"){
		rst = "arista"
	} else {
		rst = "idle"
	}
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(writer, rst)
}

func state() {
	x := http.NewServeMux()
	server := &http.Server{
		Addr:  ":8081",
		Handler: x,
	}
	x.HandleFunc("/state", status)
	server.ListenAndServe()
}

func starttest(writer http.ResponseWriter, request *http.Request) {
	selfip,_ := sendCommond("hostname | cut -d'-' -f 2")
	selfip = strings.Trim(selfip,"\n")
	fmt.Println(selfip)
	nextip,_ := strconv.Atoi(selfip)
	if nextip == 30 {
		nextip = 1
	}
	next_ip := strconv.Itoa(nextip + 1)
	go sendCommond(fmt.Sprintf("/usr/local/iperf/bin/iperf3 -c 172.16.100.%s -p 54321 -t 600", next_ip))
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	io.WriteString(writer, next_ip)
}

func trigger() {
	x := http.NewServeMux()
	server := &http.Server{
		Addr:  ":8082",
		Handler: x,
	}
	x.HandleFunc("/trigger", starttest)
	server.ListenAndServe()
}

func main() {
	go metrics("idle")
	go state()
	go trigger()
	index()
}
