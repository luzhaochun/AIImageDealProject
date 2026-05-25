package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/models"
	"backend/pkg/storage"
	"backend/services"

	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
)

func handleParseCSV(ctx context.Context, tenantId string, payload json.RawMessage) error {
	var p services.ParseCSVPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	o := orm.NewOrm()
	var batch models.ImportBatch
	if err := o.QueryTable(new(models.ImportBatch)).Filter("id", p.BatchId).One(&batch); err != nil {
		return err
	}
	f, err := os.Open(storage.FullPath(batch.StoragePath))
	if err != nil {
		batch.Status = models.BatchFailed
		batch.ErrorLog = err.Error()
		_, _ = o.Update(&batch, "Status", "ErrorLog")
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true
	header, err := r.Read()
	if err != nil {
		return err
	}
	colURL, colCat := -1, -1
	for i, h := range header {
		h = strings.ToLower(strings.TrimSpace(h))
		if h == "image_url" || h == "url" {
			colURL = i
		}
		if h == "category" {
			colCat = i
		}
	}
	if colURL < 0 || colCat < 0 {
		batch.Status = models.BatchFailed
		batch.ErrorLog = "CSV 需要列: image_url, category"
		_, _ = o.Update(&batch, "Status", "ErrorLog")
		return nil
	}

	importSvc := &services.ImportService{}
	total, success, failed := 0, 0, 0
	var errLines []string

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errLines = append(errLines, err.Error())
			failed++
			continue
		}
		total++
		if colURL >= len(row) || colCat >= len(row) {
			failed++
			continue
		}
		url := strings.TrimSpace(row[colURL])
		cat := strings.TrimSpace(row[colCat])
		if url == "" || cat == "" {
			failed++
			continue
		}
		img := &models.Image{
			Id:            uuid.New().String(),
			TenantId:      batch.TenantId,
			ImageUrl:      url,
			Category:      cat,
			Status:        models.StatusPendingStorage,
			ImportBatchId: batch.Id,
		}
		if _, err := o.Insert(img); err != nil {
			failed++
			errLines = append(errLines, err.Error())
			continue
		}
		success++
		_ = importSvc.EnqueueDownload(batch.TenantId, img.Id)
	}

	batch.TotalRows = total
	batch.SuccessCount = success
	batch.FailedCount = failed
	if len(errLines) > 0 {
		batch.ErrorLog = strings.Join(errLines[:min(20, len(errLines))], "\n")
	}
	if success > 0 {
		batch.Status = models.BatchProcessing
	} else {
		batch.Status = models.BatchFailed
	}
	_, _ = o.Update(&batch, "TotalRows", "SuccessCount", "FailedCount", "ErrorLog", "Status")
	return nil
}

func handleDownload(ctx context.Context, tenantId string, payload json.RawMessage) error {
	var p services.DownloadPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}
	o := orm.NewOrm()
	var img models.Image
	if err := o.QueryTable(new(models.Image)).Filter("id", p.ImageId).One(&img); err != nil {
		return err
	}
	if img.Status != models.StatusPendingStorage {
		return nil
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(img.ImageUrl)
	if err != nil {
		markDownloadFail(o, &img)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		markDownloadFail(o, &img)
		return nil
	}

	ext := ".jpg"
	rel := filepath.Join("originals", img.TenantId, img.Id+ext)
	_, err = storage.SaveReader(rel, resp.Body)
	if err != nil {
		markDownloadFail(o, &img)
		return err
	}
	thumbRel := filepath.Join("thumbs", img.TenantId, img.Id+ext)
	// MVP: thumb = copy path reference same file
	img.StoragePath = rel
	img.ThumbPath = thumbRel
	img.Status = models.StatusPendingAssign
	_, err = o.Update(&img, "StoragePath", "ThumbPath", "Status")
	if err != nil {
		return err
	}

	// update batch progress
	if img.ImportBatchId != "" {
		var batch models.ImportBatch
		if err := o.QueryTable(new(models.ImportBatch)).Filter("id", img.ImportBatchId).One(&batch); err == nil {
			done, _ := o.QueryTable(new(models.Image)).Filter("import_batch_id", batch.Id).Filter("status", models.StatusPendingAssign).Count()
			pending, _ := o.QueryTable(new(models.Image)).Filter("import_batch_id", batch.Id).Filter("status", models.StatusPendingStorage).Count()
			batch.SuccessCount = int(done)
			if pending == 0 && batch.Status == models.BatchProcessing {
				batch.Status = models.BatchCompleted
			}
			_, _ = o.Update(&batch, "SuccessCount", "Status")
		}
	}
	return nil
}

func markDownloadFail(o orm.Ormer, img *models.Image) {
	var batch models.ImportBatch
	if img.ImportBatchId != "" {
		_ = o.QueryTable(new(models.ImportBatch)).Filter("id", img.ImportBatchId).One(&batch)
		batch.FailedCount++
		_, _ = o.Update(&batch, "FailedCount")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
