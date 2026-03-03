package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"authentication test server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Add tool to get user email
	s.AddTool(
		mcp.NewTool("get_authentication_info",
			mcp.WithDescription("Returns the authentication information of the currently connected user"),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Retrieve the Authorization header
			authHeader := req.Header.Get("Authorization")

			if authHeader == "" {
				var sb strings.Builder
				sb.WriteString("anonymous@example.com\n\nHeaders:\n")
				for name, values := range req.Header {
					for _, value := range values {
						sb.WriteString(fmt.Sprintf("%s: %s\n", name, value))
					}
				}
				return mcp.NewToolResultText(sb.String()), nil
			}

			// Extract Bearer token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				// No "Bearer " prefix found
				return mcp.NewToolResultError("Invalid Authorization header format. Expected: Bearer <token>"), nil
			}

			// Parse JWT token without verification (to read claims)
			// In production, you should verify the token with the proper secret/public key
			token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to parse JWT token: %v", err)), nil
			}

			// Extract and return all claims as key-value pairs
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				var sb strings.Builder
				for key, value := range claims {
					sb.WriteString(fmt.Sprintf("%s: %v\n", key, value))
				}
				return mcp.NewToolResultText(sb.String()), nil
			}
			return mcp.NewToolResultError("Failed to extract claims from JWT token"), nil
		},
	)

	// Create streamable HTTP server
	httpServer := server.NewStreamableHTTPServer(s)

	addr := ":3000"
	log.Printf("MCP Server starting on http://localhost%s", addr)

	// Start the HTTP server
	if err := httpServer.Start(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
