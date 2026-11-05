package orders

import "testing"

func TestValidateCreateOrder(t *testing.T) {
	tests := []struct {
		name    string
		params  createOrderParams
		wantErr string
	}{
		{
			name: "valid order",
			params: createOrderParams{
				CustomerID: 1,
				Items: []orderItem{
					{ProductID: 1, Quantity: 2},
				},
			},
		},
		{
			name: "valid order with multiple items",
			params: createOrderParams{
				CustomerID: 1,
				Items: []orderItem{
					{ProductID: 1, Quantity: 1},
					{ProductID: 2, Quantity: 2},
				},
			},
		},
		{
			name: "zero customer ID",
			params: createOrderParams{
				CustomerID: 0,
				Items: []orderItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			wantErr: "customer ID must be positive",
		},
		{
			name: "negative customer ID",
			params: createOrderParams{
				CustomerID: -1,
				Items: []orderItem{
					{ProductID: 1, Quantity: 2},
				},
			},
			wantErr: "customer ID must be positive",
		},
		{
			name: "no items",
			params: createOrderParams{
				CustomerID: 1,
			},
			wantErr: "at least one item is required",
		},
		{
			name: "zero product ID",
			params: createOrderParams{
				CustomerID: 1,
				Items: []orderItem{
					{ProductID: 0, Quantity: 2},
				},
			},
			wantErr: "item 1: product ID must be positive",
		},
		{
			name: "negative product ID",
			params: createOrderParams{
				CustomerID: 1,
				Items: []orderItem{
					{ProductID: -1, Quantity: 2},
				},
			},
			wantErr: "item 1: product ID must be positive",
		},
		{
			name: "zero quantity",
			params: createOrderParams{
				CustomerID: 1,
				Items: []orderItem{
					{ProductID: 1, Quantity: 0},
				},
			},
			wantErr: "item 1: quantity must be positive",
		},
		{
			name: "invalid second item",
			params: createOrderParams{
				CustomerID: 1,
				Items: []orderItem{
					{ProductID: 1, Quantity: 2},
					{ProductID: 2, Quantity: -1},
				},
			},
			wantErr: "item 2: quantity must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateOrder(tt.params)

			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}

			if err.Error() != tt.wantErr {
				t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}
