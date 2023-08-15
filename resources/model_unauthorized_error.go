/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UnauthorizedError struct {
	// Application-specific error code, expressed as a string value  - invalid_access_token error - 1 - invalid_signature - 2 - nonce_not_found - 3 - not_likely_human - 4 - score_too_low - 5 - invalid_gateway_token - 6
	Code string `json:"code"`
	// Human-readable explanation specific to this occurrence of the problem
	Detail *string `json:"detail,omitempty"`
	// Object containing non-standard meta-information about the error
	Meta *map[string]interface{} `json:"meta,omitempty"`
	// HTTP status code applicable to this problem
	Status int32 `json:"status"`
	// Short, human-readable summary of the problem
	Title string `json:"title"`
}
