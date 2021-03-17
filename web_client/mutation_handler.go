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

type MutationTest struct {
	ClientsTotal string
	ServerName   string
	Time         string
}

var mutationTests []MutationTest
var runningMutation = false

func runMutation(serverType string, clientsTotal string) {
	if runningMutation == true {
		return
	}

	runningMutation = true
	fmt.Println("Running Mutation...", serverType, clientsTotal)

	req := graphql.NewRequest(`
		mutation CreateClient($client: ClientInput!){
			createClient(client: $client) {
				id
				name
			}
		}
	`)

	var serverClient = JSClient
	if serverType == "PYTHON" {
		serverClient = PythonClient
	} else if serverType == "GO" {
		serverClient = GoClient
	}

	clientsTotalAsInt, errConv := strconv.Atoi(clientsTotal)
	if errConv != nil {
		runningMutation = false
		log.Fatalln("Error converting clientsTotal", errConv)
		return
	}

	var startTime = time.Now()
	var iso8601Layout = "2006-01-02T15:04:05Z"
	var createdAt = startTime.Format(iso8601Layout)

	for i := 1; i <= clientsTotalAsInt; i++ {
		clientInput := make(map[string]interface{})

		clientInput["name"] = fmt.Sprintf("Juan GO %d", i)
		clientInput["email"] = fmt.Sprintf("jp%d@mail.com", i)
		clientInput["last_name"] = fmt.Sprintf("P\u00E9rez %d", i)
		clientInput["birthday"] = time.Date(2020, time.Month((i%12)+1), (i%28)+1, 10, 0, 0, 0, time.UTC).Format(iso8601Layout)
		clientInput["address"] = fmt.Sprintf("Avenida de las Merceditas y Calle XYZ-%d", i)
		clientInput["created_at"] = createdAt

		req.Var("client", clientInput)

		// set header fields
		req.Header.Set("Cache-Control", "no-cache")

		// define a Context for the request
		ctx := context.Background()

		// run it and capture the response
		var respData map[string]interface{}

		if err := serverClient.Run(ctx, req, &respData); err != nil {
			runningMutation = false
			log.Fatalln("Error running mutation", err)
			return
		}
	}

	var totalTime = time.Now().Sub(startTime).Seconds()
	mutationTests = append([]MutationTest{
		{
			ClientsTotal: clientsTotal,
			ServerName:   Servers[serverType].Name,
			Time:         fmt.Sprintf("%.3f", totalTime),
		},
	}, mutationTests...)
	fmt.Println(clientsTotal+" MUTATIONS COMPLETED", totalTime)
	runningMutation = false
}

var clientsTotal = "1"
var mutationServerType = GoServerType

func GraphQLMutationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var invalidForm = false

		if r.Method == "POST" {
			fmt.Println("PROCESING MUTATION FORM...")
			err := r.ParseForm()

			if err != nil {
				log.Fatal("Error parsing form", err)
			}

			println("FORM", r.Body)

			clientsTotal = strings.Trim(r.Form["clientsTotal"][0], " ")
			if firstUsers == "" {
				invalidForm = true
			}

			mutationServerType = r.Form["serverType"][0]

			if !invalidForm {
				runMutation(mutationServerType, clientsTotal)

				http.Redirect(w, r, "/go-client/mutation", http.StatusSeeOther)
				return
			}
		} else {
			fmt.Println("NOT PROCESING FORM...")
		}

		t, error := template.ParseFiles("go-client-mutation.html")
		if error != nil {
			log.Fatal("Error reading HTML file", error)
		}

		data := make(map[string]interface{})

		data["ClientsTotal"] = clientsTotal
		data["ServerType"] = mutationServerType
		data["Tests"] = mutationTests
		data["RunningMutation"] = runningMutation
		data["InvalidForm"] = invalidForm

		t.Execute(w, data)
	})
}
