package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB
var StartTime = time.Now()

func InitDB() (*sql.DB, error) {
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "root")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "archival_vault")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err := DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	if err := migrate(); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	if err := seed(); err != nil {
		log.Printf("[WARN] seed error: %v", err)
	}

	log.Printf("[DB: OK] Connected to MySQL (%s:%s/%s)", dbHost, dbPort, dbName)
	return DB, nil
}

func migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'ARCHIVIST',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		`CREATE TABLE IF NOT EXISTS records (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			code VARCHAR(50) NOT NULL UNIQUE,
			title VARCHAR(255) NOT NULL,
			classification VARCHAR(20) NOT NULL DEFAULT 'INTERNAL',
			category VARCHAR(50) NOT NULL DEFAULT 'SPEC',
			status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
			tags VARCHAR(255) NOT NULL DEFAULT '',
			summary TEXT,
			payload LONGTEXT,
			created_by VARCHAR(50) NOT NULL DEFAULT 'SYSTEM',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_category (category),
			INDEX idx_classification (classification),
			INDEX idx_status (status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func seed() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = DB.Exec(
			"INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)",
			"archivist", "archivist@raw-vault.local", string(hash), "ARCHIVIST",
		)
		if err != nil {
			return err
		}
		log.Println("[SEED] Created default user: archivist (password: password123)")
	}

	var recCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM records").Scan(&recCount)
	if err != nil {
		return err
	}

	if recCount == 0 {
		seeds := []struct {
			code, title, class, cat, status, tags, summary, payload, by string
		}{
			{
				code:    "ARC-2026-001",
				title:   "Distributed Kernel Synchronization Protocol",
				class:   "RESTRICTED",
				cat:     "SPEC",
				status:  "ACTIVE",
				tags:    "golang,raft,consensus,v2",
				summary: "Architectural blueprint detailing zero-overhead state replication over raw TCP framing.",
				payload: "{\n  \"protocol_version\": \"2.4.0-raw\",\n  \"frame_size\": 4096,\n  \"heartbeat_ms\": 50,\n  \"nodes\": [\"node-alpha\", \"node-bravo\", \"node-charlie\"],\n  \"encryption\": \"ChaCha20-Poly1305\"\n}",
				by:      "archivist",
			},
			{
				code:    "ARC-2026-002",
				title:   "Core Mesh Socket Leak Postmortem",
				class:   "INTERNAL",
				cat:     "INCIDENT",
				status:  "ARCHIVED",
				tags:    "network,leak,epoll,resolved",
				summary: "Root cause analysis of file descriptor starvation during high burst ingress at edge proxy.",
				payload: "INCIDENT_ID: INC-7719\nIMPACT: 2.1% ingress drops for 142s\nCAUSE: Missing defer conn.Close() on TLS handshake failure path.\nRESOLUTION: Patched in commit a9f4c1; file descriptor budget expanded to 65536.",
				by:      "archivist",
			},
			{
				code:    "ARC-2026-003",
				title:   "Identity Federation & JWT Claims Schema",
				class:   "CONFIDENTIAL",
				cat:     "SCHEMA",
				status:  "ACTIVE",
				tags:    "jwt,auth,claims,security",
				summary: "Standardized claim signatures and asymmetric token verification criteria for fullstack clients.",
				payload: "{\n  \"issuer\": \"archival.vault.gateway\",\n  \"audience\": [\"vue-client\", \"cli-agent\"],\n  \"algorithm\": \"HS256\",\n  \"ttl_hours\": 24,\n  \"required_claims\": [\"uid\", \"usr\", \"rol\", \"exp\"]\n}",
				by:      "archivist",
			},
			{
				code:    "ARC-2026-004",
				title:   "Telemetry Wireframe Grid Layout Matrix",
				class:   "PUBLIC",
				cat:     "DISPATCH",
				status:  "ACTIVE",
				tags:    "ui,wireframe,monospace,design",
				summary: "Design specification for zero-radius high contrast monospaced archival interfaces.",
				payload: "FONT_STACK: 'Courier New', 'Consolas', monospace\nBORDER_SPEC: 1px solid var(--border-wire)\nPALETTE: { bg: '#0c0d0e', surface: '#14171a', border: '#2b313a', text: '#e6edf3' }\nACCENT_CLASS: { CONFIDENTIAL: '#ff5370', RESTRICTED: '#ffcb6b', INTERNAL: '#82aaff', PUBLIC: '#c3e88d' }",
				by:      "archivist",
			},
			{
				code:    "ARC-2026-005",
				title:   "Quantum Resistant Key Exchange Proposal",
				class:   "CONFIDENTIAL",
				cat:     "SPEC",
				status:  "DRAFT",
				tags:    "crypto,mlkem,post-quantum",
				summary: "Exploratory draft integrating ML-KEM-768 key encapsulation into transport layer.",
				payload: "STAGE: DRAFT-01\nPREREQUISITE: NIST FIPS 203 conformance\nTESTBED: Local test cluster 127.0.0.1:8080\nSCHEDULED_AUDIT: Q4 2026",
				by:      "archivist",
			},
		}

		for _, s := range seeds {
			_, err := DB.Exec(
				`INSERT INTO records (code, title, classification, category, status, tags, summary, payload, created_by)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				s.code, s.title, s.class, s.cat, s.status, s.tags, s.summary, s.payload, s.by,
			)
			if err != nil {
				return err
			}
		}
		log.Println("[SEED] Inserted 5 initial archival wireframe records")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
