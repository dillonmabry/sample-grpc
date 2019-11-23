package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	proto "github.com/dillonmabry/sample-grpc/proto"
	"google.golang.org/grpc"

	"github.com/go-playground/lars"
	mw "github.com/go-playground/lars/_examples/middleware/logging-recovery"
)

// ApplicationGlobals ...
type ApplicationGlobals struct {
	Client proto.AddServiceClient
	Log    *log.Logger
}

func newServiceClient() proto.AddServiceClient {
	conn, err := grpc.Dial("localhost:4000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := proto.NewAddServiceClient(conn)
	return client
}

func newGlobals() *ApplicationGlobals {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	client := newServiceClient()
	return &ApplicationGlobals{
		Log:    logger,
		Client: client,
	}
}

// MyContext ...
type MyContext struct {
	*lars.Ctx
	AppContext *ApplicationGlobals
}

func newContext(l *lars.LARS) lars.Context {
	return &MyContext{
		Ctx:        lars.NewContext(l),
		AppContext: newGlobals(),
	}
}

func castCustomContext(c lars.Context, handler lars.Handler) {
	h := handler.(func(*MyContext))
	ctx := c.(*MyContext)
	h(ctx)
}

func main() {

	// Register
	l := lars.New()
	l.RegisterContext(newContext) // Cached
	l.RegisterCustomHandler(func(*MyContext) {}, castCustomContext)

	// Middleware
	l.Use(mw.LoggingAndRecovery)

	// Groups/Routes
	operations := l.Group("/operations")
	operations.Get("/add/:a/:b", Add)
	operations.Get("/multiply/:a/:b", Multiply)

	// Serve
	if err := http.ListenAndServe(":8080", l.Serve()); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// Add ...
func Add(c *MyContext) {
	a, err := strconv.ParseUint(c.Param("a"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Bad input for a")
		return
	}

	b, err := strconv.ParseUint(c.Param("b"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Bad input for b")
		return
	}

	req := &proto.Request{A: int64(a), B: int64(b)}

	if res, err := c.AppContext.Client.Add(c, req); err == nil {
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusInternalServerError, "Error adding inputs")
		return
	}
}

// Multiply ...
func Multiply(c *MyContext) {
	a, err := strconv.ParseUint(c.Param("a"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Bad input for a")
		return
	}

	b, err := strconv.ParseUint(c.Param("b"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Bad input for b")
		return
	}

	req := &proto.Request{A: int64(a), B: int64(b)}

	if res, err := c.AppContext.Client.Multiply(c, req); err == nil {
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusInternalServerError, "Error adding inputs")
		return
	}
}
