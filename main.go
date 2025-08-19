package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
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
		log.Println("[Proxy-Image] Request URL:", imgUrl)

		if imgUrl == "" {
			log.Println("[Proxy-Image] Missing URL parameter")
			return c.Status(400).SendString("Missing url parameter")
		}

		parsedUrl, err := url.QueryUnescape(imgUrl)
		if err != nil {
			log.Println("[Proxy-Image] Invalid URL:", err)
			return c.Status(400).SendString("Invalid url parameter")
		}

		resp, err := httpClient.Get(parsedUrl)
		if err != nil {
			log.Println("[Proxy-Image] Failed to fetch image:", err)
			return c.Status(500).SendString("Failed to fetch image")
		}
		defer resp.Body.Close()

		c.Set("Content-Type", resp.Header.Get("Content-Type"))
		c.Set("Access-Control-Allow-Origin", "*")
		_, err = io.Copy(c, resp.Body)
		if err != nil {
			log.Println("[Proxy-Image] Error sending image:", err)
		}
		return err
	})

	// Images endpoint
	app.Get("/images/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")
		log.Println("[Images] Scraping username:", username)

		info, err := scrape.GetIGProfileInfo(username)
		if err != nil {
			log.Println("[Images] Error scraping profile:", err)
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		imgbbKey := os.Getenv("IMGBB_KEY")
		var imgbbLinks []string

		for _, imgURL := range info.Images {
			log.Println("[Images] Uploading image:", imgURL)
			resp, err := httpClient.Get(imgURL)
			if err != nil {
				log.Println("[Images] Failed to fetch image:", err)
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
				log.Println("[Images] Failed to upload image:", err)
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

		var profileImgLink string
		if info.ProfileImage != "" {
			log.Println("[Images] Uploading profile image")
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
			} else {
				log.Println("[Images] Failed to fetch profile image:", err)
			}
		}

		log.Println("[Images] Finished scraping for:", username)
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
	log.Println("Server starting on port", port)
	app.Listen(":" + port)
}
