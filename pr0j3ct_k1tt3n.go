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
)

// Sleep Flag
var ggl_flag bool = true
var pass_flag bool = true
var yahoo_flag bool = true

// log file path
const log_f_path string = "logs/logs.txt"

// GoogleEnum Variable
const ggl_name string = "Google"
const google_url string = "https://google.com/search?q={query}&btnG=Search&hl=en-US&biw=&bih=&gbv=1&start={page_no}&filter=0"
const google_re string = "https://www.google.com/recaptcha/api.js"

// PassiveDNS Variable
const passive_name string = "PassiveDNS"
const passive_url string = "http://ptrarchive.com/tools/search.htm?label="

// Yahoo Variable
const yahoo_name = "Yahoo"
const yahoo_url  = "https://search.yahoo.com/search?p={query}&b={page_no}"

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

	is_it = false
	re := regexp.MustCompile("")
	if strings.Contains(search_engine, ggl_name) {
		re = regexp.MustCompile("(http|ftp|https)://([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])?")
	} else if strings.Contains(search_engine, passive_name) {
		re = regexp.MustCompile("<td>(.*?)</td>")
	} else if strings.Contains(search_engine, yahoo_name) {
		re = regexp.MustCompile(`<span class="txt"><span class=" cite fw-xl fz-15px">(.*?)</span>`)
	}

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
	fmt.Println("Parser Exit")
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
			fmt.Println(i)
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

func ( Yahoo Engine ) Yahoo(domain string) {

        base_url := strings.Replace(yahoo_url, "{query}", "*.*." + domain, -1)
        new_url := ""

        Yahoo.loger("RUNNING : ", "Yahoo function is runnig")

        for i := 0; i < 100; i++ {
                new_url = strings.Replace(base_url, "{page_no}", strconv.Itoa(i*10), -1)
                defer func() {
                        yahoo_flag = false
                }()
                if Yahoo.parser(Yahoo.Meow(new_url), domain, yahoo_name) {
                        fmt.Println(i)
                } else if strings.Contains(Yahoo.Meow(new_url), google_re) {
                        Yahoo.loger("ERROR : ", "Yahoo funciton is caught google recaptcha")
                        break
                } else {
                        break
                }
        }
        Yahoo.loger("SUCCESFLY : ", "Yahoo function have finished")
}

func main() {
	var domain string

	engine := Engine{}
	fmt.Print("Enter your Domain : ")
	fmt.Scan(&domain)

//	go engine.GoogleEnum(domain)
//	go engine.PassiveDNS(domain)

	go engine.Yahoo(domain)

	for yahoo_flag {
		time.Sleep(0)
	}

}
