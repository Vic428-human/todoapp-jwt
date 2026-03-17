package utils

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
)

func ObjNameFromURL(imageURL string, originalFilename string) (string, error) {
	if strings.TrimSpace(imageURL) == "" {
		objID, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}

		ext := strings.ToLower(path.Ext(originalFilename))
		if ext == "" {
			ext = ".jpg"
		}

		return fmt.Sprintf("profile/%s%s", objID.String(), ext), nil
	}

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		log.Printf("failed to parse image url %q: %v\n", imageURL, err)
		return "", err
	}

	trimmed := strings.TrimPrefix(parsedURL.Path, "/")
	parts := strings.SplitN(trimmed, "/", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid GCS image URL path: %s", parsedURL.Path)
	}

	return parts[1], nil
}
