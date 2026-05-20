package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sulthancom/pos-backend/internal/database"
	"github.com/sulthancom/pos-backend/internal/models"
	"github.com/sulthancom/pos-backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServiceHandler struct {
	db *database.Database
}

func NewServiceHandler(db *database.Database) *ServiceHandler {
	return &ServiceHandler{db: db}
}

func (h *ServiceHandler) GetAll(c *fiber.Ctx) error {
	page, limit := utils.ParsePaginationQuery(c.Query("page"), c.Query("limit"))
	search := c.Query("search")
	status := c.Query("status")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"ticket_number": primitive.Regex{Pattern: search, Options: "i"}},
			{"customer_name": primitive.Regex{Pattern: search, Options: "i"}},
			{"complaint": primitive.Regex{Pattern: search, Options: "i"}},
		}
	}
	if status != "" {
		filter["status"] = status
	}

	total, err := h.db.Collection("services").CountDocuments(ctx, filter)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to count services")
	}

	var services []models.ServiceTicket
	cursor, err := h.db.Collection("services").Find(ctx, filter, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch services")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &services); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode services")
	}

	return utils.RespondPaginated(c, services, total, page, limit)
}

func (h *ServiceHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid service ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var service models.ServiceTicket
	err = h.db.Collection("services").FindOne(ctx, bson.M{"_id": objectID}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Service ticket not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch service")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", service)
}

func (h *ServiceHandler) Create(c *fiber.Ctx) error {
	var service models.ServiceTicket
	if err := c.BodyParser(&service); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	service.ID = primitive.NewObjectID()
	service.TicketNumber = utils.GenerateServiceNumber()
	service.Status = models.ServiceWaiting
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	service.SparepartsUsed = []models.ServiceSparepart{}
	service.TechnicianNotes = []models.TechnicianNote{}
	service.Photos = []string{}

	technicianID := c.Locals("user_id").(string)
	technicianName := c.Locals("user_name").(string)
	techOID, _ := primitive.ObjectIDFromHex(technicianID)
	service.TechnicianID = techOID
	service.TechnicianName = technicianName

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("services").InsertOne(ctx, service)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create service ticket")
	}

	service.ID = result.InsertedID.(primitive.ObjectID)
	return utils.RespondSuccess(c, fiber.StatusCreated, "Service ticket created successfully", service)
}

func (h *ServiceHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid service ID")
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	updates["updated_at"] = time.Now()
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "ticket_number")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("services").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update service")
	}

	if result.MatchedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Service ticket not found")
	}

	return h.GetByID(c)
}

func (h *ServiceHandler) AddNote(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid service ID")
	}

	var req struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil || req.Note == "" {
		return utils.RespondError(c, fiber.StatusBadRequest, "Note is required")
	}

	note := models.TechnicianNote{
		ID:        primitive.NewObjectID(),
		Note:      req.Note,
		CreatedAt: time.Now(),
		CreatedBy: c.Locals("user_name").(string),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.db.Collection("services").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{
		"$push": bson.M{"technician_notes": note},
		"$set":  bson.M{"updated_at": time.Now()},
	})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to add note")
	}

	return h.GetByID(c)
}

func (h *ServiceHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid service ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("services").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to delete service")
	}

	if result.DeletedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Service ticket not found")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Service ticket deleted successfully", nil)
}
