package services

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"backend/models"
	"backend/pkg/storage"

	"github.com/beego/beego/v2/client/orm"
	"github.com/google/uuid"
)

type ExportRequest struct {
	Category string `json:"category"`
	Status   string `json:"status"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type ExportJob struct {
	Id           string    `json:"id"`
	Status       string    `json:"status"`
	FilterLabel  string    `json:"filter_label"`
	FileName     string    `json:"file_name"`
	DownloadPath string    `json:"download_path,omitempty"`
	ImageCount   int       `json:"image_count"`
	CreatedAt    time.Time `json:"created_at"`
}

func exportMetaPath(tenantId, exportId string) string {
	return filepath.Join("exports", tenantId, exportId+".json")
}

func (s *ArchiveService) ListExports(tenantId string) ([]ExportJob, error) {
	dir := storage.FullPath(filepath.Join("exports", tenantId))
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []ExportJob{}, nil
		}
		return nil, err
	}
	out := make([]ExportJob, 0)
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var job ExportJob
		if json.Unmarshal(b, &job) == nil {
			out = append(out, job)
		}
	}
	return out, nil
}

func (s *ArchiveService) CreateExport(tenantId string, req ExportRequest) (*ExportJob, error) {
	status := req.Status
	if status == "" {
		status = models.StatusCompleted
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(models.Image)).Filter("tenant_id", tenantId).Filter("status", status)
	if req.Category != "" {
		qs = qs.Filter("category", req.Category)
	}
	if req.DateFrom != "" {
		qs = qs.Filter("updated_at__gte", req.DateFrom)
	}
	if req.DateTo != "" {
		qs = qs.Filter("updated_at__lte", req.DateTo+" 23:59:59")
	}
	var images []models.Image
	if _, err := qs.OrderBy("-updated_at").Limit(500, 0).All(&images); err != nil {
		return nil, err
	}

	exportId := uuid.New().String()
	zipRel := filepath.Join("exports", tenantId, exportId+".zip")
	zipFull := storage.FullPath(zipRel)
	if err := storage.EnsureDir(zipFull); err != nil {
		return nil, err
	}

	zf, err := os.Create(zipFull)
	if err != nil {
		return nil, err
	}
	zw := zip.NewWriter(zf)
	manifest, _ := zw.Create("manifest.csv")
	cw := csv.NewWriter(manifest)
	_ = cw.Write([]string{"global_no", "category", "status", "storage_path"})

	for _, img := range images {
		_ = cw.Write([]string{img.GlobalNo, img.Category, img.Status, img.StoragePath})
		if img.StoragePath != "" {
			src := storage.FullPath(img.StoragePath)
			if st, err := os.Stat(src); err == nil && !st.IsDir() {
				f, err := os.Open(src)
				if err == nil {
					w, _ := zw.Create(filepath.Base(img.StoragePath))
					_, _ = io.Copy(w, f)
					f.Close()
				}
			}
		}
	}
	cw.Flush()
	if err := zw.Close(); err != nil {
		zf.Close()
		return nil, err
	}
	zf.Close()

	label := req.Category
	if label == "" {
		label = "全部"
	}
	job := &ExportJob{
		Id:           exportId,
		Status:       "completed",
		FilterLabel:  fmt.Sprintf("%s · %s", label, status),
		FileName:     exportId + ".zip",
		DownloadPath: zipRel,
		ImageCount:   len(images),
		CreatedAt:    time.Now(),
	}
	metaFull := storage.FullPath(exportMetaPath(tenantId, exportId))
	_ = storage.EnsureDir(metaFull)
	b, _ := json.Marshal(job)
	_ = os.WriteFile(metaFull, b, 0644)
	return job, nil
}
