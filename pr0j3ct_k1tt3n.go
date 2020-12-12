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
	"time"
	"net/url"
)

// Sleep Flag
var ggl_flag bool = true
var pass_flag bool = true

// log and output file path
const log_f_path string = "logs/logs.txt"
const output_f_path string = "output/output.txt"

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
	writer()
	reader()
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

func (e Engine) writer(link string) {
	file, err := os.OpenFile(output_f_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	e.control(err)

	defer file.Close()

	if link != "" {
	_, err = file.WriteString(link + "\n")
	e.control(err)
	}
}

func (e Engine) reader() {
	content, err := ioutil.ReadFile(output_f_path)
	e.control(err)
	fmt.Println(string(content))
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
				e.writer(links)

				is_it = true
			} else {
				conv_link, _ := url.Parse(links)
				e.writer(conv_link.Host)
				fmt.Println("|----> : " + conv_link.Host)
				is_it = true
			}
		}
	}
	return
}

func (e Engine) Meow(url string) string {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	e.control(err)

	req.Header.Set("User-Agent", "Sadeceben_Kitten_Bot/3.1")

	resp, err := client.Do(req)
	e.control(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	e.control(err)

	return string(body)
}

func (Google Engine) GoogleEnum(query string) {

	base_url := strings.Replace(google_url, "{query}", "site:*.*."+query, -1)
	new_url := ""

	Google.loger("RUNNING : ", "GoogleEnum function is runnig")

	for i := 0; i < 100; i++ {
		new_url = strings.Replace(base_url, "{page_no}", strconv.Itoa(i*10), -1)
		defer func() {
			ggl_flag = false
		}()
		if Google.parser(Google.Meow(new_url), query, ggl_name) {
		} else if strings.Contains(Google.Meow(new_url), google_re) {
			Google.loger("ERROR : ", "GoogleEnum funciton is caught google recaptcha")
			break
		} else {
			break
		}
	}
	Google.loger("SUCCESFLY : ", "GoogleEnum function have finished")
}

func (Passive Engine) PassiveDNS(domain string) {
	base_url := passive_url + domain
	Passive.loger("RUNNING : ", "PassiveDNS running")

	defer func() {
		pass_flag = false
	}()

	if Passive.parser(Passive.Meow(base_url), domain, passive_name) {
		Passive.loger("SUCCESFLY : ", "PassiveDNS have finished")
	} else {
		Passive.loger("ERROR : ", "cannot runnig PassiveDNS")
	}

}

func main() {
	var domain string

	engine := Engine{}
	fmt.Print("Enter your Domain : ")
	fmt.Scan(&domain)

	go engine.GoogleEnum(domain)
	go engine.PassiveDNS(domain)

	for ggl_flag || pass_flag {
		time.Sleep(0)
	}

//	engine.reader()
}
