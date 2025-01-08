package parser_test

//import (
//	"encoding/json"
//	"fmt"
//	"log"
//
//	"github.com/google/cel-go/cel"
//	"github.com/google/cel-go/common/types"
//)
//
//func main() {
//	// Sample scan report
//	rawJSON := `{
//		"Results": [
//			{
//				"Target": "library/busybox",
//				"Vulnerabilities": [
//					{
//						"VulnerabilityID": "CVE-2022-48174",
//						"Severity": "CRITICAL"
//					}
//				]
//			}
//		]
//	}`
//
//	var report map[string]interface{}
//	if err := json.Unmarshal([]byte(rawJSON), &report); err != nil {
//		log.Fatalf("Failed to parse JSON: %v", err)
//	}
//
//	// Condition to evaluate
//	expr := `scan.Results[0].Vulnerabilities.exists(v, v.Severity == "CRITICAL")`
//
//	// Create CEL environment
//	env, err := cel.NewEnv(cel.Variable("scan", cel.DynType))
//	if err != nil {
//		log.Fatalf("Failed to create CEL environment: %v", err)
//	}
//
//	// Compile expression
//	parsed, issues := env.Parse(expr)
//	if issues.Err() != nil {
//		log.Fatalf("Failed to parse CEL expression: %v", issues.Err())
//	}
//
//	checked, issues := env.Check(parsed)
//	if issues.Err() != nil {
//		log.Fatalf("Failed to check CEL expression: %v", issues.Err())
//	}
//
//	prg, err := env.Program(checked)
//	if err != nil {
//		log.Fatalf("Failed to create CEL program: %v", err)
//	}
//
//	// Evaluate the condition
//	out, _, err := prg.Eval(map[string]interface{}{
//		"scan": report,
//	})
//	if err != nil {
//		log.Fatalf("Failed to evaluate CEL expression: %v", err)
//	}
//
//	fmt.Printf("Condition evaluated to: %v\n", out == types.True)
//}
