package chudocs_test

import (
	"testing"

	"github.com/jjkoh95/chudriver/chuauth"
	"github.com/jjkoh95/chudriver/chudocs"
	"google.golang.org/api/docs/v1"
)

const (
	publicDocID  = "1fb-qHX9zXbWf5Ah1yeIDqCLHAjRtwpg5QUo8mBPm5Q0"
	privateDocID = "1diGwqJRVrUsMKNdZHnNdcguy7HQQSl7h9D3160NSAIc"
)

func TestReadingDoc(t *testing.T) {
	var chudocsWrapper chudocs.Chudocs
	var err error
	chudocsWrapper.Docs, err = docs.New(chuauth.GetClientFromJSON("credentials.json", docs.DocumentsReadonlyScope))
	if err != nil {
		t.Error("Expected to create client service without error")
	}

	// publicly shared docs
	_, err = chudocsWrapper.ReadDocByID(publicDocID)
	if err != nil {
		t.Error("Expected to access to public docs without error")
	}

	// private docs (with access gained)
	_, err = chudocsWrapper.ReadDocByID(privateDocID)
	if err != nil {
		t.Error("Expected to access to private docs without error")
	}
}
