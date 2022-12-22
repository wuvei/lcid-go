package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var problems_all map[string]interface{}

func go_leetcode_us(w http.ResponseWriter, r *http.Request) {
	pid := strings.Split(r.URL.Path, "/")[1]
	if problem_info, ok := problems_all[pid]; ok {
		problem_title := problem_info.(map[string]interface{})["titleSlug"].(string)
		link := fmt.Sprintf("https://leetcode.com/problems/%s/", problem_title)
		log.Printf("Get Leetcode Problem: %s\n", link)
		http.Redirect(w, r, link, 302)
	} else {
		fmt.Fprintf(w, "Fail to redirect to leetcode problem %s page.", pid)
	}
}

func go_leetcode_cn(w http.ResponseWriter, r *http.Request) {
	pid := strings.Split(r.URL.Path, "/")[2]
	if problem_info, ok := problems_all[pid]; ok {
		problem_title := problem_info.(map[string]interface{})["titleSlug"].(string)
		link := fmt.Sprintf("https://leetcode.cn/problems/%s/", problem_title)
		log.Printf("Get Leetcode Problem: %s\n", link)
		http.Redirect(w, r, link, 302)
	} else {
		fmt.Fprintf(w, "Fail to redirect to leetcode-cn problem %s page.", pid)
	}
}

func info(w http.ResponseWriter, r *http.Request) {
	pid := strings.Split(r.URL.Path, "/")[2]
	if problem_info, ok := problems_all[pid]; ok {
		prob_info, _ := json.Marshal(problem_info)
		fmt.Fprintf(w, string(prob_info))
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Fail to get info of leetcode problem %s.", pid)
	}
}

// 解析url 函数
func router(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "//" || r.URL.Path == "/index.html" || r.URL.Path == "/index.htm" {
		log.Printf("Get index.html\n")
		http.ServeFile(w, r, "./index.html")
	} else if r.URL.Path == "/favicon.ico" {
		log.Printf("Get favicon.ico\n")
		http.ServeFile(w, r, "./favicon.ico")
	} else if match, _ := regexp.MatchString(`/info/[1-9]+[0-9]*`, r.URL.Path); match {
		info(w, r)
	} else if match, _ := regexp.MatchString(`/cn/[1-9]+[0-9]*`, r.URL.Path); match {
		go_leetcode_cn(w, r)
	} else if match, _ := regexp.MatchString(`/[1-9]+[0-9]*`, r.URL.Path); match {
		go_leetcode_us(w, r)
	} else {
		w.WriteHeader(404)
	}
}

func main() {
	port := 9191
	if len(os.Args) > 2 {
		fmt.Println("At most 1 argument (port)!")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		temp, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Input port is not an Interger")
			os.Exit(1)
		}
		port = temp
	}
	content, _ := os.ReadFile("./problems_all.json")
	problems_all = make(map[string]interface{})
	_ = json.Unmarshal(content, &problems_all)

	http.ListenAndServe(fmt.Sprintf(":%d", port), http.HandlerFunc(router))
}
