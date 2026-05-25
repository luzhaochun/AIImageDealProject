package services

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"backend/models"
	"backend/pkg/redisx"
	"backend/pkg/storage"

	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
)

type ImportService struct{}

type ParseCSVPayload struct {
	BatchId string `json:"batch_id"`
}

func (s *ImportService) CreateBatch(tenantId, userId, fileName string, body io.Reader) (*models.ImportBatch, error) {
	batchId := uuid.New().String()
	rel := filepath.Join("imports", tenantId, batchId, fileName)
	full := storage.FullPath(rel)
	if err := storage.EnsureDir(full); err != nil {
		return nil, err
	}
	f, err := os.Create(full)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(f, body); err != nil {
		f.Close()
		return nil, err
	}
	f.Close()

	batch := &models.ImportBatch{
		Id:          batchId,
		TenantId:    tenantId,
		FileName:    fileName,
		StoragePath: rel,
		Status:      models.BatchProcessing,
		CreatedBy:   userId,
	}
	o := orm.NewOrm()
	if _, err := o.Insert(batch); err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(ParseCSVPayload{BatchId: batchId})
	_ = redisx.Enqueue(context.Background(), redisx.QueueParseCSV, redisx.TaskMessage{
		Type:     "parse_csv",
		TenantId: tenantId,
		Payload:  payload,
	})
	return batch, nil
}

func (s *ImportService) GetBatch(tenantId, batchId string) (*models.ImportBatch, error) {
	o := orm.NewOrm()
	var b models.ImportBatch
	err := o.QueryTable(new(models.ImportBatch)).Filter("tenant_id", tenantId).Filter("id", batchId).One(&b)
	return &b, err
}

func (s *ImportService) ListBatches(tenantId string, limit, offset int) ([]models.ImportBatch, int64, error) {
	o := orm.NewOrm()
	var list []models.ImportBatch
	qs := o.QueryTable(new(models.ImportBatch)).Filter("tenant_id", tenantId).OrderBy("-created_at")
	total, _ := qs.Count()
	_, err := qs.Limit(limit, offset).All(&list)
	return list, total, err
}

type DownloadPayload struct {
	ImageId string `json:"image_id"`
}

func (s *ImportService) EnqueueDownload(tenantId, imageId string) error {
	payload, _ := json.Marshal(DownloadPayload{ImageId: imageId})
	return redisx.Enqueue(context.Background(), redisx.QueueDownload, redisx.TaskMessage{
		Type:     "download",
		TenantId: tenantId,
		Payload:  payload,
	})
}
