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

// log file path
const log_f_path string = "logs/logs.txt"

// GoogleEnum Variable
const ggl_name string = "Google"
const google_url string = "https://google.com/search?q={query}&btnG=Search&hl=en-US&biw=&bih=&gbv=1&start={page_no}&filter=0"
const google_re string = "https://www.google.com/recaptcha/api.js"

// PassiveDNS Variable
const passive_url string = "http://ptrarchive.com/tools/search.htm?label="
const passive_name string = "PassiveDNS"

type SearchEngine interface {
	GoogleEnum()
	PassiveDNS()
	control()
	loger()
	parser()
}

type Engine struct {
}

func (e Engine) control(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (e Engine) loger(prefix, description string) {
	file, err := os.OpenFile(log_f_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	e.control(err)
	defer file.Close()
	logger := log.New(file, prefix, log.LstdFlags)
	logger.Println(description)
}

func (e Engine) parser(body, domain, search_engine string) (is_it bool) {
	var reg_x string
	if strings.Contains(search_engine, ggl_name) {
		reg_x = "(http|ftp|https)://([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])?"
	} else if strings.Contains(search_engine, passive_name) {
		reg_x = "<td>(.*?)</td>"
	}

	is_it = false
	re := regexp.MustCompile(reg_x)
	list := re.FindAllString(body, -1)
	for _, links := range list {
		if strings.Contains(links, domain) {
			if strings.Contains(search_engine, passive_name) {
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

func (e Engine) GoogleEnum(query string) {
	base_url := google_url
	base_url = strings.Replace(base_url, "{query}", "site:*.*."+query, -1)
	new_url := ""

	e.loger("RUNNING : ", "GoogleEnum function is runnig")

	for i := 0; i < 100; i++ {
		new_url = strings.Replace(base_url, "{page_no}", strconv.Itoa(i*10), -1)
		resp, err := http.Get(new_url)
		e.control(err)
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		e.control(err)

		if e.parser(string(body), query, ggl_name) {
			fmt.Println(i)
		} else if strings.Contains(string(body), google_re) {
			e.loger("ERROR : ", "GoogleEnum funciton is caught google recaptcha")
			break
		} else {
			break
		}
	}
	e.loger("SUCCESFLY : ", "GoogleEnum function have finished")
}

func (e Engine) PassiveDNS(domain string) {
	base_url := passive_url + domain
	e.loger("RUNNING : ", "PassiveDNS running")

	resp, err := http.Get(base_url)
	e.control(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	e.control(err)
	if e.parser(string(body), domain, passive_name) {
		e.loger("SUCCESFLY : ", "PassiveDNS have finished")
	} else {
		e.loger("ERROR : ", "cannot runnig PassiveDNS")
	}

}

func main() {
	var domain string
	engine := Engine{}
	fmt.Print("Enter your Domain : ")
	fmt.Scan(&domain)

	engine.PassiveDNS(domain)
	engine.GoogleEnum(domain)
}
