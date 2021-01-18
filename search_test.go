package googlesearch_test

import (
	"context"
	"testing"

	"github.com/psclil/google-search"
)

var ctx = context.Background()

func TestSearch(t *testing.T) {

	q := "Hello World"

	opts := googlesearch.SearchOptions{
		Limit: 20,
	}

	returnLinks, err := googlesearch.Search(ctx, q, opts)
	if err != nil {
		t.Errorf("something went wrong: %v", err)
		return
	}

	if len(returnLinks) == 0 {
		t.Errorf("no results returned: %v", returnLinks)
	}
}

func TestChrome87(t *testing.T) {
	q := `site:"dropbox.com" ("Api Reference" OR "Api Documentation" OR "API Documentation")`
	userAgent := `"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"`

	returnLinks, err := googlesearch.Search(ctx, q, googlesearch.SearchOptions{
		UserAgent: userAgent,
	})
	if err != nil {
		t.Errorf("something went wrong: %v", err)
		return
	}

	if len(returnLinks) == 0 {
		t.Errorf("no results returned: %v", returnLinks)
	}
}
