package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PreviewVideo represents a preview video for a page
type PreviewVideo struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
	Videos      []*PreviewVideo `json:"videos,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*
	  - Close the response HTML stream so that you don't leak resources.
	  - Finally, respond with a JSON-encoded version of the PageSummary
	    struct. That way the client can easily parse the JSON back into
	    an object. Remember to tell the client that the response content
	    type is JSON.

	  Helpful Links:
	  https://golang.org/pkg/net/http/#Request.FormValue
	  https://golang.org/pkg/net/http/#Error
	  https://golang.org/pkg/encoding/json/#NewEncoder
	*/
	url := r.FormValue("url")
	if url == "" {
		log.Println("Url Param 'url' is missing")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		bodyStream, err0 := fetchHTML(url)

		if err0 != nil {
			http.Error(w, err0.Error(), 400)
			return
		}
		summaryStruct, err1 := extractSummary(url, bodyStream)
		if err1 != nil {
			http.Error(w, err1.Error(), 400)
			return
		}
		bodyStream.Close()
		var jsonData []byte
		jsonData, err3 := json.Marshal(summaryStruct)
		if err3 != nil {
			http.Error(w, err3.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	resp, respErr := http.Get(pageURL)
	if respErr != nil {
		return nil, respErr
	} else if resp.StatusCode >= 400 {
		return nil, errors.New("BadStatusCode")
	} else if !strings.HasPrefix(resp.Header.Get("Content-type"), "text/html") {
		return nil, errors.New("InvalidContentType")
	} else {
		return resp.Body, nil
	}
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*according to the assignment description.

	  To test your implementation of this function, run the TestExtractSummary
	  test in summary_test.go. You can do that directly in Visual Studio Code,
	  or at the command line by running:
	          go test -run TestExtractSummary

	  Helpful Links:
	  https://drstearns.github.io/tutorials/tokenizing/
	  http://ogp.me/
	  https://developers.facebook.com/docs/reference/opengraph/
	  https://golang.org/pkg/net/url/#URL.ResolveReference
	*/
	tokenizer := html.NewTokenizer(htmlStream)
	summaryStruct := new(PageSummary)
	previewIcon := new(PreviewImage)
	isMetaImage := false
	isMetaVideo := false
	//Keep running the for loop to check until there is no more token left
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
			//otherwise, there was an error tokenizing,
			//which likely means the HTML was malformed.
			//since this is a simple command-line utility,
			//we can report the error and exit the process
			//with a non-zero status code
			return nil, tokenizer.Err()
		}
		//get the token
		token := tokenizer.Token()
		//process the token according to the token type...
		//if this is a start tag token...
		if tokenType == html.StartTagToken ||
			tokenType == html.SelfClosingTagToken {
			//if the name of the element is "title"
			if "meta" == token.Data {
				metaProperty := ""
				metaContent := ""
				for _, attr := range token.Attr {
					if "property" == attr.Key {
						metaProperty = attr.Val
					} else if "content" == attr.Key {
						metaContent = attr.Val
					} else if "name" == attr.Key {
						//What is the name in attr.key
						metaProperty = attr.Val
					}
				}
				//check if this is a continued image property that already exists, if it is , directly add to this current imagestruct
				if isMetaImage && strings.HasPrefix(metaProperty, "og:image:") {
					previewImages := summaryStruct.Images
					tempImage := previewImages[len(previewImages)-1]
					imageProperty := metaProperty[9:]
					switch imageProperty {
					case "url":
						url, err := relURLToAbsURL(pageURL, metaContent)
						if err != nil {
							return nil, err
						}
						tempImage.URL = url
					case "secure_url":
						secureUrl, err := relURLToAbsURL(pageURL, metaContent)
						if err != nil {
							return nil, err
						}
						tempImage.SecureURL = secureUrl
					case "type":
						tempImage.Type = metaContent
					case "width":
						widthInt, widthErr := strconv.Atoi(metaContent)
						if widthErr != nil {
							return nil, widthErr
						}
						tempImage.Width = widthInt
					case "height":
						heightInt, heightErr := strconv.Atoi(metaContent)
						if heightErr != nil {
							return nil, heightErr
						}
						tempImage.Height = heightInt
					case "alt":
						tempImage.Alt = metaContent
					}
					continue
				}
				//Check if this a continued video property that belongs to the latest video object, if it is, directly add to this current videoStruct
				if isMetaVideo && strings.HasPrefix(metaProperty, "og:video:") {
					previewVideos := summaryStruct.Videos
					tempVideo := previewVideos[len(previewVideos)-1]
					videoProperty := metaProperty[9:]
					switch videoProperty {
					case "secure_url":
						secureUrl, err := relURLToAbsURL(pageURL, metaContent)
						if err != nil {
							return nil, err
						}
						tempVideo.SecureURL = secureUrl
					case "type":
						tempVideo.Type = metaContent
					case "width":
						widthInt, widthErr := strconv.Atoi(metaContent)
						if widthErr != nil {
							return nil, widthErr
						}
						tempVideo.Width = widthInt
					case "height":
						heightInt, heightErr := strconv.Atoi(metaContent)
						if heightErr != nil {
							return nil, heightErr
						}
						tempVideo.Height = heightInt
					}
					continue
				}
				switch metaProperty {
				case "author":
					isMetaVideo = false
					isMetaImage = false
					summaryStruct.Author = metaContent
				case "og:type":
					isMetaImage = false
					isMetaVideo = false
					summaryStruct.Type = metaContent
				case "og:url":
					isMetaImage = false
					isMetaVideo = false
					url, err := relURLToAbsURL(pageURL, metaContent)
					if err != nil {
						return nil, err
					}
					summaryStruct.URL = url
				case "og:title":
					isMetaImage = false
					isMetaVideo = false
					summaryStruct.Title = metaContent
				case "og:site_name":
					isMetaImage = false
					isMetaVideo = false
					summaryStruct.SiteName = metaContent
				case "og:description":
					isMetaImage = false
					isMetaVideo = false
					summaryStruct.Description = metaContent
				case "description":
					isMetaImage = false
					isMetaVideo = false
					if summaryStruct.Description == "" {
						summaryStruct.Description = metaContent
					}
				case "keywords":
					isMetaImage = false
					isMetaVideo = false
					keywords := strings.Split(metaContent, ",")
					for i := range keywords {
						keywords[i] = strings.TrimSpace(keywords[i])
					}
					summaryStruct.Keywords = keywords
				case "og:image":
					isMetaImage = true
					isMetaVideo = false
					url, err := relURLToAbsURL(pageURL, metaContent)
					if err != nil {
						return nil, err
					}
					if summaryStruct.Images == nil {
						summaryStruct.Images = []*PreviewImage{}
					}
					singleImage := new(PreviewImage)
					singleImage.URL = url
					summaryStruct.Images = append(summaryStruct.Images, singleImage)
				case "og:video":
					isMetaImage = false
					isMetaVideo = true
					url, err := relURLToAbsURL(pageURL, metaContent)
					if err != nil {
						return nil, err
					}
					if summaryStruct.Videos == nil {
						summaryStruct.Videos = []*PreviewVideo{}
					}
					singleVideo := new(PreviewVideo)
					singleVideo.URL = url
					summaryStruct.Videos = append(summaryStruct.Videos, singleVideo)
				}
			} else if "link" == token.Data {
				isIcon := false
				for _, attr := range token.Attr {
					if "rel" == attr.Key {
						if "icon" == attr.Val {
							isIcon = true
						}
					} else if "href" == attr.Key {
						url, err := relURLToAbsURL(pageURL, attr.Val)
						if err != nil {
							return nil, err
						}
						previewIcon.URL = url
					} else if "type" == attr.Key {
						previewIcon.Type = attr.Val
					} else if "sizes" == attr.Key {
						if attr.Val != "any" {
							sizes := strings.Split(attr.Val, "x")
							heightInt, err := strconv.Atoi(sizes[0])
							if err != nil {
								return nil, err
							}
							widthInt, err := strconv.Atoi(sizes[1])
							if err != nil {
								return nil, err
							}
							previewIcon.Width = widthInt
							previewIcon.Height = heightInt
						}
					}
				}
				if isIcon {
					summaryStruct.Icon = previewIcon
				} else {
					isIcon = false
					previewIcon = new(PreviewImage)
				}
			} else if "title" == token.Data {
				//the next token should be the page title
				tokenType = tokenizer.Next()
				//just make sure it's actually a text token
				if tokenType == html.TextToken {
					if summaryStruct.Title == "" {
						summaryStruct.Title = tokenizer.Token().Data
					}
				}
			}
		}
	}
	return summaryStruct, nil
}

func relURLToAbsURL(baseURL string, relURL string) (string, error) {
	parsedRelURL, err := url.Parse(relURL)
	if err != nil {
		return "", err
	}
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	absURL := parsedBaseURL.ResolveReference(parsedRelURL)
	return absURL.String(), nil
}
