package orders

import "fmt"

func validateCreateOrder(params createOrderParams) error {
	if params.CustomerID <= 0 {
		return fmt.Errorf("customer ID must be positive")
	}

	if len(params.Items) == 0 {
		return fmt.Errorf("at least one item is required")
	}

	for i, item := range params.Items {
		if item.ProductID <= 0 {
			return fmt.Errorf("item %d: product ID must be positive", i+1)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item %d: quantity must be positive", i+1)
		}
	}

	return nil
}
