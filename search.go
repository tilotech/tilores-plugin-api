package api

// SearchParameters represents the structured properties that can be used for
// searching.
//
// Since the actual search parameters can be different for each customers, it
// is no further defined how the parameters may look like. For this reason, the
// SearchParameters must only be used in combination with matching rules that
// know how to access the relevant properties.
type SearchParameters map[string]any
