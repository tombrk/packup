package restic

import "testing"

func TestParseBackupResultSummary(t *testing.T) {
	out := []byte(`{"message_type":"status","total_files":1,"files_done":1}
{"message_type":"summary","files_new":1,"total_files_processed":1,"total_bytes_processed":42,"data_added":7,"data_added_packed":5,"snapshot_id":"abc123"}
`)

	result, err := parseBackupResult(out)
	if err != nil {
		t.Fatal(err)
	}
	if result.Summary == nil {
		t.Fatal("expected summary")
	}
	if result.Summary.SnapshotID != "abc123" {
		t.Fatalf("snapshot id = %q", result.Summary.SnapshotID)
	}
	if result.Summary.TotalBytesProcessed != 42 {
		t.Fatalf("total bytes = %d", result.Summary.TotalBytesProcessed)
	}
	if result.Summary.DataAddedPacked != 5 {
		t.Fatalf("packed bytes = %d", result.Summary.DataAddedPacked)
	}
}

func TestParseBackupResultNoSummary(t *testing.T) {
	result, err := parseBackupResult([]byte(`{"message_type":"status","total_files":1}
`))
	if err != nil {
		t.Fatal(err)
	}
	if result.Summary != nil {
		t.Fatalf("unexpected summary: %+v", result.Summary)
	}
}
