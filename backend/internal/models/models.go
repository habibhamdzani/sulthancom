package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleKasir      Role = "kasir"
	RoleTeknisi    Role = "teknisi"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Role      Role               `bson:"role" json:"role"`
	Phone     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Avatar    string             `bson:"avatar,omitempty" json:"avatar,omitempty"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
	User      User   `json:"user"`
}

type Category struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Icon         string             `bson:"icon,omitempty" json:"icon,omitempty"`
	Color        string             `bson:"color,omitempty" json:"color,omitempty"`
	ParentID     *primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	ProductCount int                `bson:"product_count" json:"product_count"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type Brand struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Logo         string             `bson:"logo,omitempty" json:"logo,omitempty"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
	ProductCount int                `bson:"product_count" json:"product_count"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type ProductType string

const (
	ProductTypePhysical ProductType = "physical"
	ProductTypeService  ProductType = "service"
	ProductTypeSparepart ProductType = "sparepart"
	ProductTypeBundle   ProductType = "bundle"
)

type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	SKU           string             `bson:"sku" json:"sku"`
	Barcode       string             `bson:"barcode,omitempty" json:"barcode,omitempty"`
	CategoryID    primitive.ObjectID `bson:"category_id" json:"category_id"`
	BrandID       *primitive.ObjectID `bson:"brand_id,omitempty" json:"brand_id,omitempty"`
	Type          ProductType        `bson:"type" json:"type"`
	CostPrice     float64            `bson:"cost_price" json:"cost_price"`
	SellingPrice  float64            `bson:"selling_price" json:"selling_price"`
	Stock         int                `bson:"stock" json:"stock"`
	MinStock      int                `bson:"min_stock" json:"min_stock"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
	Images        []string           `bson:"images" json:"images"`
	IsActive      bool               `bson:"is_active" json:"is_active"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type Supplier struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Phone     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Email     string             `bson:"email,omitempty" json:"email,omitempty"`
	Address   string             `bson:"address,omitempty" json:"address,omitempty"`
	IsActive  bool               `bson:"is_active" json:"is_active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Customer struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Phone            string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Email            string             `bson:"email,omitempty" json:"email,omitempty"`
	Address          string             `bson:"address,omitempty" json:"address,omitempty"`
	Notes            string             `bson:"notes,omitempty" json:"notes,omitempty"`
	TotalTransactions int               `bson:"total_transactions" json:"total_transactions"`
	TotalSpent       float64            `bson:"total_spent" json:"total_spent"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}

type StockMovementType string

const (
	StockIn        StockMovementType = "in"
	StockOut       StockMovementType = "out"
	StockAdjustment StockMovementType = "adjustment"
	StockTransfer  StockMovementType = "transfer"
)

type StockMovement struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProductID      primitive.ObjectID  `bson:"product_id" json:"product_id"`
	ProductName    string              `bson:"product_name" json:"product_name"`
	Type           StockMovementType   `bson:"type" json:"type"`
	Quantity       int                 `bson:"quantity" json:"quantity"`
	PreviousStock  int                 `bson:"previous_stock" json:"previous_stock"`
	NewStock       int                 `bson:"new_stock" json:"new_stock"`
	Notes          string              `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt      time.Time           `bson:"created_at" json:"created_at"`
	CreatedBy      string              `bson:"created_by" json:"created_by"`
}

type PaymentMethod string

const (
	PaymentCash      PaymentMethod = "cash"
	PaymentDANA      PaymentMethod = "dana"
	PaymentGoPay     PaymentMethod = "gopay"
	PaymentShopeePay PaymentMethod = "shopeepay"
	PaymentBCA       PaymentMethod = "bca"
	PaymentBRI       PaymentMethod = "bri"
	PaymentBNI       PaymentMethod = "bni"
	PaymentMandiri   PaymentMethod = "mandiri"
	PaymentOther     PaymentMethod = "other"
)

type PaymentStatus string

const (
	PaymentPaid     PaymentStatus = "paid"
	PaymentDP       PaymentStatus = "dp"
	PaymentUnpaid   PaymentStatus = "unpaid"
	PaymentPartial  PaymentStatus = "partial"
)

type CartItem struct {
	ProductID  primitive.ObjectID `bson:"product_id" json:"product_id"`
	Name       string             `bson:"name" json:"name"`
	Price      float64            `bson:"price" json:"price"`
	Quantity   int                `bson:"quantity" json:"quantity"`
	Discount   float64            `bson:"discount" json:"discount"`
	Subtotal   float64            `bson:"subtotal" json:"subtotal"`
}

type Transaction struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	InvoiceNumber  string             `bson:"invoice_number" json:"invoice_number"`
	CustomerID     *primitive.ObjectID `bson:"customer_id,omitempty" json:"customer_id,omitempty"`
	CustomerName   string             `bson:"customer_name,omitempty" json:"customer_name,omitempty"`
	Items          []CartItem         `bson:"items" json:"items"`
	Subtotal       float64            `bson:"subtotal" json:"subtotal"`
	Discount       float64            `bson:"discount" json:"discount"`
	Tax            float64            `bson:"tax" json:"tax"`
	Total          float64            `bson:"total" json:"total"`
	PaymentMethod  PaymentMethod      `bson:"payment_method" json:"payment_method"`
	PaymentStatus  PaymentStatus      `bson:"payment_status" json:"payment_status"`
	PaidAmount     float64            `bson:"paid_amount" json:"paid_amount"`
	RemainingAmount float64           `bson:"remaining_amount" json:"remaining_amount"`
	DPAmount       *float64           `bson:"dp_amount,omitempty" json:"dp_amount,omitempty"`
	DueDate        *time.Time         `bson:"due_date,omitempty" json:"due_date,omitempty"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`
	CashierID      primitive.ObjectID `bson:"cashier_id" json:"cashier_id"`
	CashierName    string             `bson:"cashier_name" json:"cashier_name"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type ServiceStatus string

