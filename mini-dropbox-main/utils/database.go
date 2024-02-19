package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/manishlpu/assignment/models"
	_ "github.com/go-sql-driver/mysql"
)

type PersistenceDBLayer struct {
	db *sql.DB
	sync.Mutex
}

type MetadataOps interface {
	Exists(id int64) (bool, error)
	SaveRecord(record models.Metadata) (int64, error)
	UpdateRecord(id int64, record models.Metadata) error
	FetchRecords() ([]models.Metadata, error)
	GetRecord(id int64) (*models.Metadata, error)
	DeactivateRecord(id int64) error
	FetchInactiveRecords() ([]models.Metadata, error)
}

func NewPersistenceDBLayer() (MetadataOps, error) {
	database := GetEnvValue("METADATA_DATABASE", "dbname")
	username := GetEnvValue("METADATA_USERNAME", "app-username")
	password := GetEnvValue("METADATA_PASSWORD", "app-password")
	host := GetEnvValue("METADATA_HOST", "dbhost")
	port := GetEnvValue("METADATA_PORT", "3306")

	// Create a DSN (Data Source Name) for the MySQL connection.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)

	// Open a connection to the MySQL database.
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		ErrorLog("could not open database connection: ", err)
		return nil, err
	}

	// Verify the connection by pinging the database.
	if err := db.PingContext(context.TODO()); err != nil {
		ErrorLog("could not ping database: ", err)
		return nil, err
	}

	return &PersistenceDBLayer{
		db: db,
	}, nil
}

func (pdb *PersistenceDBLayer) Exists(id int64) (bool, error) {
	// Query to check if a record with the given ID exists
	query := "SELECT 1 FROM file_metadata WHERE id = ? AND status = 1 LIMIT 1"

	// Execute the query with the target ID
	var exists bool
	err := pdb.db.QueryRow(query, id).Scan(&exists)

	// Check for errors
	if err == sql.ErrNoRows {
		DebugLog("No record found with ID ", id)
	} else if err != nil {
		return false, err
	}
	return exists, nil
}

// Insert a new metadata record into the database
func (pdb *PersistenceDBLayer) SaveRecord(record models.Metadata) (int64, error) {
	_, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Insert new metadata into the "file_metadata" table.
	stmt, err := pdb.db.Prepare("INSERT INTO file_metadata (filename, size_in_bytes, s3_object_key, description, mime_type, status) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return int64(-1), err
	}
	defer stmt.Close()

	pdb.Lock()
	defer pdb.Unlock()
	// Execute the SQL statement to insert the new row
	res, err := stmt.Exec(record.Filename, record.SizeInBytes, record.S3ObjectKey, record.Description, record.MimeType, record.Status)
	if err != nil {
		return int64(-1), err
	}

	return res.LastInsertId()
}

// Update an existing metadata row in the database.
func (pdb *PersistenceDBLayer) UpdateRecord(id int64, record models.Metadata) error {
	// Replace with your update statement
	updateSQL := "UPDATE file_metadata SET filename = ?, size_in_bytes = ?, s3_object_key = ?, mime_type = ?, description = ? WHERE id = ? AND status = 1"

	// Execute the update statement
	result, err := pdb.db.Exec(updateSQL, record.Filename, record.SizeInBytes, record.S3ObjectKey, record.MimeType, record.Description, id)
	if err != nil {
		return err
	}

	// Check the number of rows affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		InfoLog("Row updated successfully with id: ", id)
	} else {
		WarnLog("No rows were updated for ID ", id)
	}
	return nil
}

// Returns all the active metadata records from Database.
func (pdb *PersistenceDBLayer) FetchRecords() ([]models.Metadata, error) {
	// Query to retrieve records with "filename" and "description" fields.
	query := "SELECT id, filename, size_in_bytes, s3_object_key, description, mime_type, created_at, updated_at FROM file_metadata WHERE status = 1"

	// Execute the query and retrieve the results.
	rows, err := pdb.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate through the rows and store results in a slice of File structs.
	var files []models.Metadata
	for rows.Next() {
		var file models.Metadata
		if err := rows.Scan(&file.ID, &file.Filename, &file.SizeInBytes, &file.S3ObjectKey, &file.Description, &file.MimeType, &file.CreatedAt, &file.UpdatedAt); err != nil {
			ErrorLog("unable to get file metadata")
			continue
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (pdb *PersistenceDBLayer) GetRecord(id int64) (*models.Metadata, error) {
	// Query to fetch the metadata associated with the given identifier.
	query := "SELECT id, filename, size_in_bytes, s3_object_key, description, mime_type, created_at, updated_at FROM file_metadata WHERE id = ? AND status = 1"

	// Execute the query with the primary key value
	var metadata models.Metadata
	err := pdb.db.QueryRow(query, id).Scan(
		&metadata.ID, &metadata.Filename, &metadata.SizeInBytes, &metadata.S3ObjectKey,
		&metadata.Description, &metadata.MimeType, &metadata.CreatedAt, &metadata.UpdatedAt,
	)

	// Check for errors
	if err == sql.ErrNoRows {
		// No record found with the given identifier, not an error.
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func (pdb *PersistenceDBLayer) DeactivateRecord(id int64) error {
	query := "UPDATE file_metadata SET status = 0 WHERE id =?"
	pdb.Lock()
	defer pdb.Unlock()

	res, err := pdb.db.Exec(query, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no rows affected")
	}
	return err
}

func (pdb *PersistenceDBLayer) FetchInactiveRecords() ([]models.Metadata, error) {
	// Calculate the date 30 days ago
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02 15:04:05")

	// Query to delete rows with status = 1 and updatedAt < 30 days ago
	query := "SELECT id, filename, s3_object_key, status, created_at updated_at FROM file_metadata WHERE status = 0 AND updatedAt < " + thirtyDaysAgo

	// Execute the query and retrieve the results.
	rows, err := pdb.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and store results in a slice of File structs.
	var files []models.Metadata
	for rows.Next() {
		var file models.Metadata
		if err := rows.Scan(&file.ID, &file.Filename, &file.S3ObjectKey, &file.Status, &file.CreatedAt, &file.UpdatedAt); err != nil {
			ErrorLog("unable to get file metadata")
			continue
		}

		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}
