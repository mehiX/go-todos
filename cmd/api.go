package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mehix/go-todos/internal/db"
	"github.com/mehix/go-todos/pkg/todos"
)

var addr = flag.String("http", "127.0.0.1:8080", "Address to serve HTTP")
var dsn = flag.String("dsn", "test:test@tcp(127.0.0.1)/test", "Database connection string (MariaDB)")

func init() {
	flag.Parse()
}

func Execute() {

	var dbRepo todos.Repository
	var conn *sql.DB
	var err error

	if *dsn == "" {
		dbRepo = todos.NewInMemoryRepository()
	} else {
		conn, err = db.ConnWithRetry(db.Conn, 5, time.Second, time.Minute)(context.Background(), *dsn)
		if err != nil {
			log.Println("DB connection failed", err)
			log.Println("Using in-memory storage")
			dbRepo = todos.NewInMemoryRepository()
		} else {
			fmt.Println("Connected to database")
			dbRepo = todos.NewDbRepository(conn)
		}
	}

	svc := todos.NewService(todos.WithRepo(dbRepo))

	srvr := http.Server{
		Addr:              *addr,
		Handler:           todos.Handler(svc),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	done := make(chan bool)

	go func() {
		defer close(done)

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

		<-ch

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		fmt.Println("\033[K\rInitiate shutdown ...")
		if err := srvr.Shutdown(ctx); err != nil {
			log.Printf("Shutting down: %v\n", err)
		}

		if conn != nil {
			fmt.Println("Closing DB...")
			conn.Close()
		}
	}()

	fmt.Printf("Listening on %s\n", *addr)
	printRoutesHelp(srvr.Handler)

	if err := srvr.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("Cleaning up ...")
			<-done
		} else {
			log.Println(err)
		}
	}

	fmt.Println("Done")
}

func printRoutesHelp(h http.Handler) {
	r, ok := h.(*chi.Mux)
	if !ok {
		return
	}
	var out = tabwriter.NewWriter(os.Stdout, 10, 8, 0, '\t', 0)
	printHelp(out, "", r.Routes())
	out.Flush()

}

func printHelp(out *tabwriter.Writer, parentPattern string, routes []chi.Route) {

	fmt.Fprintln(out)

	for _, r := range routes {
		ptrn := strings.TrimSuffix(r.Pattern, "/*")
		if r.SubRoutes != nil {
			printHelp(out, parentPattern+ptrn, r.SubRoutes.Routes())
		} else {
			for m := range r.Handlers {
				fmt.Fprintf(out, "[%s]\t%s\n", m, parentPattern+ptrn)
			}
			fmt.Fprintln(out)
		}

	}

}
