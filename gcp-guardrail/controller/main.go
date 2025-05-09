package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hieptle/gcp-guardrail/pkg/rego"
	
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
	
	// Policy paths
	defaultPolicyDirs = []string{
		"policies/kubernetes",
	}
)

type admissionController struct {
	evaluator *rego.Evaluator
}

// Initialize the admission controller
func newAdmissionController(policyDirs []string) (*admissionController, error) {
	evaluator, err := rego.NewEvaluator(context.Background(), policyDirs)
	if err != nil {
		return nil, fmt.Errorf("failed to create policy evaluator: %w", err)
	}
	
	return &admissionController{
		evaluator: evaluator,
	}, nil
}

// Handle the admission review request
func (ac *admissionController) handleAdmissionRequest(w http.ResponseWriter, r *http.Request) {
	var body []byte
	if r.Body != nil {
		if data, err := r.Body.Read(body); err == nil {
			body = data
		}
	}
	
	// Verify the content type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		log.Printf("contentType=%s, expected application/json", contentType)
		http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	
	// Parse the AdmissionReview request
	reviewRequest := admissionv1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &reviewRequest); err != nil {
		log.Printf("Could not decode body: %v", err)
		http.Error(w, "Invalid AdmissionReview request", http.StatusBadRequest)
		return
	}
	
	// Initialize response
	reviewResponse := admissionv1.AdmissionReview{
		TypeMeta: reviewRequest.TypeMeta,
		Response: &admissionv1.AdmissionResponse{
			UID: reviewRequest.Request.UID,
		},
	}
	
	// Evaluate the request against policies
	allowed, reason, err := ac.evaluateRequest(reviewRequest.Request)
	if err != nil {
		log.Printf("Error evaluating request: %v", err)
		reviewResponse.Response.Allowed = false
		reviewResponse.Response.Result = &metav1.Status{
			Message: fmt.Sprintf("Error evaluating request: %v", err),
		}
	} else {
		reviewResponse.Response.Allowed = allowed
		if !allowed {
			reviewResponse.Response.Result = &metav1.Status{
				Message: reason,
			}
		}
	}
	
	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviewResponse)
}

// Evaluate the admission request against the policies
func (ac *admissionController) evaluateRequest(request *admissionv1.AdmissionRequest) (bool, string, error) {
	// Convert the request to a format the OPA evaluator can process
	input := map[string]interface{}{
		"kind":       request.Kind.Kind,
		"name":       request.Name,
		"namespace":  request.Namespace,
		"operation":  request.Operation,
		"object":     request.Object.Raw,
		"oldObject":  request.OldObject.Raw,
		"parameters": request.Options,
	}
	
	// Determine the package path based on the resource kind
	packagePath := fmt.Sprintf("kubernetes.admission.%s", normalizeKind(request.Kind.Kind))
	
	// Evaluate the policies
	result, err := ac.evaluator.Evaluate(packagePath, input)
	if err != nil {
		return false, "", fmt.Errorf("policy evaluation failed: %w", err)
	}
	
	// If there are any violations, deny the request
	if len(result.Violations) > 0 {
		var reason string
		for i, violation := range result.Violations {
			if i == 0 {
				reason = violation.Message
			} else {
				reason = fmt.Sprintf("%s; %s", reason, violation.Message)
			}
		}
		return false, reason, nil
	}
	
	return true, "", nil
}

// Normalize the Kubernetes kind name
func normalizeKind(kind string) string {
	return strings.ToLower(kind)
}

// Server health check handler
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	// Get policy directories from environment variable or use default
	policyDirsEnv := os.Getenv("POLICY_DIRS")
	policyDirs := defaultPolicyDirs
	if policyDirsEnv != "" {
		policyDirs = filepath.SplitList(policyDirsEnv)
	}
	
	// Create the admission controller
	ac, err := newAdmissionController(policyDirs)
	if err != nil {
		log.Fatalf("Failed to create admission controller: %v", err)
	}
	
	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/validate", ac.handleAdmissionRequest)
	
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8443"
	}
	
	// Start the HTTP server
	log.Printf("Starting server on port %s", port)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}
	
	// Start the server with TLS if certificates are provided
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	if certFile != "" && keyFile != "" {
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		log.Fatal(server.ListenAndServe())
	}
} 