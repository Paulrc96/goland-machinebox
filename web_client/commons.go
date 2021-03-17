package web_client

import "github.com/machinebox/graphql"

type Case struct {
	Name  string
	Query string
}

type Server struct {
	Name string
}

var JSServerURL = "http://localhost:4000/"
var PythonServerURL = "http://localhost:8000/graphql_server"
var GoServerURL = "http://localhost:8080/query"

// create a JSClient (safe to share across requests)
var JSClient = graphql.NewClient(JSServerURL)
var PythonClient = graphql.NewClient(PythonServerURL)
var GoClient = graphql.NewClient(GoServerURL)

var Servers = make(map[string]Server)
var Cases = make(map[string]Case)

var GoServerType = "GO"

func InitMaps() {
	Cases["case1"] = Case{
		Name: "Caso 1",
		Query: `
			query ($first: Int!) {
				users (first: $first) {
					id
					name
					last_name
					email
					address
					birthday
				}
			}
		`,
	}
	Cases["case2"] = Case{
		Name: "Caso 2",
		Query: `
			query ($first: Int!) {
				users (first: $first) {
					id
					name
					last_name
					email
					address
					birthday
					posts {
						title
						description
					}
				}
			}
		`,
	}
	Cases["case3"] = Case{
		Name: "Caso 3",
		Query: `
			query ($first: Int!) {
				users (first: $first) {
					id
					name
					last_name
					email
					address
					birthday
					posts {
						title
						description
						comments {
							post_id
							description
						}
					}
				}
			}
		`,
	}

	Servers["JS"] = Server{
		Name: "Javascript",
	}
	Servers["PYTHON"] = Server{
		Name: "Python",
	}
	Servers[GoServerType] = Server{
		Name: "Go",
	}
}
