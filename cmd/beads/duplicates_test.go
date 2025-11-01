package main

import (
	"context"
	"testing"

	"github.com/shaneholloman/beads/internal/types"
)

func TestFindDuplicateGroups(t *testing.T) {
	tests := []struct {
		name           string
		issues         []*types.Issue
		expectedGroups int
	}{
		{
			name: "no duplicates",
			issues: []*types.Issue{
				{ID: "beads-1", Title: "Task 1", Status: types.StatusOpen},
				{ID: "beads-2", Title: "Task 2", Status: types.StatusOpen},
			},
			expectedGroups: 0,
		},
		{
			name: "simple duplicate",
			issues: []*types.Issue{
				{ID: "beads-1", Title: "Task 1", Status: types.StatusOpen},
				{ID: "beads-2", Title: "Task 1", Status: types.StatusOpen},
			},
			expectedGroups: 1,
		},
		{
			name: "duplicate with different status ignored",
			issues: []*types.Issue{
				{ID: "beads-1", Title: "Task 1", Status: types.StatusOpen},
				{ID: "beads-2", Title: "Task 1", Status: types.StatusClosed},
			},
			expectedGroups: 0,
		},
		{
			name: "multiple duplicates",
			issues: []*types.Issue{
				{ID: "beads-1", Title: "Task 1", Status: types.StatusOpen},
				{ID: "beads-2", Title: "Task 1", Status: types.StatusOpen},
				{ID: "beads-3", Title: "Task 2", Status: types.StatusOpen},
				{ID: "beads-4", Title: "Task 2", Status: types.StatusOpen},
			},
			expectedGroups: 2,
		},
		{
			name: "different descriptions are duplicates if title matches",
			issues: []*types.Issue{
				{ID: "beads-1", Title: "Task 1", Description: "Desc 1", Status: types.StatusOpen},
				{ID: "beads-2", Title: "Task 1", Description: "Desc 2", Status: types.StatusOpen},
			},
			expectedGroups: 0, // Different descriptions = not duplicates
		},
		{
			name: "exact content match",
			issues: []*types.Issue{
				{ID: "beads-1", Title: "Task 1", Description: "Desc 1", Design: "Design 1", AcceptanceCriteria: "AC 1", Status: types.StatusOpen},
				{ID: "beads-2", Title: "Task 1", Description: "Desc 1", Design: "Design 1", AcceptanceCriteria: "AC 1", Status: types.StatusOpen},
			},
			expectedGroups: 1,
		},
		{
			name: "three-way duplicate",
			issues: []*types.Issue{
				{ID: "beads-1", Title: "Task 1", Status: types.StatusOpen},
				{ID: "beads-2", Title: "Task 1", Status: types.StatusOpen},
				{ID: "beads-3", Title: "Task 1", Status: types.StatusOpen},
			},
			expectedGroups: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := findDuplicateGroups(tt.issues)
			if len(groups) != tt.expectedGroups {
				t.Errorf("findDuplicateGroups() returned %d groups, want %d", len(groups), tt.expectedGroups)
			}
		})
	}
}

func TestChooseMergeTarget(t *testing.T) {
	tests := []struct {
		name      string
		group     []*types.Issue
		refCounts map[string]int
		wantID    string
	}{
		{
			name: "choose by reference count",
			group: []*types.Issue{
				{ID: "beads-2", Title: "Task"},
				{ID: "beads-1", Title: "Task"},
			},
			refCounts: map[string]int{
				"beads-1": 5,
				"beads-2": 0,
			},
			wantID: "beads-1",
		},
		{
			name: "choose by lexicographic order if same references",
			group: []*types.Issue{
				{ID: "beads-2", Title: "Task"},
				{ID: "beads-1", Title: "Task"},
			},
			refCounts: map[string]int{
				"beads-1": 0,
				"beads-2": 0,
			},
			wantID: "beads-1",
		},
		{
			name: "prefer higher references even with larger ID",
			group: []*types.Issue{
				{ID: "beads-1", Title: "Task"},
				{ID: "beads-100", Title: "Task"},
			},
			refCounts: map[string]int{
				"beads-1":   1,
				"beads-100": 10,
			},
			wantID: "beads-100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := chooseMergeTarget(tt.group, tt.refCounts)
			if target.ID != tt.wantID {
				t.Errorf("chooseMergeTarget() = %v, want %v", target.ID, tt.wantID)
			}
		})
	}
}

