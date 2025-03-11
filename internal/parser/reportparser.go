package parser

import (
	"fmt"
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/shinobistack/gokakashi/internal/restapi/v1/scans"
	"github.com/shinobistack/gokakashi/pkg/scanner/v1"
	"regexp"
)

// TodO: Need to better this logic.

func ReportParser(scanCondition string, scanData *scans.GetScanResponse) (bool, []string, error) {
	scannerInstance, err := scanner.NewScanner(scanData.Scanner)
	if err != nil {
		return false, nil, fmt.Errorf("unsupported scanner: %s", scanData.Scanner)
	}

	reportMap, err := scannerInstance.ParseReport(scanData.Report)
	if err != nil {
		return false, nil, err
	}

	// Initialize CEL environment
	env, err := cel.NewEnv(cel.Variable("report", cel.DynType))
	if err != nil {
		return false, nil, fmt.Errorf("failed to create CEL environment: %w", err)
	}

	// Compile CEL expression
	parsed, issues := env.Parse(scanCondition)
	if issues.Err() != nil {
		return false, nil, fmt.Errorf("failed to parse CEL condition: %w", issues.Err())
	}

	checked, issues := env.Check(parsed)
	if issues.Err() != nil {
		return false, nil, fmt.Errorf("failed to check CEL condition: %w", issues.Err())
	}

	prg, err := env.Program(checked)
	if err != nil {
		return false, nil, fmt.Errorf("failed to create CEL program: %w", err)
	}

	// Evaluate the condition
	out, _, err := prg.Eval(map[string]interface{}{
		"report": reportMap,
	})
	if err != nil {
		return false, nil, fmt.Errorf("failed to evaluate CEL expression: %w", err)
	}

	// Extract severities from the condition dynamically
	severities := extractSeverities(scanCondition)

	// Return evaluation result and extracted severities
	return out == types.True, severities, nil
}

// ValidateCELExpression checks if the CEL expression references fields that exist in the JSON report
func ValidateCELExpression(expression string, jsonData map[string]interface{}) error {
	// Extract fields used in the CEL expression
	fields := ExtractFieldsFromCEL(expression)
	availableKeys := FlattenJSONKeys(jsonData, "")

	for _, field := range fields {
		if _, exists := availableKeys[field]; !exists {
			return fmt.Errorf("CEL expression references missing field: '%s'. This field may not exist in the scanner's JSON structure. Please use a compatible CEL expression", field)
		}
	}

	return nil

}

func ExtractFieldsFromCEL(expression string) []string {
	// Match words that look like field names (e.g., report.Results, v.Severity)
	re := regexp.MustCompile(`\b[a-zA-Z_][a-zA-Z0-9_]*\b`)
	matches := re.FindAllString(expression, -1)
	// fmt.Print(matches, "matches\n")

	// Remove CEL operators and keywords (e.g., exists, &&, ||, CRITICAL, HIGH)
	ignoredWords := map[string]bool{
		"exists": true, "true": true, "false": true, "in": true, "and": true, "or": true, "CRITICAL": true, "HIGH": true, "MEDIUM": true, "LOW": true, "UNKNOWN": true,
	}

	var fields []string
	for _, match := range matches {
		// Ignore CEL keywords, single-letter variables, and "report"
		if !ignoredWords[match] && match != "report" && len(match) > 1 {
			fields = append(fields, match)
		}
	}
	return fields
}

func FlattenJSONKeys(data interface{}, prefix string) map[string]bool {
	keys := make(map[string]bool)

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fullKey := key
			if prefix != "" {
				fullKey = prefix + "." + key
			}
			keys[key] = true     // Store key as is
			keys[fullKey] = true // Store key with full path for flexibility

			// Recurse into nested objects
			for subKey := range FlattenJSONKeys(value, fullKey) {
				keys[subKey] = true
			}
		}

	case []interface{}:
		for _, item := range v {
			for subKey := range FlattenJSONKeys(item, prefix) {
				keys[subKey] = true
			}
		}
	}

	return keys
}

// extractSeverities parses the CEL condition to find severity levels
func extractSeverities(condition string) []string {
	// Regular expression to match severity values (e.g., 'CRITICAL', 'HIGH')
	re := regexp.MustCompile(`v\.Severity\s*==\s*'([A-Z]+)'`)
	matches := re.FindAllStringSubmatch(condition, -1)

	// Extract unique severity levels
	severityMap := make(map[string]struct{})
	for _, match := range matches {
		if len(match) > 1 {
			severityMap[match[1]] = struct{}{}
		}
	}

	// Convert map keys to a slice
	var severities []string
	for severity := range severityMap {
		severities = append(severities, severity)
	}
	return severities
}
