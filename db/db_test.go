package db

import (
	"os"
	"testing"
)

func TestInsertAndGetFile(t *testing.T) {
	// Setup temporary database
	tmpDB := "test.db"
	os.Setenv("DB_PATH", tmpDB)

	// Initialize DB
	// We need to suppress logs or accept them during tests.
	// InitDB calls openDB which reads env.DBPath.

	// Since InitDB calls os.Exit on failure, we need to be careful.
	// But it shouldn't fail if we give it a writable path.

	// We might need to ensure the directory exists if it's in a subfolder,
	// but "test.db" is in current dir.

	// Reset global db variable if possible, but it's private.
	// However, InitDB overwrites it.

	// Silence logger for cleaner test output (optional)
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))

	InitDB()
	defer os.Remove(tmpDB)
	defer db.Close() // Assuming we can close it, but db is private and *sql.DB.

	f := &File{
		PubID: "test_pub_id",
		Meta: FileMeta{
			OriginalName:  "test.txt",
			Size:          10,
			Hash:          "test_hash",
			LocalFileName: "test.txt",
			FileType:      "text/plain",
		},
		Issuer: "tester",
	}

	err := InsertFile(f)
	if err != nil {
		t.Fatalf("Failed to insert file: %v", err)
	}

	if f.ID == 0 {
		t.Errorf("Expected ID to be populated, got 0")
	}

	retrieved, err := GetFileByPubID("test_pub_id")
	if err != nil {
		t.Fatalf("Failed to get file: %v", err)
	}

	if retrieved == nil {
		t.Fatalf("Expected file to be found, got nil")
	}

	if retrieved.PubID != f.PubID {
		t.Errorf("Expected PubID %s, got %s", f.PubID, retrieved.PubID)
	}
	if retrieved.Meta.Hash != f.Meta.Hash {
		t.Errorf("Expected Meta %s, got %s", f.Meta.Hash, retrieved.Meta.Hash)
	}
	if retrieved.Issuer != f.Issuer {
		t.Errorf("Expected Issuer %s, got %s", f.Issuer, retrieved.Issuer)
	}
}

func TestInsertAndGetPost(t *testing.T) {
	// Setup temporary database
	tmpDB := "test.db"
	os.Setenv("DB_PATH", tmpDB)

	// Initialize DB
	// We need to suppress logs or accept them during tests.
	// InitDB calls openDB which reads env.DBPath.

	// Since InitDB calls os.Exit on failure, we need to be careful.
	// But it shouldn't fail if we give it a writable path.

	// We might need to ensure the directory exists if it's in a subfolder,
	// but "test.db" is in current dir.

	// Reset global db variable if possible, but it's private.
	// However, InitDB overwrites it.

	// Silence logger for cleaner test output (optional)
	// slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError})))

	InitDB()
	defer os.Remove(tmpDB)
	defer db.Close() // Assuming we can close it, but db is private and *sql.DB.

	// Post structure is:
	// pubID
	// content
	// issuer
	// Insert post
	p := &Post{
		PubID:   "test_pub_id",
		Content: `{"files": ["12", "34"],"description": "description"}`,
		Issuer:  "tester",
	}

	err := InsertPost(p)
	if err != nil {
		t.Fatalf("Failed to insert post: %v", err)
	}

	retrieved, err := GetPostByPubID("test_pub_id")
	if err != nil {
		t.Fatalf("Failed to get post: %v", err)
	}

	if retrieved == nil {
		t.Fatalf("Expected post to be found, got nil")
	}

	if retrieved.ID == 0 {
		t.Errorf("Expected ID to be populated, got 0")
	}

	if retrieved.PubID != p.PubID {
		t.Errorf("Expected PubID %s, got %s", p.PubID, retrieved.PubID)
	}
	if retrieved.Content != p.Content {
		t.Errorf("Expected Content %s, got %s", p.Content, retrieved.Content)
	}
	if retrieved.Issuer != p.Issuer {
		t.Errorf("Expected Issuer %s, got %s", p.Issuer, retrieved.Issuer)
	}
}

func TestCounting(t *testing.T) {
	// Setup temporary database
	tmpDB := "test.db"
	os.Setenv("DB_PATH", tmpDB)

	// Initialize DB
	InitDB()
	defer os.Remove(tmpDB)
	defer db.Close()

	// Insert some files and posts
	f := &File{
		PubID: "test_pub_id",
		Meta: FileMeta{
			OriginalName:  "test.txt",
			Size:          10,
			Hash:          "test_hash",
			LocalFileName: "test.txt",
			FileType:      "text/plain",
		},
		Issuer: "tester",
	}
	f2 := &File{
		PubID: "test_pub_id_2",
		Meta: FileMeta{
			OriginalName:  "test.txt",
			Size:          10,
			Hash:          "test_hash",
			LocalFileName: "test.txt",
			FileType:      "text/plain",
		},
		Issuer: "tester",
	}
	p := &Post{
		PubID:   "test_pub_id",
		Content: `{"files": ["12", "34"],"description": "description"}`,
		Issuer:  "tester",
	}
	p2 := &Post{
		PubID:   "test_pub_id_2",
		Content: `{"files": ["12", "34"],"description": "description"}`,
		Issuer:  "tester",
	}

	InsertFile(f)
	InsertFile(f)
	InsertFile(f)

	InsertPost(p)
	InsertPost(p)
	InsertPost(p)

	// Check file count
	count, err := GetFileEntries()
	if err != nil {
		t.Fatalf("Failed to get file count: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected file count to be 1, got %d", count)
	}

	// Check post count
	count, err = GetPostEntries()
	if err != nil {
		t.Fatalf("Failed to get post count: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected post count to be 1, got %d", count)
	}

	InsertFile(f2)
	InsertPost(p2)

	// Check file count
	count, err = GetFileEntries()
	if err != nil {
		t.Fatalf("Failed to get file count: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected file count to be 2, got %d", count)
	}

	// Check post count
	count, err = GetPostEntries()
	if err != nil {
		t.Fatalf("Failed to get post count: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected post count to be 2, got %d", count)
	}
}