const (
	ServiceWaiting      ServiceStatus = "waiting"
	ServiceChecked      ServiceStatus = "checked"
	ServiceInProgress   ServiceStatus = "in_progress"
	ServiceWaitingParts ServiceStatus = "waiting_parts"
	ServiceCompleted    ServiceStatus = "completed"
	ServicePickedUp     ServiceStatus = "picked_up"
	ServiceCancelled    ServiceStatus = "cancelled"
)

type ServiceSparepart struct {
	ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
	Name      string             `bson:"name" json:"name"`
	Quantity  int                `bson:"quantity" json:"quantity"`
	Price     float64            `bson:"price" json:"price"`
}

type TechnicianNote struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Note      string             `bson:"note" json:"note"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	CreatedBy string             `bson:"created_by" json:"created_by"`
}

type ServiceTicket struct {
	ID                   primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	TicketNumber         string               `bson:"ticket_number" json:"ticket_number"`
	CustomerID           *primitive.ObjectID  `bson:"customer_id,omitempty" json:"customer_id,omitempty"`
	CustomerName         string               `bson:"customer_name" json:"customer_name"`
	CustomerPhone        string               `bson:"customer_phone,omitempty" json:"customer_phone,omitempty"`
	DeviceType           string               `bson:"device_type" json:"device_type"`
	DeviceBrand          string               `bson:"device_brand,omitempty" json:"device_brand,omitempty"`
	SerialNumber         string               `bson:"serial_number,omitempty" json:"serial_number,omitempty"`
	Complaint            string               `bson:"complaint" json:"complaint"`
	Diagnosis            string               `bson:"diagnosis,omitempty" json:"diagnosis,omitempty"`
	EstimatedCost        float64              `bson:"estimated_cost" json:"estimated_cost"`
	ActualCost           float64              `bson:"actual_cost,omitempty" json:"actual_cost,omitempty"`
	EstimatedCompletion  *time.Time           `bson:"estimated_completion,omitempty" json:"estimated_completion,omitempty"`
	Status               ServiceStatus        `bson:"status" json:"status"`
	SparepartsUsed       []ServiceSparepart   `bson:"spareparts_used" json:"spareparts_used"`
	TechnicianNotes      []TechnicianNote     `bson:"technician_notes" json:"technician_notes"`
	Photos               []string             `bson:"photos" json:"photos"`
	TechnicianID         primitive.ObjectID   `bson:"technician_id" json:"technician_id"`
	TechnicianName       string               `bson:"technician_name" json:"technician_name"`
	CreatedAt            time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time            `bson:"updated_at" json:"updated_at"`
}

type ReceivablePayment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Method    PaymentMethod      `bson:"method" json:"method"`
	Date      time.Time          `bson:"date" json:"date"`
	Notes     string             `bson:"notes,omitempty" json:"notes,omitempty"`
}

type ReceivableStatus string

const (
	ReceivableUnpaid   ReceivableStatus = "unpaid"
	ReceivablePartial  ReceivableStatus = "partial"
	ReceivablePaid     ReceivableStatus = "paid"
	ReceivableOverdue  ReceivableStatus = "overdue"
)

type Receivable struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	TransactionID   primitive.ObjectID   `bson:"transaction_id" json:"transaction_id"`
	InvoiceNumber   string               `bson:"invoice_number" json:"invoice_number"`
	CustomerID      primitive.ObjectID   `bson:"customer_id" json:"customer_id"`
	CustomerName    string               `bson:"customer_name" json:"customer_name"`
	TotalAmount     float64              `bson:"total_amount" json:"total_amount"`
	PaidAmount      float64              `bson:"paid_amount" json:"paid_amount"`
	RemainingAmount float64              `bson:"remaining_amount" json:"remaining_amount"`
	DueDate         time.Time            `bson:"due_date" json:"due_date"`
	Status          ReceivableStatus     `bson:"status" json:"status"`
	Payments        []ReceivablePayment  `bson:"payments" json:"payments"`
	CreatedAt       time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updated_at"`
}

type StoreSettings struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Address       string             `bson:"address" json:"address"`
	Phone         string             `bson:"phone" json:"phone"`
	Email         string             `bson:"email" json:"email"`
	Website       string             `bson:"website,omitempty" json:"website,omitempty"`
	Logo          string             `bson:"logo,omitempty" json:"logo,omitempty"`
	TaxRate       float64            `bson:"tax_rate" json:"tax_rate"`
	Currency      string             `bson:"currency" json:"currency"`
	InvoiceFooter string             `bson:"invoice_footer,omitempty" json:"invoice_footer,omitempty"`
	Theme         string             `bson:"theme" json:"theme"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type PaginationQuery struct {
	Page    int    `query:"page"`
	Limit   int    `query:"limit"`
	Search  string `query:"search"`
	SortBy  string `query:"sort_by"`
	SortDir string `query:"sort_dir"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

type DashboardStats struct {
	TodaySales          float64 `json:"today_sales"`
	TodayTransactions   int64   `json:"today_transactions"`
	ActiveServices      int64   `json:"active_services"`
	TotalReceivables    float64 `json:"total_receivables"`
	TotalProducts       int64   `json:"total_products"`
	LowStockProducts    int64   `json:"low_stock_products"`
	BestSellingProduct  ProductInfo `json:"best_selling_product"`
	SlowMovingProduct   ProductInfo `json:"slow_moving_product"`
}

type ProductInfo struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
