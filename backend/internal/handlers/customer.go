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

type CustomerHandler struct {
	db *database.Database
}

func NewCustomerHandler(db *database.Database) *CustomerHandler {
	return &CustomerHandler{db: db}
}

func (h *CustomerHandler) GetAll(c *fiber.Ctx) error {
	page, limit := utils.ParsePaginationQuery(c.Query("page"), c.Query("limit"))
	search := c.Query("search")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": primitive.Regex{Pattern: search, Options: "i"}},
			{"phone": primitive.Regex{Pattern: search, Options: "i"}},
			{"email": primitive.Regex{Pattern: search, Options: "i"}},
		}
	}

	total, err := h.db.Collection("customers").CountDocuments(ctx, filter)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to count customers")
	}

	var customers []models.Customer
	cursor, err := h.db.Collection("customers").Find(ctx, filter, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch customers")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &customers); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode customers")
	}

	return utils.RespondPaginated(c, customers, total, page, limit)
}

func (h *CustomerHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid customer ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var customer models.Customer
	err = h.db.Collection("customers").FindOne(ctx, bson.M{"_id": objectID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Customer not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch customer")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", customer)
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var customer models.Customer
	if err := c.BodyParser(&customer); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	customer.ID = primitive.NewObjectID()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("customers").InsertOne(ctx, customer)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create customer")
	}

	customer.ID = result.InsertedID.(primitive.ObjectID)
	return utils.RespondSuccess(c, fiber.StatusCreated, "Customer created successfully", customer)
}

func (h *CustomerHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid customer ID")
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	updates["updated_at"] = time.Now()
	delete(updates, "id")
	delete(updates, "created_at")
	delete(updates, "total_transactions")
	delete(updates, "total_spent")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("customers").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": updates})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to update customer")
	}

	if result.MatchedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Customer not found")
	}

	return h.GetByID(c)
}

func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid customer ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := h.db.Collection("customers").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to delete customer")
	}

	if result.DeletedCount == 0 {
		return utils.RespondError(c, fiber.StatusNotFound, "Customer not found")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Customer deleted successfully", nil)
}

type ReceivableHandler struct {
	db *database.Database
}

func NewReceivableHandler(db *database.Database) *ReceivableHandler {
	return &ReceivableHandler{db: db}
}

func (h *ReceivableHandler) GetAll(c *fiber.Ctx) error {
	page, limit := utils.ParsePaginationQuery(c.Query("page"), c.Query("limit"))
	search := c.Query("search")
	status := c.Query("status")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"invoice_number": primitive.Regex{Pattern: search, Options: "i"}},
			{"customer_name": primitive.Regex{Pattern: search, Options: "i"}},
		}
	}
	if status != "" {
		filter["status"] = status
	}

	total, err := h.db.Collection("receivables").CountDocuments(ctx, filter)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to count receivables")
	}

	var receivables []models.Receivable
	cursor, err := h.db.Collection("receivables").Find(ctx, filter, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)).SetSort(bson.M{"due_date": 1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch receivables")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &receivables); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode receivables")
	}

	return utils.RespondPaginated(c, receivables, total, page, limit)
}

func (h *ReceivableHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid receivable ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var receivable models.Receivable
	err = h.db.Collection("receivables").FindOne(ctx, bson.M{"_id": objectID}).Decode(&receivable)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Receivable not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch receivable")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", receivable)
}

func (h *ReceivableHandler) MakePayment(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid receivable ID")
	}

	var req struct {
		Amount float64             `json:"amount"`
		Method models.PaymentMethod `json:"method"`
		Notes  string              `json:"notes"`
	}

	if err := c.BodyParser(&req); err != nil || req.Amount <= 0 {
		return utils.RespondError(c, fiber.StatusBadRequest, "Valid amount is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var receivable models.Receivable
	err = h.db.Collection("receivables").FindOne(ctx, bson.M{"_id": objectID}).Decode(&receivable)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Receivable not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch receivable")
	}

	if req.Amount > receivable.RemainingAmount {
		return utils.RespondError(c, fiber.StatusBadRequest, "Payment amount exceeds remaining balance")
	}

	payment := models.ReceivablePayment{
		ID:     primitive.NewObjectID(),
		Amount: req.Amount,
		Method: req.Method,
		Date:   time.Now(),
		Notes:  req.Notes,
	}

	receivable.PaidAmount += req.Amount
	receivable.RemainingAmount -= req.Amount
	receivable.Payments = append(receivable.Payments, payment)
	receivable.UpdatedAt = time.Now()

	if receivable.RemainingAmount <= 0 {
		receivable.Status = models.ReceivablePaid
	} else {
		receivable.Status = models.ReceivablePartial
	}

	_, err = h.db.Collection("receivables").UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{
		"paid_amount":      receivable.PaidAmount,
		"remaining_amount": receivable.RemainingAmount,
		"status":           receivable.Status,
		"payments":         receivable.Payments,
		"updated_at":       receivable.UpdatedAt,
	}})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to process payment")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Payment processed successfully", receivable)
}

func (h *ReceivableHandler) GetStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := h.db.Collection("receivables").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
			"total": bson.M{"$sum": "$remaining_amount"},
		}}},
	})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch stats")
	}
	defer cursor.Close(ctx)

	var result []bson.M
	cursor.All(ctx, &result)

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", result)
}
