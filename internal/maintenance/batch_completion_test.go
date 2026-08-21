package maintenance

import (
	"context"
	"reflect"
	"testing"
	"time"
)

func completionFixture(id string, now time.Time) (WorkOrder, Asset, ServiceReport) {
	started := now.Add(-time.Hour)
	order := WorkOrder{ID: "order-" + id, TenantID: "farm", AssetID: "asset-" + id, AssignedTo: "tech", Status: "in_progress", MeterAtOpen: 100, StartedAt: &started, Version: 4}
	asset := Asset{ID: "asset-" + id, TenantID: "farm", BarnID: "barn", Name: id, Category: "milking", Status: AssetServicing, MeterHours: 100, ServiceEvery: 200, LastServiceHour: 80, Version: 7}
	report := ServiceReport{WorkOrderID: order.ID, AssetID: asset.ID, Technician: "tech", LaborMins: 40, Resolution: "replaced worn coupling", MeterHours: 101, CompletedAt: now}
	return order, asset, report
}

func TestCompleteBatchRejectsPartialMaintenanceState(t *testing.T) {
	now := time.Now().UTC()
	firstOrder, firstAsset, firstReport := completionFixture("first", now)
	secondOrder, secondAsset, secondReport := completionFixture("second", now)
	secondReport.Resolution = ""
	orders := []WorkOrder{firstOrder, secondOrder}
	assets := []Asset{firstAsset, secondAsset}
	wantOrders := append([]WorkOrder(nil), orders...)
	wantAssets := append([]Asset(nil), assets...)
	requests := []CompletionRequest{
		{OrderIndex: 0, AssetIndex: 0, Report: firstReport, ExpectedOrderVersion: 4, ExpectedAssetVersion: 7},
		{OrderIndex: 1, AssetIndex: 1, Report: secondReport, ExpectedOrderVersion: 4, ExpectedAssetVersion: 7},
	}
	if err := CompleteBatch(context.Background(), orders, assets, requests); err == nil {
		t.Fatal("invalid second service report was accepted")
	}
	if !reflect.DeepEqual(orders, wantOrders) || !reflect.DeepEqual(assets, wantAssets) {
		t.Fatalf("failed batch exposed partial state: orders=%+v assets=%+v", orders, assets)
	}
}
