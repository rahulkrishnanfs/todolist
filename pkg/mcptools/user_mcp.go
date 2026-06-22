package mcptools

import (
	"net/http"
	"todolist/pkg/controller"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func SetTools(mcpserver *mcp.Server, mux *http.ServeMux, controller controller.UserController) {
	mcp.AddTool(mcpserver, &mcp.Tool{
		Name:        "create_user",
		Description: "Create a new user account for the todolist service"},
		controller.CreateUserTool)
	// mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return mcpserver },
	// 	nil))

	// this will be only use in non-prod
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(
		func(*http.Request) *mcp.Server { return mcpserver },
		&mcp.StreamableHTTPOptions{DisableLocalhostProtection: true},
	))
}
