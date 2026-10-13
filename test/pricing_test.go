package test

import (
	"testing"
)

func TestPricingCalculation(t *testing.T) {
	price := 850000.0
	quantity := 2
	expected := 1700000.0
	actual := price * float64(quantity)

	if actual != expected {
		t.Errorf("Tính toán sai giá tiền: mong đợi %f, thực tế %f", expected, actual)
	}
}

func TestStockDeduction(t *testing.T) {
	currentStock := 15
	bought := 3
	remaining := currentStock - bought

	if remaining != 12 {
		t.Errorf("Trừ tồn kho sai: mong đợi 12, thực tế %d", remaining)
	}
}
