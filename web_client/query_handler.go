package web_client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/machinebox/graphql"
)

type Test struct {
	FirstUsers string
	ServerName string
	CaseName   string
	Time       string
}

var tests []Test
var querying = false

func query(caseId string, serverType string, firstUsers string) {
	querying = true
	fmt.Println("Querying...", caseId, serverType, firstUsers)

	req := graphql.NewRequest(Cases[caseId].Query)

	var serverClient = JSClient
	if serverType == "PYTHON" {
		serverClient = PythonClient
	} else if serverType == "GO" {
		req.Header.Set("X-FirstUsers", firstUsers)
		serverClient = GoClient
	}

	firstVariable, errConv := strconv.Atoi(firstUsers)
	if errConv != nil {
		querying = false
		log.Fatalln("Error setting variables query", errConv)
		return
	}
	req.Var("first", firstVariable)

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	var respData map[string]interface{}

	var startTime = time.Now()

	if err := serverClient.Run(ctx, req, &respData); err != nil {
		querying = false
		log.Fatalln("Error running query", err)
		return
	}

	var totalTime = time.Now().Sub(startTime).Seconds()
	tests = append([]Test{
		{
			FirstUsers: firstUsers,
			ServerName: Servers[serverType].Name,
			CaseName:   Cases[caseId].Name,
			Time:       fmt.Sprintf("%.3f", totalTime),
		},
	}, tests...)
	fmt.Println("QUERY COMPLETED", totalTime)
	querying = false
}

var firstUsers = "1"
var queryServerType = GoServerType
var caseId = "case1"

func GraphQLQueryHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var invalidForm = false

		if r.Method == "POST" {
			fmt.Println("PROCESING FORM...")
			err := r.ParseForm()

			if err != nil {
				log.Fatal("Error parsing form", err)
			}

			firstUsers = strings.Trim(r.Form["firstUsers"][0], " ")
			if firstUsers == "" {
				invalidForm = true
			}

			queryServerType = r.Form["serverType"][0]
			caseId = r.Form["caseId"][0]

			if !invalidForm {
				query(caseId, queryServerType, firstUsers)

				http.Redirect(w, r, "/go-client", http.StatusSeeOther)
				return
			}
		} else {
			fmt.Println("NOT PROCESING FORM...")
		}

		t, error := template.ParseFiles("go-client.html")
		if error != nil {
			log.Fatal("Error reading HTML file", error)
		}

		data := make(map[string]interface{})

		data["firstUsers"] = firstUsers
		data["serverType"] = queryServerType
		data["caseId"] = caseId
		data["tests"] = tests
		data["Querying"] = querying
		data["InvalidForm"] = invalidForm

		t.Execute(w, data)
	})
}
