package prow

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var URL = "https://prow.k8s.io"

type Prow struct {
	ProwURL  string
	JunitURL string
	Error    string
}

// BuildLog holds the build log
type BuildLog struct {
	Error   string
	LensURL string
}

type ProwInterface interface {
	GetSpyGlassLens() (*BuildLog, error)
}

func NewProw(prowUrl string) ProwInterface {
	if prowUrl == "" {
		prowUrl = URL
	}
	return &Prow{ProwURL: prowUrl}
}

// GetSpyGlassLens returns a jUnit object with parsed error from the build
// spyglass pane. This requires multiple requests to scrape JS files
// rendered in the main page, and used later for next pages.
func (t *Prow) GetSpyGlassLens() (*BuildLog, error) {
	body, err := getHTTPResponse(t.ProwURL)
	if err != nil {
		return nil, err
	}

	// Extract Junit Iframe data and build the next URL.
	var lensURL string
	if lensURL, err = extractBuildLensURL(body); err != nil {
		return nil, err
	}

	// Extract Lens build body from the glass pane request.
	if body, err = getHTTPResponse(lensURL); err != nil {
		return nil, err
	}
	buildLog, err := extractBuildLogs(body)
	if err != nil {
		return nil, err
	}
	buildLog.LensURL = lensURL
	return buildLog, nil
}

// extractBuildLogs returns the buildlog serialization error logs from
// build log glass pane.
func extractBuildLogs(body io.Reader) (*BuildLog, error) {
	var chunks, errors = 0, ""

	// Parse the HTML body
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.TextNode && node.Data != "" && !strings.Contains(node.Data, "\n") && len(node.Data) > 10 {
			if node.Parent.Data == "span" && chunks == 1 {
				errors += fmt.Sprintf("%s\n", node.Data)
			} else if node.Type == html.TextNode && node.Data == "expand_less" {
				chunks += 1
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	return &BuildLog{Error: errors}, nil
}

// extractBuildLensURL returns the final URL for next phase build lens logs extraction.
func extractBuildLensURL(body io.Reader) (string, error) {
	var (
		content  string
		buildLen string
	)

	// Parse the HTML body
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		// Search the first script data and save on a local variable
		if node.Data == "script" && node.FirstChild != nil && strings.Contains(node.FirstChild.Data, "lensArtifacts") {
			content = node.FirstChild.Data
		}
		// Set the iframe index for extracting the script data
		if node.DataAtom == 195590 && node.Type == html.ElementNode {
			if getAttributeValue(node.Attr, "data-lens-name") == "buildlog" {
				buildLen = getAttributeValue(node.Attr, "data-lens-index")
				return
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)

	// Extract JSON path from script lensArtifacts
	var jsonData map[string][]string
	artifacts := GetRegexParameter(`var lensArtifacts = (?<JSON>{"\d+":\[(?:"[^"]+",?)*\](?:,\s*"\d+":\[(?:"[^"]+",?)*\])*\});`, content)["JSON"]
	if err := json.Unmarshal([]byte(artifacts), &jsonData); err != nil {
		return "", err
	}

	// Extract artifacts list from specific lensArtifacts buildlog index.
	var artifactsList []byte
	if artifactsList, err = json.Marshal(jsonData[buildLen]); err != nil {
		return "", err
	}

	// Returns the final buildlog URL rendered in the iframe.
	gsURL := GetRegexParameter(`var src = "(?<URL>[^"]+)"`, content)["URL"]
	data := fmt.Sprintf(`{"artifacts": %s,"index": %s,"src": "%s"}`, artifactsList, buildLen, gsURL)
	return URL + "/spyglass/lens/buildlog/iframe?req=" + url.QueryEscape(data), nil
}

func getHTTPResponse(url string) (io.Reader, error) {
	var (
		response *http.Response
		err      error
	)
	if response, err = http.Get(url); err != nil {
		return nil, err
	}
	return response.Body, nil
}

func getAttributeValue(attrs []html.Attribute, name string) string {
	for _, attr := range attrs {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}
