package handlers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sulthancom/pos-backend/internal/database"
	"github.com/sulthancom/pos-backend/internal/models"
	"github.com/sulthancom/pos-backend/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportHandler struct {
	db *database.Database
}

func NewReportHandler(db *database.Database) *ReportHandler {
	return &ReportHandler{db: db}
}

func (h *ReportHandler) SalesReport(c *fiber.Ctx) error {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	groupBy := c.Query("group_by", "day")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter["created_at"] = bson.M{"$gte": t}
		}
	}
	if dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			if existing, ok := filter["created_at"].(bson.M); ok {
				existing["$lte"] = t
			} else {
				filter["created_at"] = bson.M{"$lte": t}
			}
		}
	}

	var dateField string
	switch groupBy {
	case "month":
		dateField = "%Y-%m"
	case "year":
		dateField = "%Y"
	default:
		dateField = "%Y-%m-%d"
	}

	var results []bson.M
	cursor, err := h.db.Collection("transactions").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":             bson.M{"$dateToString": bson.M{"format": dateField, "date": "$created_at"}},
			"total_sales":     bson.M{"$sum": "$total"},
			"total_discount":  bson.M{"$sum": "$discount"},
			"total_tax":       bson.M{"$sum": "$tax"},
			"transaction_count": bson.M{"$sum": 1},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"_id": 1}}},
	})
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "Failed to generate sales report")
	}
	defer cursor.Close(ctx)

	cursor.All(ctx, &results)

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", results)
}

func (h *ReportHandler) ProductReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var bestSelling []bson.M
	cursor, _ := h.db.Collection("stock_movements").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"type": "out"}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":            "$product_id",
			"product_name":   bson.M{"$first": "$product_name"},
			"total_sold":     bson.M{"$sum": bson.M{"$abs": "$quantity"}},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"total_sold": -1}}},
		bson.D{{Key: "$limit", Value: 10}},
	})
	cursor.All(ctx, &bestSelling)

	var lowStock []models.Product
	cursor3, _ := h.db.Collection("products").Find(ctx, bson.M{"$expr": bson.M{"$lte": []string{"$stock", "$min_stock"}}})
	cursor3.All(ctx, &lowStock)

	var categoryStats []bson.M
	cursor2, _ := h.db.Collection("products").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.M{
			"_id":       "$category_id",
			"count":     bson.M{"$sum": 1},
			"total_stock": bson.M{"$sum": "$stock"},
		}}},
	})
	cursor2.All(ctx, &categoryStats)

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", bson.M{
		"best_selling":  bestSelling,
		"low_stock":     lowStock,
		"category_stats": categoryStats,
	})
}

func (h *ReportHandler) ServiceReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var statusStats []bson.M
	cursor, _ := h.db.Collection("services").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}}},
	})
	cursor.All(ctx, &statusStats)

	var technicianStats []bson.M
	cursor2, _ := h.db.Collection("services").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.M{
			"_id":            "$technician_id",
			"technician_name": bson.M{"$first": "$technician_name"},
			"total_services":  bson.M{"$sum": 1},
			"completed": bson.M{"$sum": bson.M{"$cond": bson.M{
				"if":   bson.M{"$eq": []string{"$status", "completed"}},
				"then": 1,
				"else": 0,
			}}},
		}}},
	})
	cursor2.All(ctx, &technicianStats)

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", bson.M{
		"status_stats":    statusStats,
		"technician_stats": technicianStats,
	})
}

func (h *ReportHandler) FinancialReport(c *fiber.Ctx) error {
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	if dateFrom != "" {
		if t, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter["created_at"] = bson.M{"$gte": t}
		}
	}
	if dateTo != "" {
		if t, err := time.Parse("2006-01-02", dateTo); err == nil {
			if existing, ok := filter["created_at"].(bson.M); ok {
				existing["$lte"] = t
			} else {
				filter["created_at"] = bson.M{"$lte": t}
			}
		}
	}

	var totalIncome float64
	cursor, _ := h.db.Collection("transactions").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$paid_amount"}}}},
	})
	if cursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		cursor.Decode(&result)
		totalIncome = result.Total
	}

	var totalReceivables float64
	cursor2, _ := h.db.Collection("receivables").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"status": bson.M{"$in": []string{"unpaid", "partial"}}}}},
		bson.D{{Key: "$group", Value: bson.M{"_id": nil, "total": bson.M{"$sum": "$remaining_amount"}}}},
	})
	if cursor2.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		cursor2.Decode(&result)
		totalReceivables = result.Total
	}

	var paymentMethodStats []bson.M
	cursor3, _ := h.db.Collection("transactions").Aggregate(ctx, mongo.Pipeline{
		bson.D{{Key: "$group", Value: bson.M{
			"_id":   "$payment_method",
			"count": bson.M{"$sum": 1},
			"total": bson.M{"$sum": "$total"},
		}}},
	})
	cursor3.All(ctx, &paymentMethodStats)

	return utils.RespondSuccess(c, fiber.StatusOK, "Success", bson.M{
		"total_income":      totalIncome,
		"total_receivables": totalReceivables,
		"payment_methods":   paymentMethodStats,
	})
}
