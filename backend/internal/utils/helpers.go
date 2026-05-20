package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.New().String()
}

func GenerateInvoiceNumber() string {
	now := time.Now()
	return fmt.Sprintf("INV-%s-%03d", now.Format("20060102"), now.UnixMilli()%1000)
}

func GenerateServiceNumber() string {
	now := time.Now()
	return fmt.Sprintf("SRV-%s-%03d", now.Format("20060102"), now.UnixMilli()%1000)
}

func ParsePaginationQuery(page, limit string) (int, int) {
	p := 1
	l := 20

	if page != "" {
		fmt.Sscanf(page, "%d", &p)
	}
	if limit != "" {
		fmt.Sscanf(limit, "%d", &l)
	}

	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 20
	}

	return p, l
}
