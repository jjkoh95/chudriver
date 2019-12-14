package chulink

import "google.golang.org/api/firebasedynamiclinks/v1"

// Chulink - Firebase Dynamic Links wrapper.
type Chulink struct {
	Link *firebasedynamiclinks.Service
}

// GenerateLink is the function to shorten longURL.
func (chulink *Chulink) GenerateLink(longURL, firebaseRegisteredDomain string) (string, error) {
	shortLinkRequest := &firebasedynamiclinks.CreateShortDynamicLinkRequest{
		DynamicLinkInfo: &firebasedynamiclinks.DynamicLinkInfo{
			DomainUriPrefix: firebaseRegisteredDomain,
			Link:            longURL,
		},
	}

	shortLinkCall := chulink.Link.ShortLinks.Create(shortLinkRequest)

	res, err := shortLinkCall.Do()
	if err != nil {
		return "", err
	}

	return res.ShortLink, nil
}
