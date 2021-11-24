package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"getenv"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Query struct {
	User User `json:"user"`
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Usid     string `json:"usid"`
}

// Get cookie

var cookie []*http.Cookie

// credentials

var custUser string = getenv.GoDotEnvVariable("customerUser")
var custPass string = getenv.GoDotEnvVariable("customerPassword")
var dbUser string = getenv.GoDotEnvVariable("dbUser")
var dbPassword string = getenv.GoDotEnvVariable("dbPassword")
var dbHost string = getenv.GoDotEnvVariable("dbHost")
var dbPort string = getenv.GoDotEnvVariable("dbPort")
var dbName string = getenv.GoDotEnvVariable("dbName")

// urls
var urlLogin string = "http://customer.wd.xco.devel.ifx/zeus/users/sign_in"
var urlPodft string = "http://customer.wd.xco.devel.ifx/zeus/profiles/podft/210.json"
var urlStat string = "http://customer.wd.xco.devel.ifx/zeus/statistics"

func main() {

	lgn_struct := Query{
		User: User{
			Login:    custUser,
			Password: custPass,
			Usid:     "",
		},
	}

	jsonData, _ := json.Marshal(lgn_struct)
	jsonStr := []byte(jsonData)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}

	/// FIRST REQUEST ///

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {

		req_lgn, err := http.NewRequest("POST", urlLogin, bytes.NewBuffer(jsonStr))
		req_lgn.Header.Set("Content-Type", "application/json")
		resp_lgn, err := client.Do(req_lgn)
		if err != nil {
			panic(err)
		}
		cookie = resp_lgn.Cookies()

		for _, c := range cookie {
			fmt.Println(c.Name, c.Value)
		}

		defer resp_lgn.Body.Close()

		body_lgn, _ := ioutil.ReadAll(resp_lgn.Body)
		bodyStatus_lgn := resp_lgn.StatusCode
		fmt.Printf("%+v\n", string(body_lgn), bodyStatus_lgn)

		if bodyStatus_lgn == 201 {
			io.WriteString(writer, "1")
		} else {
			io.WriteString(writer, "0")
		}

	})

	mux.HandleFunc("/podft", func(writer http.ResponseWriter, request *http.Request) {

		req_podft, err := http.NewRequest("GET", urlPodft, nil)
		if err != nil {
			panic(err)
		}
		for i := range cookie {
			req_podft.AddCookie(cookie[i])
		}
		resp_podft, err := client.Do(req_podft)
		if err != nil {
			panic(err)
		}
		defer resp_podft.Body.Close()

		body_podft, err := ioutil.ReadAll(resp_podft.Body)
		if err != nil {
			panic(err)
		}
		bodyStatus_podft := resp_podft.StatusCode

		var m map[string]interface{}
		body_s := string(body_podft)
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_podft == 200 {
			f := m["data"]
			j := f.(map[string]interface{})
			if j["id"] != "" && j["type"] != "" {
				io.WriteString(writer, "1")
			}
		} else {
			io.WriteString(writer, "0")
		}
	})

	mux.HandleFunc("/time", func(writer http.ResponseWriter, request *http.Request) {

		start := time.Now()
		req_podft, err := http.NewRequest("GET", urlPodft, nil)
		if err != nil {
			panic(err)
		}

		for i := range cookie {
			req_podft.AddCookie(cookie[i])
		}
		resp_podft, err := client.Do(req_podft)
		if err != nil {
			panic(err)
		}

		defer resp_podft.Body.Close()

		elapsed := time.Since(start).Seconds()
		elapsed_int := elapsed * 1000
		elapsed_str := strconv.FormatFloat(elapsed_int, 'f', -1, 64)
		io.WriteString(writer, elapsed_str)
	})

	mux.HandleFunc("/query-tasks", func(writer http.ResponseWriter, request *http.Request) {

		req_stat, err := http.NewRequest("GET", urlStat, nil)
		if err != nil {
			panic(err)
		}
		for i := range cookie {
			req_stat.AddCookie(cookie[i])
		}
		resp_stat, err := client.Do(req_stat)
		if err != nil {
			panic(err)
		}
		defer resp_stat.Body.Close()

		body_stat, err := ioutil.ReadAll(resp_stat.Body)
		if err != nil {
			panic(err)
		}
		bodyStatus_stat := resp_stat.StatusCode

		var m map[string]interface{}
		body_s := string(body_stat)
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_stat == 200 {
			f := m["meta"]
			j := f.(map[string]interface{})
			total_int := j["total"].(float64)
			if total_int <= 200 {
				io.WriteString(writer, "1")
			}
		} else {
			io.WriteString(writer, "0")
		}
	})

	mux.HandleFunc("/total-tasks", func(writer http.ResponseWriter, request *http.Request) {

		req_stat, err := http.NewRequest("GET", urlStat, nil)
		if err != nil {
			panic(err)
		}
		for i := range cookie {
			req_stat.AddCookie(cookie[i])
		}
		resp_stat, err := client.Do(req_stat)
		if err != nil {
			panic(err)
		}
		defer resp_stat.Body.Close()

		body_stat, err := ioutil.ReadAll(resp_stat.Body)
		if err != nil {
			panic(err)
		}
		bodyStatus_stat := resp_stat.StatusCode

		var m map[string]interface{}
		body_s := string(body_stat)
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_stat == 200 {
			f := m["meta"]
			j := f.(map[string]interface{})
			total_int := j["total"].(float64)
			total_str := strconv.FormatFloat(total_int, 'f', -1, 64)
			io.WriteString(writer, total_str)
		} else {
			io.WriteString(writer, "-1")
		}
	})

	mux.HandleFunc("/batch", func(writer http.ResponseWriter, request *http.Request) {

		req_stat, err := http.NewRequest("GET", urlStat, nil)
		if err != nil {
			panic(err)
		}
		for i := range cookie {
			req_stat.AddCookie(cookie[i])
		}
		resp_stat, err := client.Do(req_stat)
		if err != nil {
			panic(err)
		}
		defer resp_stat.Body.Close()

		body_stat, err := ioutil.ReadAll(resp_stat.Body)
		if err != nil {
			panic(err)
		}
		bodyStatus_stat := resp_stat.StatusCode

		var m map[string]interface{}
		body_s := string(body_stat)
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_stat == 200 {
			f := m["meta"]
			j := f.(map[string]interface{})
			u := j["total_by_type"]
			us := u.(map[string]interface{})
			qu := us["search_lists"]
			if _, ok := us["search_lists"]; ok {
				users := qu.(map[string]interface{})

				var str1 string
				str1 = fmt.Sprintf("%v", users)
				fmt.Println("str", str1)

				reg, _ := regexp.Compile("[^0-9]+")
				prcs := reg.ReplaceAllString(str1, "\n")

				pr := strings.Split(prcs, "\n")

				rest := 0
				for _, v := range pr {
					if v != "" {
						v_int, _ := strconv.Atoi(v)
						rest = rest + v_int
						//fmt.Println(x,v)
					}
				}
				//fmt.Println(rest)

				total_int := float64(rest)
				total_str := strconv.FormatFloat(total_int, 'f', -1, 64)
				io.WriteString(writer, total_str)
			} else {
				io.WriteString(writer, "0")
			}
		} else {
			io.WriteString(writer, "-1")
		}
	})

	mux.HandleFunc("/query-users", func(writer http.ResponseWriter, request *http.Request) {

		req_stat, err := http.NewRequest("GET", urlStat, nil)
		if err != nil {
			panic(err)
		}
		for i := range cookie {
			req_stat.AddCookie(cookie[i])
		}
		resp_stat, err := client.Do(req_stat)
		if err != nil {
			panic(err)
		}
		defer resp_stat.Body.Close()

		body_stat, err := ioutil.ReadAll(resp_stat.Body)
		if err != nil {
			panic(err)
		}
		bodyStatus_stat := resp_stat.StatusCode

		var m map[string]interface{}
		body_s := string(body_stat)
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_stat == 200 {
			f := m["meta"]
			j := f.(map[string]interface{})
			u := j["total_by_user"]
			users := u.(map[string]interface{})
			l := len(users)
			if l <= 30 {
				io.WriteString(writer, "1")
			}
		} else {
			io.WriteString(writer, "0")
		}
	})

	mux.HandleFunc("/total-users", func(writer http.ResponseWriter, request *http.Request) {

		req_query_t, err := http.NewRequest("GET", urlStat, nil)
		if err != nil {
			panic(err)
		}
		for i := range cookie {
			req_query_t.AddCookie(cookie[i])
		}
		resp_query_t, err := client.Do(req_query_t)
		if err != nil {
			panic(err)
		}
		defer resp_query_t.Body.Close()

		body_query_t, err := ioutil.ReadAll(resp_query_t.Body)
		if err != nil {
			panic(err)
		}
		bodyStatus_query_t := resp_query_t.StatusCode
		//    fmt.Printf("%+v\n", string(body_query_t),bodyStatus_query_t)

		var m map[string]interface{}
		body_s := string(body_query_t)
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_query_t == 200 {
			f := m["meta"]
			j := f.(map[string]interface{})
			u := j["total_by_user"]
			users := u.(map[string]interface{})
			res := 0
			for _, v := range users {
				i := reflect.ValueOf(v).Len()
				res = i + res
			}
			fmt.Println(res)
			l_int := float64(res)
			l_str := strconv.FormatFloat(l_int, 'f', -1, 64)
			io.WriteString(writer, l_str)
		} else {
			io.WriteString(writer, "-1")
		}

		time.Sleep(time.Second * 1)

	})

	mux.HandleFunc("/collect", func(writer http.ResponseWriter, request *http.Request) {

		req_stat, err := http.NewRequest("GET", urlStat, nil)
		if err != nil {
			panic(err)
		}
		for i := range cookie {
			req_stat.AddCookie(cookie[i])
		}
		resp_stat, err := client.Do(req_stat)
		if err != nil {
			panic(err)
		}
		defer resp_stat.Body.Close()

		body_stat, err := ioutil.ReadAll(resp_stat.Body)
		if err != nil {
			panic(err)
		}
		bodyStatus_stat := resp_stat.StatusCode

		var m map[string]interface{}
		body_s := string(body_stat)
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_stat == 200 {
			io.WriteString(writer, body_s)
			//}
		} else {
			io.WriteString(writer, "-1")
		}
	})

	mux.HandleFunc("/uniq-session", func(writer http.ResponseWriter, request *http.Request) {

		connStr := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName

		/*
			cfg := mysql.Config{
				User:                dbUser,
				Passwd:              dbPassword,
				Net:                 "tcp",
				Addr:                "127.0.0.1:3306",
				DBName:              "world",
				AllowNativePassword: true,
			}
		*/

		db, err := sql.Open("mysql", connStr)
		if err != nil {
			panic(err.Error())
		}

		defer db.Close()

		var result string

		ses, err := db.Query("select count(unique_session_id) from users where unique_session_id IS NOT NULL")
		if err != nil {
			panic(err.Error())
		}
		defer ses.Close()

		for ses.Next() {
			if err := ses.Scan(&result); err != nil {
				panic(err.Error())
			}
		}
		io.WriteString(writer, result)
	})

	http.ListenAndServe(":8080", mux)
}
