package handlers

import (
	"context"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/bytebufferpool"
)

func RenderComponent(c *fiber.Ctx, com templ.Component) error {
	// Get new buffer from pool
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	// Renders component en puts resulting html into buf
	err := com.Render(context.Background(), buf)
	if err != nil {
		return err
	}

	// Set Content-Type to text/html
	c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	// Set rendered template to body
	c.Context().SetBody(buf.Bytes())
	return nil
}
