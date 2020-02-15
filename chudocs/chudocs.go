package chudocs

import (
	"strings"

	"google.golang.org/api/docs/v1"
)

// Chudocs - Google docs wrapper.
type Chudocs struct {
	Docs *docs.Service
}

// ReadDocByID - Read doc as raw string.
func (chudocs *Chudocs) ReadDocByID(documentID string) (string, error) {
	doc, err := chudocs.Docs.Documents.Get(documentID).Do()
	if err != nil {
		return "", err
	}
	return readStructuralElements(doc.Body.Content), nil
}

func readStructuralElements(struturalElements []*docs.StructuralElement) string {
	var sb strings.Builder
	for _, structuralElement := range struturalElements {
		if structuralElement.Paragraph == nil {
			continue
		}
		for _, paragraphElement := range structuralElement.Paragraph.Elements {
			textRun := paragraphElement.TextRun
			sb.WriteString(textRun.Content)
		}
	}
	return sb.String()
}
