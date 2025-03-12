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

	jsonMap, err := scannerInstance.ParseReport(scanData.Report)
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
		"report": jsonMap,
	})
	if err != nil {
		return false, nil, fmt.Errorf("failed to evaluate CEL expression: %w", err)
	}

	// Extract severities from the condition dynamically
	severities := extractSeverities(scanCondition)

	// Return evaluation result and extracted severities
	return out == types.True, severities, nil
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
