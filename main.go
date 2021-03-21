package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"embed"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/mmcdole/gofeed"
)

//go:embed index.html

var content embed.FS

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	app := fiber.New(fiber.Config{})

	app.Get("/", func(c *fiber.Ctx) error {
		feedUrl := c.Query("url")
		if feedUrl == "" {
			c.Type("html", "UTF8")
			body, _ := content.ReadFile("index.html")
			return c.Send(body)
		}

		feeds, statusCode := getFeeds(feedUrl)

		if c.Is("json") {
			if statusCode >= 400 || statusCode == 300 {
				c.Status(statusCode)
			}
			return c.JSON(feeds)
		} else {
			c.Status(statusCode)
			c.Location(feeds[0])

			if len(feeds) > 1 {
				responseBody := "Multiple Choices\n\n"
				for _, feed := range feeds {
					responseBody += feed + "\n"
				}
				return c.SendString(responseBody)
			}

			return c.Send(nil)
		}
	})

	fmt.Println(app.Listen(fmt.Sprintf(":%s", port)))
}

func absoluteUrl(requestUrl, foundUrl string) string {
	if !strings.HasPrefix(foundUrl, "http") {
		parsedUrl, _ := url.Parse(requestUrl)
		foundUrl = fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, foundUrl)
	}

	return foundUrl
}

func getFeeds(requestURL string) ([]string, int) {
	feeds := []string{}

	fp := gofeed.NewParser()
	_, err := fp.ParseURL(requestURL)
	if err == nil {
		feeds = []string{requestURL}
	} else if err != nil && err == gofeed.ErrFeedTypeNotDetected {
		res, err := http.Get(requestURL)
		if err != nil {
			fmt.Println("Failed to fetch URL")
			return feeds, fiber.StatusInternalServerError
		}
		defer res.Body.Close()

		if res.StatusCode >= 400 {
			fmt.Println("Provided URL returned an error status code")
			return feeds, res.StatusCode
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			fmt.Println("Failed to parse response body")
			return feeds, fiber.StatusInternalServerError
		}

		matches := doc.Find(`[rel="alternate"][type="application/rss+xml"]`)
		if matches.Length() == 0 {
			fmt.Println("No RSS feeds found on page")
			return feeds, fiber.StatusNotFound
		}

		matches.Each(func(i int, s *goquery.Selection) {
			feeds = append(feeds, absoluteUrl(requestURL, s.AttrOr("href", "")))
		})

		if matches.Length() > 1 {
			fmt.Println("Multiple feeds found on page")
			return feeds, fiber.StatusMultipleChoices
		} else {
			fmt.Println("Feed found on page")
			return feeds, fiber.StatusTemporaryRedirect
		}
	} else if err != nil {
		fmt.Println("Failed while attempting to parse feed")
		return feeds, fiber.StatusInternalServerError
	}

	return feeds, 200
}
