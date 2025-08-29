package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"


	"github.com/lib/pq"  // postgreSQL driver
)

// DB wraps our database connection
type DB strcut {
	*sql.DB
}

//NewConnection creates a connection pool
//This is THE MOST IMPORTANT pattern for production

func NewConnection(databseURL string) (*DB, error) {

	//sql.Open doesn't actually connect, just prepares
	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		return nil, fmt.Errorf("error opening databse: %w", err)
	}



	// CRITICAL: Configure connection pool
	// These settings prevent databse overload

	db.SetMaxOpenConns(25)				   // Max 25 connections
	db.SetMaxIdleConns(5)				   // Keep 5 ready
	db.SetConnMaxLifetime(5 * time.Minute) // Refresh connection

	// Actually test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return &DB{db}, nil
}


// SaveVerification stores our verification result
// Notice: We're using the DB receiver (d *DB)

func (d *DB) SaveVErification(v VErificationsRecord) error {
	query := `
		INSERT INTO verifications
		(id, tenant_name, tenant_email, income, risk_score, status, verified_at, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	//Execute query with params (prevents SQL injection)
	_, err := d.Exec(
		query,
		v.ID,
		v.TenantName,
		v.TenantEmail,
		v.Income,
		v.RiskScore,
		v.Status,
		v.VerifiedAt,
		v.Details,
	)

	if err != nil {
		return fmt.Errorf("error saving verification: %w", err)
	}

	log.Printf("💾 Saved verifications: %s", v.ID)
	return nil

}

// VerificationRecord matches our database schema
type VerificationRecord struct {
	ID string
	TenantName string
	TenantEmail string
	Income string
	RiskScore int
	Status string
	VerifiedAt time.Time
	Details []string
}