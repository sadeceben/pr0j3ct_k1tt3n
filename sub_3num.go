package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func control(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func loger(prefix, description string) {
	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	control(err)
	defer file.Close()
	logger := log.New(file, prefix, log.LstdFlags)
	logger.Println(description)
}

func parser(body, domain, search_engine string) (is_it bool) {
	var reg_x string
	if strings.Contains(search_engine, "Google") {
		reg_x = "(http|ftp|https)://([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])?"
	} else if strings.Contains(search_engine, "PassiveDNS") {
		reg_x = "<td>(.*?)</td>"
	}

	is_it = false
	re := regexp.MustCompile(reg_x)
	list := re.FindAllString(body, -1)
	for _, links := range list {
		if strings.Contains(links, domain) {
			if strings.Contains(search_engine, "PassiveDNS") {
				links = strings.Replace(links, "<td>", "", -1)
				links = strings.Replace(links, " [TR]</td>", "", -1)
				fmt.Println("|----> : " + "http://" + links + " OR " + "https://" + links)
				is_it = true
			} else {
				fmt.Println("|----> : " + links)
				is_it = true
			}
		}
	}
	return
}

func GoogleEnum(query string) {
	base_url := "https://google.com/search?q={query}&btnG=Search&hl=en-US&biw=&bih=&gbv=1&start={page_no}&filter=0"
	base_url = strings.Replace(base_url, "{query}", "site:*.*."+query, -1)
	new_url := ""

	loger("RUNNING : ", "GoogleEnum function is runnig")

	for i := 0; i < 100; i++ {
		new_url = strings.Replace(base_url, "{page_no}", strconv.Itoa(i*10), -1)
		resp, err := http.Get(new_url)
		control(err)
		if resp.Status == "404 Not Found" {
			loger("ERROR : ", "GoogleEnum function get request status is : "+resp.Status)
			break
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		control(err)

		if parser(string(body), query, "Google") {
			fmt.Println(i)
		} else if strings.Contains(string(body), "https://www.google.com/recaptcha/api.js") {
			loger("ERROR : ", "GoogleEnum funciton is caught google recaptcha")
			break
		} else {
			loger("ERROR : ", "Last Page")
			break
		}
	}
	loger("STOP : ", "GoogleEnum function is stop")
}

func PassiveDNS(domain string) {
	base_url := "http://ptrarchive.com/tools/search.htm?label=" + domain
	loger("RUNNING : ", "PassiveDNS running")

	resp, err := http.Get(base_url)
	control(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	control(err)
	if parser(string(body), domain, "PassiveDNS") {
		loger("SUCCESFLY : ", "PassiveDNS have finished")
	} else {
		loger("ERROR : ", "cannot runnig PassiveDNS")
	}

}

func main() {
	var domain string
	fmt.Print("Enter your Domain : ")
	fmt.Scan(&domain)

	//	GoogleEnum("firat.edu.tr")
	PassiveDNS(domain)
}
