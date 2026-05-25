package models

import (
	"time"
)

type Tenant struct {
	Id        string    `orm:"column(id);pk" json:"id"`
	Name      string    `orm:"column(name)" json:"name"`
	Slug      string    `orm:"column(slug)" json:"slug"`
	Settings  string    `orm:"column(settings);type(json);null" json:"settings,omitempty"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (t *Tenant) TableName() string { return "tenants" }

type User struct {
	Id           string    `orm:"column(id);pk" json:"id"`
	TenantId     string    `orm:"column(tenant_id)" json:"tenant_id"`
	Email        string    `orm:"column(email)" json:"email"`
	PasswordHash string    `orm:"column(password_hash)" json:"-"`
	DisplayName  string    `orm:"column(display_name)" json:"display_name"`
	Role         string    `orm:"column(role)" json:"role"`
	IsActive     int8      `orm:"column(is_active)" json:"is_active"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (u *User) TableName() string { return "users" }

type ImportBatch struct {
	Id           string    `orm:"column(id);pk" json:"id"`
	TenantId     string    `orm:"column(tenant_id)" json:"tenant_id"`
	FileName     string    `orm:"column(file_name)" json:"file_name"`
	StoragePath  string    `orm:"column(storage_path)" json:"storage_path,omitempty"`
	Status       string    `orm:"column(status)" json:"status"`
	TotalRows    int       `orm:"column(total_rows)" json:"total_rows"`
	SuccessCount int       `orm:"column(success_count)" json:"success_count"`
	FailedCount  int       `orm:"column(failed_count)" json:"failed_count"`
	ErrorLog     string    `orm:"column(error_log);null" json:"error_log,omitempty"`
	CreatedBy    string    `orm:"column(created_by)" json:"created_by"`
	CreatedAt    time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt    time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (b *ImportBatch) TableName() string { return "import_batches" }

type Image struct {
	Id            string    `orm:"column(id);pk" json:"id"`
	TenantId      string    `orm:"column(tenant_id)" json:"tenant_id"`
	GlobalNo      string    `orm:"column(global_no)" json:"global_no,omitempty"`
	ImageUrl      string    `orm:"column(image_url);null" json:"image_url,omitempty"`
	StoragePath   string    `orm:"column(storage_path)" json:"storage_path,omitempty"`
	ThumbPath     string    `orm:"column(thumb_path)" json:"thumb_path,omitempty"`
	Category      string    `orm:"column(category)" json:"category,omitempty"`
	Status        string    `orm:"column(status)" json:"status"`
	AssignedTo    string    `orm:"column(assigned_to);null" json:"assigned_to,omitempty"`
	ImportBatchId string    `orm:"column(import_batch_id);null" json:"import_batch_id,omitempty"`
	ExternalId    string    `orm:"column(external_id)" json:"external_id,omitempty"`
	Metadata      string    `orm:"column(metadata);type(json);null" json:"metadata,omitempty"`
	DiscardReason string    `orm:"column(discard_reason);null" json:"discard_reason,omitempty"`
	CreatedAt     time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
	UpdatedAt     time.Time `orm:"column(updated_at);auto_now;type(datetime)" json:"updated_at"`
}

func (i *Image) TableName() string { return "images" }

type ImageVersion struct {
	Id        string    `orm:"column(id);pk" json:"id"`
	ImageId   string    `orm:"column(image_id)" json:"image_id"`
	VersionNo int       `orm:"column(version_no)" json:"version_no"`
	FilePath  string    `orm:"column(file_path)" json:"file_path,omitempty"`
	LayerData string    `orm:"column(layer_data);type(json);null" json:"layer_data,omitempty"`
	IsCurrent int8      `orm:"column(is_current)" json:"is_current"`
	CreatedBy string    `orm:"column(created_by)" json:"created_by"`
	CreatedAt time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
}

func (v *ImageVersion) TableName() string { return "image_versions" }

type ReviewRecord struct {
	Id         string    `orm:"column(id);pk" json:"id"`
	ImageId    string    `orm:"column(image_id)" json:"image_id"`
	VersionId  string    `orm:"column(version_id);null" json:"version_id,omitempty"`
	Round      int       `orm:"column(round)" json:"round"`
	Result     string    `orm:"column(result)" json:"result"`
	Comment    string    `orm:"column(comment);null" json:"comment,omitempty"`
	ReviewerId string    `orm:"column(reviewer_id)" json:"reviewer_id"`
	CreatedAt  time.Time `orm:"column(created_at);auto_now_add;type(datetime)" json:"created_at"`
}

func (r *ReviewRecord) TableName() string { return "review_records" }
