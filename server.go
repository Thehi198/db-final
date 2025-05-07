package main

import (
	"log"
	"vectordb/vectordb"

	"github.com/gofiber/fiber/v2"
)

var db *vectordb.Database

type createReq struct {
	Dimension int               `json:"dimension"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}
type insertReq struct {
	Values   []float64         `json:"values"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
type queryReq struct {
	Values         []float64         `json:"values"`
	K              int               `json:"k"`
	MetadataFilter map[string]string `json:"metadata_filter,omitempty"`
}
type insertResp struct {
	UUID string `json:"uuid"`
}
type queryByUUIDReq struct {
	UUID string `json:"uuid"`
}
type updateReq struct {
	UUID     string            `json:"uuid"`
	Values   []float64         `json:"values"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
type deleteReq struct {
	UUID string `json:"uuid"`
}

func serve() {
	app := fiber.New()

	app.Post("/create", func(c *fiber.Ctx) error {
		var req createReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		db = vectordb.NewDatabase(req.Dimension)
		db.Metadata = req.Metadata
		return c.JSON(fiber.Map{"status": "ok", "dimension": req.Dimension, "metadata": req.Metadata})
	})

	app.Post("/insert", func(c *fiber.Ctx) error {
		if db == nil {
			return c.Status(400).JSON(fiber.Map{"error": "db not initialized"})
		}
		var req insertReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		id, err := db.Insert(vectordb.Vector{Values: req.Values, Metadata: req.Metadata})
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(insertResp{UUID: id})
	})

	app.Post("/query", func(c *fiber.Ctx) error {
		if db == nil {
			return c.Status(400).JSON(fiber.Map{"error": "db not initialized"})
		}
		var req queryReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		res, err := db.Query(req.Values, req.K, req.MetadataFilter)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(res)
	})

	app.Post("/query_uuid", func(c *fiber.Ctx) error {
		if db == nil {
			return c.Status(400).JSON(fiber.Map{"error": "db not initialized"})
		}
		var req queryByUUIDReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		v, err := db.QueryByUUID(req.UUID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(v)
	})

	app.Put("/update", func(c *fiber.Ctx) error {
		var req updateReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.Update(req.UUID, req.Values, req.Metadata); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "updated"})
	})

	app.Delete("/delete", func(c *fiber.Ctx) error {
		var req deleteReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := db.Delete(req.UUID); err != nil {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"status": "deleted"})
	})

	app.Get("/list", func(c *fiber.Ctx) error {
		if db == nil {
			return c.Status(400).JSON(fiber.Map{"error": "db not initialized"})
		}
		return c.JSON(db.Vectors)
	})

	log.Fatal(app.Listen(":3000"))
}

func main() {
	serve()
}
