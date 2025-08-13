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
	// Proxy-image endpoint to serve any image with CORS headers
	app.Get("/proxy-image", func(c *fiber.Ctx) error {
		imgUrl := c.Query("url")
		if imgUrl == "" {
			return c.Status(400).SendString("Missing url parameter")
		}
		parsedUrl, err := url.QueryUnescape(imgUrl)
		if err != nil {
			return c.Status(400).SendString("Invalid url parameter")
		}
		resp, err := http.Get(parsedUrl)
		if err != nil {
			return c.Status(500).SendString("Failed to fetch image")
		}
		defer resp.Body.Close()
		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		c.Set("Access-Control-Allow-Origin", "*")
		_, err = io.Copy(c, resp.Body)
		return err
	})

	// Enable CORS for all origins
	app.Use(cors.New())

	// Proxy endpoint to serve images
	app.Get("/proxy", func(c *fiber.Ctx) error {
		imgUrl := c.Query("url")
		if imgUrl == "" {
			return c.Status(400).SendString("Missing url parameter")
		}
		parsedUrl, err := url.QueryUnescape(imgUrl)
		if err != nil {
			return c.Status(400).SendString("Invalid url parameter")
		}
		resp, err := http.Get(parsedUrl)
		if err != nil {
			return c.Status(500).SendString("Failed to fetch image")
		}
		defer resp.Body.Close()
		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		_, err = io.Copy(c, resp.Body)
		return err
	})

	app.Get("/images/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")
		info, err := scrape.GetIGProfileInfo(username)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"Error": "failed to scrape profile",
			})
		}

		imgbbKey := "904775b3a745b64f07d3f6dff7407701" // Provided API key
		var imgbbLinks []string
		for _, imgURL := range info.Images {
			// Download image
			resp, err := http.Get(imgURL)
			if err != nil {
				continue
			}
			imgBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				continue
			}
			// Encode image to base64
			imgBase64 := base64.StdEncoding.EncodeToString(imgBytes)

			// Upload to imgbb
			data := url.Values{}
			data.Set("key", imgbbKey)
			data.Set("image", imgBase64)
			req, err := http.NewRequest("POST", "https://api.imgbb.com/1/upload", bytes.NewBufferString(data.Encode()))
			if err != nil {
				continue
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			client := &http.Client{}
			res, err := client.Do(req)
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

		// Upload profile image to imgbb
		var profileImgLink string
		if info.ProfileImage != "" {
			resp, err := http.Get(info.ProfileImage)
			if err == nil {
				imgBytes, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err == nil {
					imgBase64 := base64.StdEncoding.EncodeToString(imgBytes)
					data := url.Values{}
					data.Set("key", imgbbKey)
					data.Set("image", imgBase64)
					req, err := http.NewRequest("POST", "https://api.imgbb.com/1/upload", bytes.NewBufferString(data.Encode()))
					if err == nil {
						req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
						client := &http.Client{}
						res, err := client.Do(req)
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
			}
		}

		return c.JSON(fiber.Map{
			"username":      username,
			"images":        imgbbLinks,
			"profile_image": profileImgLink,
			"followers":     info.Followers,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // fallback for local dev
	}
	app.Listen(":" + port)
}