func TestCountReferences(t *testing.T) {
	issues := []*types.Issue{
		{
			ID:          "beads-1",
			Description: "See beads-2 for details",
			Notes:       "Related to beads-3",
		},
		{
			ID:          "beads-2",
			Description: "Mentioned beads-1 twice: beads-1",
		},
		{
			ID:    "beads-3",
			Notes: "Nothing to see here",
		},
	}

	counts := countReferences(issues)

	expectedCounts := map[string]int{
		"beads-1": 2, // Referenced twice in beads-2
		"beads-2": 1, // Referenced once in beads-1
		"beads-3": 1, // Referenced once in beads-1
	}

	for id, expectedCount := range expectedCounts {
		if counts[id] != expectedCount {
			t.Errorf("countReferences()[%s] = %d, want %d", id, counts[id], expectedCount)
		}
	}
}

func TestDuplicateGroupsWithDifferentStatuses(t *testing.T) {
	issues := []*types.Issue{
		{ID: "beads-1", Title: "Task 1", Status: types.StatusOpen},
		{ID: "beads-2", Title: "Task 1", Status: types.StatusClosed},
		{ID: "beads-3", Title: "Task 1", Status: types.StatusOpen},
	}

	groups := findDuplicateGroups(issues)

	// Should have 1 group with beads-1 and beads-3 (both open)
	if len(groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(groups))
	}

	if len(groups[0]) != 2 {
		t.Fatalf("Expected 2 issues in group, got %d", len(groups[0]))
	}

	// Verify beads-2 (closed) is not in the group
	for _, issue := range groups[0] {
		if issue.ID == "beads-2" {
			t.Errorf("beads-2 (closed) should not be in group with open issues")
		}
	}
}

func TestDuplicatesIntegration(t *testing.T) {
	ctx := context.Background()
	testStore, cleanup := setupTestDB(t)
	defer cleanup()

	// Create duplicate issues
	issues := []*types.Issue{
		{
			ID:          "beads-1",
			Title:       "Fix authentication bug",
			Description: "Users can't login",
			Status:      types.StatusOpen,
			Priority:    1,
			IssueType:   types.TypeBug,
		},
		{
			ID:          "beads-2",
			Title:       "Fix authentication bug",
			Description: "Users can't login",
			Status:      types.StatusOpen,
			Priority:    1,
			IssueType:   types.TypeBug,
		},
		{
			ID:          "beads-3",
			Title:       "Different task",
			Description: "Different description",
			Status:      types.StatusOpen,
			Priority:    2,
			IssueType:   types.TypeTask,
		},
	}

	for _, issue := range issues {
		if err := testStore.CreateIssue(ctx, issue, "test"); err != nil {
			t.Fatalf("CreateIssue failed: %v", err)
		}
	}

	// Fetch all issues
	allIssues, err := testStore.SearchIssues(ctx, "", types.IssueFilter{})
	if err != nil {
		t.Fatalf("SearchIssues failed: %v", err)
	}

	// Find duplicates
	groups := findDuplicateGroups(allIssues)

	if len(groups) != 1 {
		t.Fatalf("Expected 1 duplicate group, got %d", len(groups))
	}

	if len(groups[0]) != 2 {
		t.Fatalf("Expected 2 issues in group, got %d", len(groups[0]))
	}

	// Verify the duplicate group contains beads-1 and beads-2
	ids := make(map[string]bool)
	for _, issue := range groups[0] {
		ids[issue.ID] = true
	}

	if !ids["beads-1"] || !ids["beads-2"] {
		t.Errorf("Expected duplicate group to contain beads-1 and beads-2")
	}
}
