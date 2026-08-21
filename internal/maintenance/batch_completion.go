package maintenance

import (
	"context"
	"fmt"

	"go-base/internal/domain"
)

func CompleteBatch(ctx context.Context, orders []WorkOrder, assets []Asset, requests []CompletionRequest) error {
	if len(requests) == 0 {
		return fmt.Errorf("%w: completion batch", domain.ErrInvalid)
	}

	// Phase 1: validate and evolve state against working copies. A failure in
	// any request returns early, leaving the caller's orders and assets
	// untouched so the batch stays all-or-nothing across validation.
	workOrders := append([]WorkOrder(nil), orders...)
	workAssets := append([]Asset(nil), assets...)
	for _, request := range requests {
		if err := ctx.Err(); err != nil {
			return err
		}
		if request.OrderIndex < 0 || request.OrderIndex >= len(workOrders) || request.AssetIndex < 0 || request.AssetIndex >= len(workAssets) {
			return fmt.Errorf("%w: completion batch index", domain.ErrInvalid)
		}
		completedOrder, servicedAsset, err := Complete(
			workOrders[request.OrderIndex],
			workAssets[request.AssetIndex],
			request.Report,
			request.ExpectedOrderVersion,
			request.ExpectedAssetVersion,
		)
		if err != nil {
			return err
		}
		workOrders[request.OrderIndex] = completedOrder
		workAssets[request.AssetIndex] = servicedAsset
	}

	// Phase 2: every request validated and evolved successfully, so publish the
	// whole batch back to the caller in one shot.
	copy(orders, workOrders)
	copy(assets, workAssets)
	return nil
}
