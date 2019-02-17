package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/go-toschool/helenia/assistants"
	"github.com/go-toschool/palermo/auth"
	"github.com/go-toschool/platon/talks"
	"github.com/urfave/negroni"

	"google.golang.org/grpc"

	"github.com/go-toschool/sicily/cmd/server/api"
	"github.com/go-toschool/sicily/cmd/server/healthz"
	"github.com/go-toschool/sicily/cmd/server/home"
	"github.com/go-toschool/sicily/cmd/server/prometheus"
	"github.com/go-toschool/sicily/graph"
	"github.com/go-toschool/sicily/graph/mutation"
	"github.com/go-toschool/sicily/graph/queries"
	"github.com/go-toschool/syracuse/citizens"
	"github.com/graphql-go/graphql"
)

func main() {
	citizensHost := flag.String("citizens-host", "localhost", "Citizens service host")
	citizensPort := flag.Int64("citizens-port", 8001, "Citizens service port")
	palermoHost := flag.String("palermo-host", "localhost", "Palermo service host")
	palermoPort := flag.Int64("palermo-port", 8003, "Palermo service port")
	platoHost := flag.String("plato-host", "localhost", "Plato service host")
	platoPort := flag.Int64("plato-port", 8004, "Plato service port")
	heleniaHost := flag.String("helenia-host", "localhost", "Helenia service host")
	heleniaPort := flag.Int64("helenia-port", 8005, "Helenia service port")

	flag.Parse()
	// Connect services
	citizenURL := fmt.Sprintf("%s:%d", *citizensHost, *citizensPort)
	fmt.Printf("Connecting to: %s\n", citizenURL)
	citizensConn, err := grpc.Dial(citizenURL, grpc.WithInsecure())
	check("citizens connection:", err)

	palermoURL := fmt.Sprintf("%s:%d", *palermoHost, *palermoPort)
	fmt.Printf("Connecting to: %s\n", palermoURL)
	palermoConn, err := grpc.Dial(palermoURL, grpc.WithInsecure())
	check("palermo connection:", err)

	platoURL := fmt.Sprintf("%s:%d", *platoHost, *platoPort)
	fmt.Printf("Connecting to: %s\n", platoURL)
	platoConn, err := grpc.Dial(platoURL, grpc.WithInsecure())
	check("plato connection:", err)

	heleniaURL := fmt.Sprintf("%s:%d", *heleniaHost, *heleniaPort)
	fmt.Printf("Connecting to: %s\n", heleniaURL)
	heleniaConn, err := grpc.Dial(heleniaURL, grpc.WithInsecure())
	check("helenia connection:", err)

	// Initialize citizen client
	citizenSvc := citizens.NewCitizenshipClient(citizensConn)
	palermoSvc := auth.NewAuthServiceClient(palermoConn)
	talksSvc := talks.NewTalkingClient(platoConn)
	assistantSvc := assistants.NewAssistantsClient(heleniaConn)

	graphCtx := &graph.Context{
		UserService:       citizenSvc,
		SessionService:    palermoSvc,
		TalkService:       talksSvc,
		AssistantsService: assistantSvc,
	}

	// graphql schemas
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queries.Queries(graphCtx),
		Mutation: mutation.Mutations(graphCtx),
	})
	check("session schema:", err)

	mux := http.NewServeMux()

	// public endpoint
	mux.Handle("/", home.Routes())
	mux.Handle("/metrics", prometheus.Routes())
	mux.Handle("/healthz", healthz.Routes())

	// private endpoint
	ac := &api.Context{
		User:    citizenSvc,
		Session: palermoSvc,
		Schema:  schema,
	}

	mux.Handle("/graphql", api.Routes(ac))

	log.Println("Now server is running on port 3000")
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(mux)

	check("server: ", http.ListenAndServe(":3000", n))
}

func check(section string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%s %v", section, err))
	}
}
