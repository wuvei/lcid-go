package main

import (
	"fmt"
	"errors"
	"log"
	"net/http"
	"regexp"
	"bytes"
    "encoding/json"
    "io"
    "os"
)

func getCsrftoken () (string, error) {
	req, err := http.NewRequest("GET", "https://leetcode.com/problemset/all/", nil)
	if err != nil {
		log.Fatalln(err)
		return "", errors.New("Failed to generate the http request")
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
        return http.ErrUseLastResponse
    },}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return "", errors.New("Http request failed")
	}
	defer resp.Body.Close()

	reg := regexp.MustCompile(`csrftoken=(\S+); `)
	csrftoken := reg.FindStringSubmatch(resp.Header["Set-Cookie"][1])
	if csrftoken[1] == "" {
		fmt.Errorf("Fail to parse csrftoken from headers! headers: \n %v", resp.Header)
		return "", errors.New("Fail to parse csrftoken from headers!")
	} else {
		fmt.Printf("csrftoken = %s\n", csrftoken[1])
		return csrftoken[1], nil
	}
}

func FetchProblems(csrftoken string, limit int) map[string]map[string]map[string]interface{} {
	data := make(map[string]interface{})
    data["query"] = "query problemsetQuestionList($categorySlug:String,$limit:Int,$skip:Int,$filters:QuestionListFilterInput){problemsetQuestionList:questionList(categorySlug:$categorySlug limit:$limit skip:$skip filters:$filters){total:totalNum questions:data{acRate difficulty freqBar frontendQuestionId:questionFrontendId isFavor paidOnly:isPaidOnly status title titleSlug topicTags{name id slug}hasSolution hasVideoSolution}}}"
    sub_data := map[string]interface{}{"categorySlug": "", "skip": 0, "limit": limit, "filters":map[string]interface{}{} }
    data["variables"] = sub_data
    postBody, _ := json.Marshal(data)
    req, err := http.NewRequest("POST", "https://leetcode.com/graphql/", bytes.NewBuffer(postBody))
    if err != nil {
        log.Fatalln(err)
        os.Exit(1)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("cookie", fmt.Sprintf("csrftoken=%s", csrftoken))

    client := &http.Client{}
    resp, err := client.Do(req)
     if err != nil {
         log.Fatalln(err)
         os.Exit(1)
     }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    m := make(map[string]map[string]map[string]interface{})
    err = json.Unmarshal(body, &m)

    return m
}

func main() {
	csrftoken, err := getCsrftoken()
	if err != nil {
		os.Exit(1)
	}

	response_content := FetchProblems(csrftoken, 1)
	temp := response_content["data"]["problemsetQuestionList"]["total"].(float64)
	var total_count int = int(temp)
	fmt.Printf("Found %d problems in total.\n", total_count)
    fmt.Printf("Now try fetch all %d LeetCode problems...\n", total_count)

    response_content = FetchProblems(csrftoken, total_count)
	x := response_content["data"]["problemsetQuestionList"]["questions"].([]interface{})
    questions_all := make(map[string]interface{})
    for i := range x {
        temp := x[i].(map[string]interface{})
        var num string = temp["frontendQuestionId"].(string)
        questions_all[num] = x[i]
    }
    qustions_bytes, _:= json.Marshal(questions_all)

    f, err := os.OpenFile("./problems_all.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
        fmt.Println(err)
    }
    f.Write(qustions_bytes)
    f.Close()
    fmt.Printf("All %d problems info saved into problems_all.json file.\n", total_count)
}