package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/abrshDev/insta-scrapper/scrape"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	httpClient := &http.Client{} // simple client, no proxy

	app.Use(cors.New())

	// Proxy-image endpoint (optional)
	app.Get("/proxy-image", func(c *fiber.Ctx) error {
		imgUrl := c.Query("url")
		if imgUrl == "" {
			return c.Status(400).SendString("Missing url parameter")
		}
		parsedUrl, err := url.QueryUnescape(imgUrl)
		if err != nil {
			return c.Status(400).SendString("Invalid url parameter")
		}
		resp, err := httpClient.Get(parsedUrl)
		if err != nil {
			return c.Status(500).SendString("Failed to fetch image")
		}
		defer resp.Body.Close()
		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		c.Set("Access-Control-Allow-Origin", "*")
		_, err = io.Copy(c, resp.Body)
		return err
	})

	// Images endpoint
	app.Get("/images/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")
		info, err := scrape.GetIGProfileInfo(username) // Scrape.do proxy used here
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(), // propagate actual error
			})
		}

		imgbbKey := os.Getenv("IMGBB_KEY")
		var imgbbLinks []string

		// Upload posts
		for _, imgURL := range info.Images {
			resp, err := httpClient.Get(imgURL)
			if err != nil {
				continue
			}
			imgBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			data := url.Values{}
			data.Set("key", imgbbKey)
			data.Set("image", base64.StdEncoding.EncodeToString(imgBytes))

			req, _ := http.NewRequest("POST", "https://api.imgbb.com/1/upload", bytes.NewBufferString(data.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			res, err := httpClient.Do(req)
			if err != nil {
				continue
			}
			var result map[string]interface{}
			json.NewDecoder(res.Body).Decode(&result)
			res.Body.Close()

			if data, ok := result["data"].(map[string]interface{}); ok {
				if url, ok := data["url"].(string); ok {
					imgbbLinks = append(imgbbLinks, url)
				}
			}
		}

		// Upload profile image
		var profileImgLink string
		if info.ProfileImage != "" {
			resp, err := httpClient.Get(info.ProfileImage)
			if err == nil {
				imgBytes, _ := io.ReadAll(resp.Body)
				resp.Body.Close()

				data := url.Values{}
				data.Set("key", imgbbKey)
				data.Set("image", base64.StdEncoding.EncodeToString(imgBytes))

				req, _ := http.NewRequest("POST", "https://api.imgbb.com/1/upload", bytes.NewBufferString(data.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				res, err := httpClient.Do(req)
				if err == nil {
					var result map[string]interface{}
					json.NewDecoder(res.Body).Decode(&result)
					res.Body.Close()
					if data, ok := result["data"].(map[string]interface{}); ok {
						if url, ok := data["url"].(string); ok {
							profileImgLink = url
						}
					}
				}
			}
		}

		return c.JSON(fiber.Map{
			"username":      username,
			"images":        imgbbLinks,
			"profile_image": profileImgLink,
			"followers":     info.Followers,
		})
	})

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
