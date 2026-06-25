package service

import (
	"context"
	"fmt"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/repository"
)

type CartService interface {
	Quote(ctx context.Context, req dto.CartQuoteRequest) (*dto.CartQuoteResponse, error)
}

type cartService struct {
	products repository.ProductRepository
}

func NewCartService(products repository.ProductRepository) CartService {
	return &cartService{products: products}
}

func (s *cartService) Quote(ctx context.Context, req dto.CartQuoteRequest) (*dto.CartQuoteResponse, error) {
	response := &dto.CartQuoteResponse{
		Items:       make([]dto.CartQuoteItemResponse, 0, len(req.Items)),
		Warnings:    []string{},
		ShippingFee: 0,
	}

	for _, item := range req.Items {
		product, err := s.products.FindByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		quoteItem := dto.CartQuoteItemResponse{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}

		if product == nil {
			quoteItem.IsAvailable = false
			response.Warnings = append(response.Warnings, fmt.Sprintf("Product %d no longer exists", item.ProductID))
			response.Items = append(response.Items, quoteItem)
			continue
		}

		quoteItem.Name = product.Name
		quoteItem.ImageURL = product.ImageURL
		quoteItem.UnitPrice = product.Price
		quoteItem.AvailableStock = product.Stock

		isAvailable := product.IsActive && product.Category != nil && product.Category.IsActive && product.Stock >= item.Quantity
		quoteItem.IsAvailable = isAvailable

		if !product.IsActive || product.Category == nil || !product.Category.IsActive {
			response.Warnings = append(response.Warnings, fmt.Sprintf("Product %s is inactive", product.Name))
		}
		if product.Stock < item.Quantity {
			response.Warnings = append(response.Warnings, fmt.Sprintf("Product %s only has %d item(s) left", product.Name, product.Stock))
		}

		if isAvailable {
			quoteItem.Subtotal = product.Price * int64(item.Quantity)
			response.Subtotal += quoteItem.Subtotal
		}

		response.Items = append(response.Items, quoteItem)
	}

	response.DiscountAmount = 0
	response.FinalAmount = response.Subtotal + response.ShippingFee - response.DiscountAmount
	return response, nil
}
