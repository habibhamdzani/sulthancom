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

type TransactionHandler struct {
	db *database.Database
}

func NewTransactionHandler(db *database.Database) *TransactionHandler {
	return &TransactionHandler{db: db}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var tx models.Transaction
	if err := c.BodyParser(&tx); err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	tx.ID = primitive.NewObjectID()
	tx.InvoiceNumber = utils.GenerateInvoiceNumber()
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()

	cashierID := c.Locals("user_id").(string)
	cashierName := c.Locals("user_name").(string)
	cashierOID, _ := primitive.ObjectIDFromHex(cashierID)
	tx.CashierID = cashierOID
	tx.CashierName = cashierName

	if tx.PaymentStatus == models.PaymentPaid {
		tx.PaidAmount = tx.Total
		tx.RemainingAmount = 0
	} else if tx.PaymentStatus == models.PaymentDP {
		if tx.DPAmount == nil {
			return utils.RespondError(c, fiber.StatusBadRequest, "DP amount is required")
		}
		tx.PaidAmount = *tx.DPAmount
		tx.RemainingAmount = tx.Total - tx.PaidAmount
	} else if tx.PaymentStatus == models.PaymentUnpaid {
		tx.PaidAmount = 0
		tx.RemainingAmount = tx.Total
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := h.db.Client.StartSession()
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to start session")
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		_, err := h.db.Collection("transactions").InsertOne(sc, tx)
		if err != nil {
			return nil, err
		}

		for _, item := range tx.Items {
			if item.Quantity > 0 {
				var product models.Product
				err := h.db.Collection("products").FindOne(sc, bson.M{"_id": item.ProductID}).Decode(&product)
				if err == nil {
					previousStock := product.Stock
					product.Stock -= item.Quantity
					h.db.Collection("products").UpdateOne(sc, bson.M{"_id": item.ProductID}, bson.M{"$set": bson.M{"stock": product.Stock, "updated_at": time.Now()}})

					movement := models.StockMovement{
						ID:            primitive.NewObjectID(),
						ProductID:     item.ProductID,
						ProductName:   item.Name,
						Type:          models.StockOut,
						Quantity:      item.Quantity,
						PreviousStock: previousStock,
						NewStock:      product.Stock,
						Notes:         "Terjual - " + tx.InvoiceNumber,
						CreatedAt:     time.Now(),
						CreatedBy:     cashierName,
					}
					h.db.Collection("stock_movements").InsertOne(sc, movement)
				}
			}
		}

		if tx.PaymentStatus == models.PaymentUnpaid || tx.PaymentStatus == models.PaymentDP {
			if tx.CustomerID != nil && tx.DueDate != nil {
				receivable := models.Receivable{
					ID:              primitive.NewObjectID(),
					TransactionID:   tx.ID,
					InvoiceNumber:   tx.InvoiceNumber,
					CustomerID:      *tx.CustomerID,
					CustomerName:    tx.CustomerName,
					TotalAmount:     tx.Total,
					PaidAmount:      tx.PaidAmount,
					RemainingAmount: tx.RemainingAmount,
					DueDate:         *tx.DueDate,
					Status:          models.ReceivablePartial,
					Payments:        []models.ReceivablePayment{},
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				if tx.PaidAmount == 0 {
					receivable.Status = models.ReceivableUnpaid
				}
				if tx.DPAmount != nil {
					receivable.Payments = append(receivable.Payments, models.ReceivablePayment{
						ID:     primitive.NewObjectID(),
						Amount: *tx.DPAmount,
						Method: tx.PaymentMethod,
						Date:   time.Now(),
						Notes:  "DP",
					})
				}
				h.db.Collection("receivables").InsertOne(sc, receivable)
			}
		}

		if tx.CustomerID != nil {
			h.db.Collection("customers").UpdateOne(sc, bson.M{"_id": *tx.CustomerID}, bson.M{
				"$inc": bson.M{"total_transactions": 1, "total_spent": tx.Total},
			})
		}

		return nil, nil
	})

	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to create transaction")
	}

	return utils.RespondSuccess(c, fiber.StatusCreated, "Transaction created successfully", tx)
}

func (h *TransactionHandler) GetAll(c *fiber.Ctx) error {
	page, limit := utils.ParsePaginationQuery(c.Query("page"), c.Query("limit"))
	search := c.Query("search")
	status := c.Query("status")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

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
		filter["payment_status"] = status
	}
	if dateFrom != "" || dateTo != "" {
		dateFilter := bson.M{}
		if dateFrom != "" {
			if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
				dateFilter["$gte"] = t
			}
		}
		if dateTo != "" {
			if t, err := time.Parse("2006-01-02", dateTo); err == nil {
				dateFilter["$lte"] = t
			}
		}
		if len(dateFilter) > 0 {
			filter["created_at"] = dateFilter
		}
	}

	total, err := h.db.Collection("transactions").CountDocuments(ctx, filter)
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to count transactions")
	}

	var transactions []models.Transaction
	cursor, err := h.db.Collection("transactions").Find(ctx, filter, options.Find().SetSkip(int64((page-1)*limit)).SetLimit(int64(limit)).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch transactions")
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &transactions); err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to decode transactions")
	}

	return utils.RespondPaginated(c, transactions, total, page, limit)
}

func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return utils.RespondError(c, fiber.StatusBadRequest, "Invalid transaction ID")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tx models.Transaction
	err = h.db.Collection("transactions").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.RespondError(c, fiber.StatusNotFound, "Transaction not found")
		}
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to fetch transaction")
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", tx)
}

func (h *TransactionHandler) GetStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	var todayStats []bson.M
	cursor, err := h.db.Collection("transactions").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"created_at": bson.M{"$gte": todayStart, "$lt": todayEnd}}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total_sales": bson.M{"$sum": "$total"}, "total_transactions": bson.M{"$sum": 1}}}},
	})
	if err == nil {
		cursor.Decode(&todayStats)
	}

	todaySales := 0.0
	todayTransactions := int64(0)
	if len(todayStats) > 0 {
		todaySales = todayStats[0]["total_sales"].(float64)
		todayTransactions = int64(todayStats[0]["total_transactions"].(int32))
	}

	var activeServices int64
	activeServices, _ = h.db.Collection("services").CountDocuments(ctx, bson.M{"status": bson.M{"$in": []string{"waiting", "checked", "in_progress", "waiting_parts"}}})

	var totalReceivables float64
	cursor2, _ := h.db.Collection("receivables").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"status": bson.M{"$in": []string{"unpaid", "partial"}}}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$remaining_amount"}}}},
	})
	if cursor2.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		cursor.Decode(&result)
		totalReceivables = result.Total
	}

	var totalProducts int64
	totalProducts, _ = h.db.Collection("products").CountDocuments(ctx, bson.M{"is_active": true})

	var lowStock int64
	lowStock, _ = h.db.Collection("products").CountDocuments(ctx, bson.M{"$expr": bson.M{"$lte": []string{"$stock", "$min_stock"}}})

	stats := models.DashboardStats{
		TodaySales:        todaySales,
		TodayTransactions: todayTransactions,
		ActiveServices:    activeServices,
		TotalReceivables:  totalReceivables,
		TotalProducts:     totalProducts,
		LowStockProducts:  lowStock,
	}

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", stats)
}
