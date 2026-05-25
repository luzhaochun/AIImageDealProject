-- ImageDeal schema MySQL 8
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS tenants (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  slug VARCHAR(50) NOT NULL UNIQUE,
  settings JSON,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS users (
  id CHAR(36) PRIMARY KEY,
  tenant_id CHAR(36) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  display_name VARCHAR(100) DEFAULT '',
  role VARCHAR(20) NOT NULL,
  is_active TINYINT(1) NOT NULL DEFAULT 1,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_tenant_email (tenant_id, email),
  KEY idx_users_tenant (tenant_id),
  CONSTRAINT fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS import_batches (
  id CHAR(36) PRIMARY KEY,
  tenant_id CHAR(36) NOT NULL,
  file_name VARCHAR(255) NOT NULL,
  storage_path VARCHAR(500) DEFAULT '',
  status VARCHAR(30) NOT NULL DEFAULT 'processing',
  total_rows INT NOT NULL DEFAULT 0,
  success_count INT NOT NULL DEFAULT 0,
  failed_count INT NOT NULL DEFAULT 0,
  error_log TEXT,
  created_by CHAR(36) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  KEY idx_batches_tenant (tenant_id),
  CONSTRAINT fk_batches_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS images (
  id CHAR(36) PRIMARY KEY,
  tenant_id CHAR(36) NOT NULL,
  global_no VARCHAR(50) DEFAULT '',
  image_url TEXT,
  storage_path VARCHAR(500) DEFAULT '',
  thumb_path VARCHAR(500) DEFAULT '',
  category VARCHAR(100) DEFAULT '',
  status VARCHAR(30) NOT NULL DEFAULT 'pending_storage',
  assigned_to CHAR(36) DEFAULT NULL,
  import_batch_id CHAR(36) DEFAULT NULL,
  external_id VARCHAR(100) DEFAULT '',
  metadata JSON,
  discard_reason TEXT,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  KEY idx_images_claim (tenant_id, status, category, created_at),
  KEY idx_images_tasks (tenant_id, assigned_to, status),
  KEY idx_images_review (tenant_id, status),
  KEY idx_images_global_no (global_no),
  CONSTRAINT fk_images_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS image_versions (
  id CHAR(36) PRIMARY KEY,
  image_id CHAR(36) NOT NULL,
  version_no INT NOT NULL,
  file_path VARCHAR(500) DEFAULT '',
  layer_data JSON,
  is_current TINYINT(1) NOT NULL DEFAULT 0,
  created_by CHAR(36) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  UNIQUE KEY uk_image_version (image_id, version_no),
  KEY idx_versions_image (image_id),
  CONSTRAINT fk_versions_image FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS review_records (
  id CHAR(36) PRIMARY KEY,
  image_id CHAR(36) NOT NULL,
  version_id CHAR(36) DEFAULT NULL,
  round INT NOT NULL,
  result VARCHAR(20) NOT NULL,
  comment TEXT,
  reviewer_id CHAR(36) NOT NULL,
  created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  KEY idx_review_image (image_id),
  CONSTRAINT fk_review_image FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 种子：租户 + 三角色（密码均为 password123）
INSERT INTO tenants (id, name, slug, settings) VALUES
('00000000-0000-0000-0000-000000000001', 'ACME Demo', 'ACME', '{"second_review_enabled":true,"max_discard_per_day":10}')
ON DUPLICATE KEY UPDATE name=VALUES(name);

INSERT INTO users (id, tenant_id, email, password_hash, display_name, role) VALUES
('00000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000001', 'admin@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '管理员', 'admin'),
('00000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000001', 'user@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '张三', 'user'),
('00000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000001', 'reviewer@example.com', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', '审核员', 'reviewer')
ON DUPLICATE KEY UPDATE email=VALUES(email);
