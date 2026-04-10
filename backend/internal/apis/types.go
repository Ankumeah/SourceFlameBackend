package apis

import "net/http"

func internal_server_error() (int, map[string]any) { return http.StatusInternalServerError, map[string]any { "error": "Internal server error" } }
func bad_request(err error) (int, map[string]any) { return http.StatusBadRequest, map[string]any { "error": err.Error() } }
func invalid_user() (int, map[string]any) { return http.StatusNotFound, map[string]any { "error": "Invalid user" } }
func invalid_repo() (int, map[string]any) { return http.StatusNotFound, map[string]any { "error": "Invalid repo" } }
