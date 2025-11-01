package sqlite

import (
	"context"
	"database/sql"
	"math"
	"strconv"
)

// AdaptiveIDConfig holds configuration for adaptive ID length scaling
type AdaptiveIDConfig struct {
	// MaxCollisionProbability is the threshold at which we scale up ID length (e.g., 0.25 = 25%)
	MaxCollisionProbability float64

	// MinLength is the minimum hash length to use (default 4)
	MinLength int

	// MaxLength is the maximum hash length to use (default 8)
	MaxLength int
}

// DefaultAdaptiveConfig returns sensible defaults
func DefaultAdaptiveConfig() AdaptiveIDConfig {
	return AdaptiveIDConfig{
		MaxCollisionProbability: 0.25, // 25% threshold
		MinLength:               4,
		MaxLength:               8,
	}
}

// collisionProbability calculates P(collision) using birthday paradox approximation
// P(collision) ≈ 1 - e^(-n²/2N)
// where n = number of items, N = total possible values
func collisionProbability(numIssues int, idLength int) float64 {
	const base = 36.0 // lowercase alphanumeric (0-9, a-z)
	totalPossibilities := math.Pow(base, float64(idLength))
	exponent := -float64(numIssues*numIssues) / (2.0 * totalPossibilities)
	return 1.0 - math.Exp(exponent)
}

// computeAdaptiveLength determines the optimal ID length for the current database size
func computeAdaptiveLength(numIssues int, config AdaptiveIDConfig) int {
	// Try lengths from min to max, return first that meets threshold
	for length := config.MinLength; length <= config.MaxLength; length++ {
		prob := collisionProbability(numIssues, length)
		if prob <= config.MaxCollisionProbability {
			return length
		}
	}

	// If even maxLength doesn't meet threshold, return maxLength anyway
	return config.MaxLength
}

// getAdaptiveConfig reads adaptive ID config from database, returns defaults if not set
func getAdaptiveConfig(ctx context.Context, conn *sql.Conn) AdaptiveIDConfig {
	config := DefaultAdaptiveConfig()

	// Read max_collision_prob
	var probStr string
	err := conn.QueryRowContext(ctx, `SELECT value FROM config WHERE key = ?`, "max_collision_prob").Scan(&probStr)
	if err == nil && probStr != "" {
		if prob, err := strconv.ParseFloat(probStr, 64); err == nil {
			config.MaxCollisionProbability = prob
		}
	}

	// Read min_hash_length
	var minLenStr string
	err = conn.QueryRowContext(ctx, `SELECT value FROM config WHERE key = ?`, "min_hash_length").Scan(&minLenStr)
	if err == nil && minLenStr != "" {
		if minLen, err := strconv.Atoi(minLenStr); err == nil {
			config.MinLength = minLen
		}
	}

	// Read max_hash_length
	var maxLenStr string
	err = conn.QueryRowContext(ctx, `SELECT value FROM config WHERE key = ?`, "max_hash_length").Scan(&maxLenStr)
	if err == nil && maxLenStr != "" {
		if maxLen, err := strconv.Atoi(maxLenStr); err == nil {
			config.MaxLength = maxLen
		}
	}

	return config
}

// countTopLevelIssues returns the number of top-level issues (excluding child issues)
func countTopLevelIssues(ctx context.Context, conn *sql.Conn, prefix string) (int, error) {
	var count int
	// Count only top-level issues (no dot in ID after prefix)
	err := conn.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM issues
		WHERE id LIKE ? || '-%'
		  AND instr(substr(id, length(?) + 2), '.') = 0
	`, prefix, prefix).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetAdaptiveIDLength returns the appropriate hash length based on database size
func GetAdaptiveIDLength(ctx context.Context, conn *sql.Conn, prefix string) (int, error) {
	// Get current issue count
	numIssues, err := countTopLevelIssues(ctx, conn, prefix)
	if err != nil {
		return 6, err // Fallback to 6 on error
	}

	// Get adaptive config
	config := getAdaptiveConfig(ctx, conn)

	// Compute optimal length
	length := computeAdaptiveLength(numIssues, config)

	return length, nil
}
