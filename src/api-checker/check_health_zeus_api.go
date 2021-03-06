package main

import (
	"auth"
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

var bodyStatus int = 0

func HttpQueryGet(url string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	for i := range cookie {
		req.AddCookie(cookie[i])
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	bodyStatus = resp.StatusCode

	return string(body)
}

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

	mux := http.NewServeMux()
	/// FIRST REQUEST ///
	mux.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {

		cookie = auth.HttpQueryPost(urlLogin, jsonStr)
		var c_name string

		for _, c := range cookie {
			if c.Name == "usid" {
				c_name = c.Name
			}
		}

		if c_name == "usid" {
			io.WriteString(writer, "1")
		} else {
			io.WriteString(writer, "0")
		}
	})

	mux.HandleFunc("/podft", func(writer http.ResponseWriter, request *http.Request) {

		body_s := HttpQueryGet(urlPodft)
		bodyStatus_podft := bodyStatus

		var m map[string]interface{}
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
		getTime := HttpQueryGet(urlPodft)
		if getTime != "" {
			fmt.Println(0)
		}
		elapsed := time.Since(start).Seconds()
		elapsed_int := elapsed * 1000
		elapsed_str := strconv.FormatFloat(elapsed_int, 'f', -1, 64)
		io.WriteString(writer, elapsed_str)
	})

	mux.HandleFunc("/query-tasks", func(writer http.ResponseWriter, request *http.Request) {

		body_s := HttpQueryGet(urlStat)
		bodyStatus_stat := bodyStatus

		var m map[string]interface{}
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

		body_s := HttpQueryGet(urlStat)
		bodyStatus_stat := bodyStatus

		var m map[string]interface{}
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

		body_s := HttpQueryGet(urlStat)
		bodyStatus_stat := bodyStatus

		var m map[string]interface{}
		json.Unmarshal([]byte(body_s), &m)

		if bodyStatus_stat == 200 {
			f := m["meta"]
			j := f.(map[string]interface{})
			u := j["total_by_type"]
			us := u.(map[string]interface{})
			qu := us["search_lists"]
			if _, ok := us["search_lists"]; ok {
				users := qu.(map[string]interface{})
				// TRANSFORM DATA
				//var str1 string
				str1 := ""
				str1 = fmt.Sprintf("%v", users)
				fmt.Println("str", str1)

				reg, _ := regexp.Compile("[^0-9]+")
				get_num := reg.ReplaceAllString(str1, "\n")

				get_array := strings.Split(get_num, "\n")

				rest := 0
				for _, v := range get_array {
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

		body_s := HttpQueryGet(urlStat)
		bodyStatus_stat := bodyStatus

		var m map[string]interface{}
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

		body_s := HttpQueryGet(urlStat)
		bodyStatus_query_t := bodyStatus

		var m map[string]interface{}
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

		body_s := HttpQueryGet(urlStat)
		bodyStatus_stat := bodyStatus

		var m map[string]interface{}
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
