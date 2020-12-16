package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"html/template"
)


// Target Domain
var target_domain string

// Sleep Flag
var ggl_flag bool = true
var pass_flag bool = true

// log and output file path
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
	result()
	control()
	loger()
	parser()
	formatter()
}

type Links struct {
	List []string
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


func (e Engine) formatter(content string) []string {
	re := regexp.MustCompile(`(.|\n)*?,`)
	links := re.FindAllString(content,-1)

	return links
}

func (e Engine) parser(body, search_engine string) (is_it bool,domains string) {
	var reg_x string
	domains = ""
	if strings.Contains(search_engine, ggl_name) {
		reg_x = "(http|ftp|https)://([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])?"
	} else if strings.Contains(search_engine, passive_name) {
		reg_x = "<td>(.*?)</td>"
	}

	is_it = false
	re := regexp.MustCompile(reg_x)
	list := re.FindAllString(body, -1)
	for _, links := range list {
		if strings.Contains(links, target_domain) {
			if strings.Contains(search_engine, passive_name) {
				links = strings.Replace(links, "<td>", "", -1)
				links = strings.Replace(links, " [TR]</td>", "", -1)
				domains += links + ","
				is_it = true
			} else {
				conv_link, _ := url.Parse(links)
				domains += conv_link.Host + ","
				is_it = true
			}
		}
	}
	return is_it, domains
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
func (Google Engine) GoogleEnum(query string) (string){

	base_url := strings.Replace(google_url, "{query}", "site:*.*."+query, -1)
	new_url := ""
	var domain_list string
	var is_it bool
	Google.loger("RUNNING : ", "GoogleEnum function is runnig")

	for i := 0; i < 100; i++ {
		new_url = strings.Replace(base_url, "{page_no}", strconv.Itoa(i*10), -1)
		defer func() {
			ggl_flag = false
		}()
		is_it, domain_list = Google.parser(Google.Meow(new_url), ggl_name)
		if is_it{
		} else if strings.Contains(Google.Meow(new_url), google_re) {
			Google.loger("ERROR : ", "GoogleEnum funciton is caught google recaptcha")
			break
		} else {
			break
		}
	}
	Google.loger("SUCCESFLY : ", "GoogleEnum function have finished")
	return domain_list
}

func (Passive Engine) PassiveDNS(domain string) (string){
	base_url := passive_url + domain
	Passive.loger("RUNNING : ", "PassiveDNS running")

	defer func() {
		pass_flag = false
	}()
	is_it, domain_list := Passive.parser(Passive.Meow(base_url), passive_name)
	if is_it{
		Passive.loger("SUCCESFLY : ", "PassiveDNS have finished")
	} else {
		Passive.loger("ERROR : ", "cannot runnig PassiveDNS")
	}
	return domain_list
}

func (e Engine) result(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, _ := template.ParseFiles("static/result.html")
	target_domain = r.FormValue("name")
	domains := e.start(string(target_domain))
	content := e.formatter(domains)

	links := Links {
		List: content,
	}
	err := parsedTemplate.Execute(w, links)
	e.control(err)
}

func (e Engine) start(domain string)(string){
	goog_list := e.GoogleEnum(domain)
	passive_list := e.PassiveDNS(domain)
//	go engine.YahooEnum(domain)

	for ggl_flag || pass_flag {
		time.Sleep(0)
	}
	return goog_list + passive_list
}

func main() {
	engine := Engine{}
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/result", engine.result)

	log.Println("http://127.0.0.1:4343/main.html")
	err := http.ListenAndServe(":4343", nil)
	engine.control(err)

}
