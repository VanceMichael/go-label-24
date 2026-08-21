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
	for _, request := range requests {
		if err := ctx.Err(); err != nil {
			return err
		}
		if request.OrderIndex < 0 || request.OrderIndex >= len(orders) || request.AssetIndex < 0 || request.AssetIndex >= len(assets) {
			return fmt.Errorf("%w: completion batch index", domain.ErrInvalid)
		}
		completedOrder, servicedAsset, err := Complete(
			orders[request.OrderIndex],
			assets[request.AssetIndex],
			request.Report,
			request.ExpectedOrderVersion,
			request.ExpectedAssetVersion,
		)
		if err != nil {
			return err
		}
		orders[request.OrderIndex] = completedOrder
		assets[request.AssetIndex] = servicedAsset
	}
	return nil
}
