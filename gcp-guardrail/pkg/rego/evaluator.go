package rego

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/storage/inmem"
)

// Violation represents a policy violation
type Violation struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
	Policy   string `json:"policy"`
}

// EvaluationResult contains the results of a policy evaluation
type EvaluationResult struct {
	Violations []Violation `json:"violations"`
	Warnings   []Violation `json:"warnings"`
	PassCount  int         `json:"pass_count"`
	FailCount  int         `json:"fail_count"`
	WarnCount  int         `json:"warn_count"`
}

// Evaluator is responsible for evaluating OPA policies
type Evaluator struct {
	policyDirs []string
	modules    map[string]*ast.Module
	ctx        context.Context
	store      storage.Store
}

// NewEvaluator creates a new OPA policy evaluator
func NewEvaluator(ctx context.Context, policyDirs []string) (*Evaluator, error) {
	e := &Evaluator{
		policyDirs: policyDirs,
		modules:    make(map[string]*ast.Module),
		ctx:        ctx,
		store:      inmem.New(),
	}

	if err := e.loadPolicies(); err != nil {
		return nil, err
	}

	return e, nil
}

// loadPolicies loads all .rego files from the specified directories
func (e *Evaluator) loadPolicies() error {
	for _, dir := range e.policyDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && strings.HasSuffix(path, ".rego") {
				// Read and parse the Rego file
				content, err := os.ReadFile(path)
				if err != nil {
					return fmt.Errorf("failed to read policy file %s: %w", path, err)
				}

				// Parse the module
				module, err := ast.ParseModule(path, string(content))
				if err != nil {
					return fmt.Errorf("failed to parse policy file %s: %w", path, err)
				}

				// Add the module to the map
				e.modules[path] = module
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to load policies from %s: %w", dir, err)
		}
	}

	return nil
}

// Evaluate evaluates the given input against loaded policies
func (e *Evaluator) Evaluate(packagePath string, input interface{}) (*EvaluationResult, error) {
	result := &EvaluationResult{}

	// Create a new Rego instance
	r := rego.New(
		rego.Query(fmt.Sprintf("data.%s", packagePath)),
		rego.Store(e.store),
		rego.Modules(e.modules),
		rego.Input(input),
	)

	// Run the evaluation
	rs, err := r.Eval(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("evaluation failed: %w", err)
	}

	// Process deny rules
	denyQuery := rego.New(
		rego.Query(fmt.Sprintf("data.%s.deny", packagePath)),
		rego.Store(e.store),
		rego.Modules(e.modules),
		rego.Input(input),
	)

	denyRs, err := denyQuery.Eval(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("deny rule evaluation failed: %w", err)
	}

	// Process warnings
	warnQuery := rego.New(
		rego.Query(fmt.Sprintf("data.%s.warn", packagePath)),
		rego.Store(e.store),
		rego.Modules(e.modules),
		rego.Input(input),
	)

	warnRs, err := warnQuery.Eval(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("warning rule evaluation failed: %w", err)
	}

	// Extract violations from deny rules
	if len(denyRs) > 0 && len(denyRs[0].Expressions) > 0 {
		violations := denyRs[0].Expressions[0].Value
		if violations != nil {
			for _, v := range violations.([]interface{}) {
				result.Violations = append(result.Violations, Violation{
					Message:  v.(string),
					Severity: "ERROR",
					Policy:   packagePath,
				})
			}
		}
	}

	// Extract warnings
	if len(warnRs) > 0 && len(warnRs[0].Expressions) > 0 {
		warnings := warnRs[0].Expressions[0].Value
		if warnings != nil {
			for _, w := range warnings.([]interface{}) {
				result.Warnings = append(result.Warnings, Violation{
					Message:  w.(string),
					Severity: "WARNING",
					Policy:   packagePath,
				})
			}
		}
	}

	result.FailCount = len(result.Violations)
	result.WarnCount = len(result.Warnings)
	result.PassCount = 0 // TODO: Calculate based on total rules - (failures + warnings)

	return result, nil
}

// EvaluateAll evaluates input against all loaded policies
func (e *Evaluator) EvaluateAll(input interface{}) ([]*EvaluationResult, error) {
	var results []*EvaluationResult

	// Get unique package paths from the loaded modules
	packagePaths := make(map[string]bool)
	for _, module := range e.modules {
		packagePath := strings.Join(module.Package.Path, ".")
		packagePaths[packagePath] = true
	}

	// Evaluate each package
	for packagePath := range packagePaths {
		result, err := e.Evaluate(packagePath, input)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate package %s: %w", packagePath, err)
		}
		results = append(results, result)
	}

	return results, nil
} 