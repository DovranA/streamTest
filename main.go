package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func videoHandler(c *fiber.Ctx) error {
	videoPath := filepath.Join("videos", "tmbiz-1728563552776.mp4")

	file, err := os.Open(videoPath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not get file info")
	}

	fileSize := stat.Size()

	rangeHeader := c.Get("Range")
	if rangeHeader != "" {
		rangeParts := strings.Split(strings.Replace(rangeHeader, "bytes=", "", 1), "-")
		start, err := strconv.ParseInt(rangeParts[0], 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid range")
		}

		var end int64
		if rangeParts[1] != "" {
			end, err = strconv.ParseInt(rangeParts[1], 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).SendString("Invalid range")
			}
		} else {
			end = fileSize - 1
		}

		chunkSize := end - start + 1

		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
		c.Set("Accept-Ranges", "bytes")
		c.Set("Content-Length", fmt.Sprintf("%d", chunkSize))
		c.Set("Content-Type", "video/mp4")
		c.Status(fiber.StatusPartialContent)

		buf := make([]byte, chunkSize)
		file.Seek(start, 0)
		file.Read(buf)
		return c.Send(buf)
	}

	c.Set("Content-Length", fmt.Sprintf("%d", fileSize))
	c.Set("Content-Type", "video/mp4")

	return c.SendFile(videoPath)
}

func main() {
	app := fiber.New()

	app.Get("/video", videoHandler)

	app.Listen(":8080")
}
