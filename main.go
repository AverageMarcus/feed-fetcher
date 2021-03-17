package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/mmcdole/gofeed"
)

func main() {
	fp := gofeed.NewParser()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	app := fiber.New(fiber.Config{})

	app.Get("/", func(c *fiber.Ctx) error {
		feedUrl := c.Query("url")
		if feedUrl == "" {
			fmt.Println("No URL provided")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		_, err := fp.ParseURL(feedUrl)
		if err != nil && err == gofeed.ErrFeedTypeNotDetected {
			res, err := http.Get(feedUrl)
			if err != nil {
				fmt.Println("Failed to fetch URL")
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			defer res.Body.Close()
			if res.StatusCode >= 400 {
				fmt.Println("Provided URL returned an error status code")
				return c.SendStatus(res.StatusCode)
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				fmt.Println("Failed to parse response body")
				return c.SendStatus(fiber.StatusInternalServerError)
			}

			matches := doc.Find(`[rel="alternate"][type="application/rss+xml"]`)
			if matches.Length() == 0 {
				fmt.Println("No RSS feeds found on page")
				return c.SendStatus(fiber.StatusNotFound)
			}

			foundUrl, ok := matches.First().Attr("href")
			if !ok {
				fmt.Println("href attribute missing from tag")
				return c.SendStatus(fiber.StatusNotFound)
			}
			c.Set("Location", absoluteUrl(feedUrl, foundUrl))
			if matches.Length() > 1 {
				fmt.Println("Multiple feeds found on page")
				return c.SendStatus(fiber.StatusMultipleChoices)
			} else {
				fmt.Println("Feed found on page")
				return c.SendStatus(fiber.StatusTemporaryRedirect)
			}
		} else if err != nil {
			fmt.Println("Failed while attempting to parse feed")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		fmt.Println("URL provided is already a feed")
		c.Set("Location", feedUrl)
		return c.SendStatus(fiber.StatusMovedPermanently)
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
