package middlewars

import "fmt"

var bad_request = fmt.Errorf("Bad request")
var invalid_auth_info = fmt.Errorf("%w: Invalid auth info", bad_request)
var empty_credintials = fmt.Errorf("%w: Empty credintials", bad_request)
