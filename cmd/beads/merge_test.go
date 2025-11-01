package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/shaneholloman/beads/internal/types"
)

func TestValidateMerge(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, ".beads", "issues.db")

	testStore := newTestStoreWithPrefix(t, dbFile, "beads")
	store = testStore
	ctx := context.Background()

	// Create test issues
	issue1 := &types.Issue{
		ID:          "beads-1",
		Title:       "Test issue 1",
		Description: "Test",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}
	issue2 := &types.Issue{
		ID:          "beads-2",
		Title:       "Test issue 2",
		Description: "Test",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}
	issue3 := &types.Issue{
		ID:          "beads-3",
		Title:       "Test issue 3",
		Description: "Test",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}

	if err := testStore.CreateIssue(ctx, issue1, "beads"); err != nil {
		t.Fatalf("Failed to create issue1: %v", err)
	}
	if err := testStore.CreateIssue(ctx, issue2, "beads"); err != nil {
		t.Fatalf("Failed to create issue2: %v", err)
	}
	if err := testStore.CreateIssue(ctx, issue3, "beads"); err != nil {
		t.Fatalf("Failed to create issue3: %v", err)
	}

	tests := []struct {
		name      string
		targetID  string
		sourceIDs []string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid merge",
			targetID:  "beads-1",
			sourceIDs: []string{"beads-2", "beads-3"},
			wantErr:   false,
		},
		{
			name:      "self-merge error",
			targetID:  "beads-1",
			sourceIDs: []string{"beads-1"},
			wantErr:   true,
			errMsg:    "cannot merge issue into itself",
		},
		{
			name:      "self-merge in list",
			targetID:  "beads-1",
			sourceIDs: []string{"beads-2", "beads-1"},
			wantErr:   true,
			errMsg:    "cannot merge issue into itself",
		},
		{
			name:      "nonexistent target",
			targetID:  "beads-999",
			sourceIDs: []string{"beads-1"},
			wantErr:   true,
			errMsg:    "target issue not found",
		},
		{
			name:      "nonexistent source",
			targetID:  "beads-1",
			sourceIDs: []string{"beads-999"},
			wantErr:   true,
			errMsg:    "source issue not found",
		},
		{
			name:      "multiple sources valid",
			targetID:  "beads-1",
			sourceIDs: []string{"beads-2"},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMerge(tt.targetID, tt.sourceIDs)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateMerge() expected error, got nil")
				} else if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("validateMerge() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateMerge() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateMergeMultipleSelfReferences(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, ".beads", "issues.db")

	testStore := newTestStoreWithPrefix(t, dbFile, "beads")
	store = testStore
	ctx := context.Background()

	issue1 := &types.Issue{
		ID:          "beads-10",
		Title:       "Test issue 10",
		Description: "Test",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}

	if err := testStore.CreateIssue(ctx, issue1, "beads"); err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	// Test merging multiple instances of same ID (should catch first one)
	err := validateMerge("beads-10", []string{"beads-10", "beads-10"})
	if err == nil {
		t.Error("validateMerge() expected error for duplicate self-merge, got nil")
	}
	if !contains(err.Error(), "cannot merge issue into itself") {
		t.Errorf("validateMerge() error = %v, want error containing 'cannot merge issue into itself'", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestPerformMergeIdempotent verifies that merge operations are idempotent
func TestPerformMergeIdempotent(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, ".beads", "issues.db")

	testStore := newTestStoreWithPrefix(t, dbFile, "beads")
	store = testStore
	ctx := context.Background()

	// Create test issues
	issue1 := &types.Issue{
		ID:          "beads-100",
		Title:       "Target issue",
		Description: "This is the target",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}
	issue2 := &types.Issue{
		ID:          "beads-101",
		Title:       "Source issue 1",
		Description: "This mentions beads-100",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}
	issue3 := &types.Issue{
		ID:          "beads-102",
		Title:       "Source issue 2",
		Description: "Another source",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}

	for _, issue := range []*types.Issue{issue1, issue2, issue3} {
		if err := testStore.CreateIssue(ctx, issue, "beads"); err != nil {
			t.Fatalf("Failed to create issue %s: %v", issue.ID, err)
		}
	}

	// Add a dependency from beads-101 to another issue
	issue4 := &types.Issue{
		ID:          "beads-103",
		Title:       "Dependency target",
		Description: "Dependency target",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}
	if err := testStore.CreateIssue(ctx, issue4, "beads"); err != nil {
		t.Fatalf("Failed to create issue4: %v", err)
	}

	dep := &types.Dependency{
		IssueID:     "beads-101",
		DependsOnID: "beads-103",
		Type:        types.DepBlocks,
	}
	if err := testStore.AddDependency(ctx, dep, "test"); err != nil {
		t.Fatalf("Failed to add dependency: %v", err)
	}

	// First merge - should complete successfully
	result1, err := performMerge(ctx, "beads-100", []string{"beads-101", "beads-102"})
	if err != nil {
		t.Fatalf("First merge failed: %v", err)
	}

	if result1.issuesClosed != 2 {
		t.Errorf("First merge: expected 2 issues closed, got %d", result1.issuesClosed)
	}
	if result1.issuesSkipped != 0 {
		t.Errorf("First merge: expected 0 issues skipped, got %d", result1.issuesSkipped)
	}
	if result1.depsAdded == 0 {
		t.Errorf("First merge: expected some dependencies added, got 0")
	}

	// Verify issues are closed
	closed1, _ := testStore.GetIssue(ctx, "beads-101")
	if closed1.Status != types.StatusClosed {
		t.Errorf("beads-101 should be closed after first merge")
	}
	closed2, _ := testStore.GetIssue(ctx, "beads-102")
	if closed2.Status != types.StatusClosed {
		t.Errorf("beads-102 should be closed after first merge")
	}

	// Second merge (retry) - should be idempotent
	result2, err := performMerge(ctx, "beads-100", []string{"beads-101", "beads-102"})
	if err != nil {
		t.Fatalf("Second merge (retry) failed: %v", err)
	}

	// All operations should be skipped
	if result2.issuesClosed != 0 {
		t.Errorf("Second merge: expected 0 issues closed, got %d", result2.issuesClosed)
	}
	if result2.issuesSkipped != 2 {
		t.Errorf("Second merge: expected 2 issues skipped, got %d", result2.issuesSkipped)
	}

	// Dependencies should be skipped (already exist)
	if result2.depsAdded != 0 {
		t.Errorf("Second merge: expected 0 dependencies added, got %d", result2.depsAdded)
	}

	// Text references are naturally idempotent - count may vary
	// (it will update again but result is the same)
}

// TestPerformMergePartialRetry tests retrying after partial failure
func TestPerformMergePartialRetry(t *testing.T) {
	tmpDir := t.TempDir()
	dbFile := filepath.Join(tmpDir, ".beads", "issues.db")

	testStore := newTestStoreWithPrefix(t, dbFile, "beads")
	store = testStore
	ctx := context.Background()

	// Create test issues
	issue1 := &types.Issue{
		ID:          "beads-200",
		Title:       "Target",
		Description: "Target issue",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}
	issue2 := &types.Issue{
		ID:          "beads-201",
		Title:       "Source 1",
		Description: "Source 1",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}
	issue3 := &types.Issue{
		ID:          "beads-202",
		Title:       "Source 2",
		Description: "Source 2",
		Priority:    1,
		IssueType:   types.TypeTask,
		Status:      types.StatusOpen,
	}

	for _, issue := range []*types.Issue{issue1, issue2, issue3} {
		if err := testStore.CreateIssue(ctx, issue, "beads"); err != nil {
			t.Fatalf("Failed to create issue %s: %v", issue.ID, err)
		}
	}

	// Simulate partial failure: manually close one source issue
	if err := testStore.CloseIssue(ctx, "beads-201", "Manually closed", "beads"); err != nil {
		t.Fatalf("Failed to manually close beads-201: %v", err)
	}

	// Run merge - should handle one already-closed issue gracefully
	result, err := performMerge(ctx, "beads-200", []string{"beads-201", "beads-202"})
	if err != nil {
		t.Fatalf("Merge with partial state failed: %v", err)
	}

	// Should skip the already-closed issue and close the other
	if result.issuesClosed != 1 {
		t.Errorf("Expected 1 issue closed, got %d", result.issuesClosed)
	}
	if result.issuesSkipped != 1 {
		t.Errorf("Expected 1 issue skipped, got %d", result.issuesSkipped)
	}

	// Verify both are now closed
	closed1, _ := testStore.GetIssue(ctx, "beads-201")
	if closed1.Status != types.StatusClosed {
		t.Errorf("beads-201 should remain closed")
	}
	closed2, _ := testStore.GetIssue(ctx, "beads-202")
	if closed2.Status != types.StatusClosed {
		t.Errorf("beads-202 should be closed")
	}
}
