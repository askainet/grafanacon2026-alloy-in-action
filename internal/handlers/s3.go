package handlers

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type S3Handler struct {
	endpoint     string
	bucket       string
	lokiEndpoint string
	client       *http.Client
}

func NewS3Handler(endpoint, bucket, lokiEndpoint string) *S3Handler {
	return &S3Handler{
		endpoint:     endpoint,
		bucket:       bucket,
		lokiEndpoint: lokiEndpoint,
		client:       &http.Client{Timeout: 15 * time.Second},
	}
}

type listBucketResult struct {
	Contents []s3Object `xml:"Contents"`
}

type s3Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	Size         int64  `xml:"Size"`
}

// ListObjects lists objects and returns their contents from the S3 bucket.
func (h *S3Handler) ListObjects(w http.ResponseWriter, r *http.Request) {
	objects, err := h.listObjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	if len(objects) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "no objects found in bucket",
			"bucket":  h.bucket,
		})
		return
	}

	// Sort by key descending (most recent first, since keys include timestamps)
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].Key > objects[j].Key
	})

	limit := min(10, len(objects))

	type objectEntry struct {
		Key     string `json:"key"`
		Content string `json:"content"`
	}

	entries := make([]objectEntry, 0, limit)
	for _, obj := range objects[:limit] {
		body, err := h.getObject(obj.Key)
		if err != nil {
			entries = append(entries, objectEntry{Key: obj.Key, Content: "error: " + err.Error()})
			continue
		}
		entries = append(entries, objectEntry{Key: obj.Key, Content: string(body)})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"bucket":        h.bucket,
		"total_objects": len(objects),
		"showing":       limit,
		"objects":       entries,
	})
}

// Verify checks both S3 (all logs present) and Loki (no DEBUG logs leaking).
func (h *S3Handler) Verify(w http.ResponseWriter, r *http.Request) {
	s3Result := h.verifyS3()
	lokiResult := h.verifyLoki()

	// Overall pass requires both checks to pass
	s3Pass := s3Result.Status == "pass"
	lokiPass := lokiResult.Status == "pass"

	w.Header().Set("Content-Type", "application/json")
	if s3Pass && lokiPass {
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "pass",
			"message": "S3 receiving all logs and Loki has no DEBUG logs - routing is working correctly",
		})
	} else {
		var messages []string
		messages = append(messages, "S3: "+s3Result.Message)
		messages = append(messages, "Loki: "+lokiResult.Message)
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "fail",
			"message": strings.Join(messages, "; "),
			"s3":      s3Result,
			"loki":    lokiResult,
		})
	}
}

type s3VerifyResult struct {
	Status       string         `json:"status"`
	Message      string         `json:"message"`
	TotalObjects int            `json:"total_objects"`
	TotalLines   int            `json:"total_lines"`
	Levels       map[string]int `json:"levels"`
}

func (h *S3Handler) verifyS3() s3VerifyResult {
	objects, err := h.listObjects()
	if err != nil {
		return s3VerifyResult{Status: "fail", Message: "failed to list S3 objects: " + err.Error()}
	}

	if len(objects) == 0 {
		return s3VerifyResult{Status: "fail", Message: "no objects found in S3 bucket - are logs being routed?"}
	}

	levelCounts := map[string]int{}
	totalLines := 0

	for _, obj := range objects {
		body, err := h.getObject(obj.Key)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(body), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			totalLines++
			var entry struct {
				Level string `json:"level"`
			}
			if json.Unmarshal([]byte(line), &entry) == nil && entry.Level != "" {
				levelCounts[entry.Level]++
			} else {
				levelCounts["UNKNOWN"]++
			}
		}
	}

	hasDebug := levelCounts["DEBUG"] > 0
	hasOther := false
	for level, count := range levelCounts {
		if level != "DEBUG" && count > 0 {
			hasOther = true
			break
		}
	}

	result := s3VerifyResult{
		TotalObjects: len(objects),
		TotalLines:   totalLines,
		Levels:       levelCounts,
	}

	switch {
	case hasDebug && hasOther:
		result.Status = "pass"
		result.Message = "S3 contains both DEBUG and non-DEBUG logs - routing is working correctly"
	case hasDebug:
		result.Status = "fail"
		result.Message = "S3 contains DEBUG logs but no other levels - check that all logs are being forwarded to S3"
	case hasOther:
		result.Status = "fail"
		result.Message = "S3 contains non-DEBUG logs but no DEBUG logs - is mission2 active?"
	default:
		result.Status = "fail"
		result.Message = "could not parse any log levels from S3 objects"
	}

	return result
}

type lokiVerifyResult struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	DebugCount int    `json:"debug_count"`
}

func (h *S3Handler) verifyLoki() lokiVerifyResult {
	// Query Loki for DEBUG logs in the last 2 minutes.
	// Uses filename label (always present from loki.source.file) so this works
	// even if the participant didn't extract other labels.
	query := `sum(count_over_time({filename=~".+"} | json | level = "DEBUG" [30s]))`
	queryURL := fmt.Sprintf("%s/loki/api/v1/query?query=%s", h.lokiEndpoint, url.QueryEscape(query))

	resp, err := h.client.Get(queryURL)
	if err != nil {
		return lokiVerifyResult{Status: "fail", Message: "failed to query Loki: " + err.Error()}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return lokiVerifyResult{Status: "fail", Message: "failed to read Loki response: " + err.Error()}
	}

	if resp.StatusCode != http.StatusOK {
		return lokiVerifyResult{Status: "fail", Message: fmt.Sprintf("Loki returned status %d", resp.StatusCode)}
	}

	var lokiResp struct {
		Status string `json:"status"`
		Data   struct {
			Result []struct {
				Value json.RawMessage `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &lokiResp); err != nil {
		return lokiVerifyResult{Status: "fail", Message: "failed to parse Loki response: " + err.Error()}
	}

	if lokiResp.Status != "success" {
		return lokiVerifyResult{Status: "fail", Message: "Loki query failed: " + lokiResp.Status}
	}

	// Parse the count from the vector result: [timestamp, "count"]
	debugCount := 0
	if len(lokiResp.Data.Result) > 0 {
		var valuePair [2]json.RawMessage
		if json.Unmarshal(lokiResp.Data.Result[0].Value, &valuePair) == nil {
			var countStr string
			if json.Unmarshal(valuePair[1], &countStr) == nil {
				fmt.Sscanf(countStr, "%d", &debugCount)
			}
		}
	}

	if debugCount > 0 {
		return lokiVerifyResult{
			Status:     "fail",
			Message:    fmt.Sprintf("%d DEBUG log(s) found in Loki in the last 2 minutes - stage.drop filter needs work", debugCount),
			DebugCount: debugCount,
		}
	}

	return lokiVerifyResult{
		Status:  "pass",
		Message: "no DEBUG logs found in Loki - filtering is working correctly",
	}
}

func (h *S3Handler) listObjects() ([]s3Object, error) {
	listURL := fmt.Sprintf("%s/%s", h.endpoint, h.bucket)
	resp, err := http.Get(listURL)
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 objects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("S3 list failed (%d): %s", resp.StatusCode, string(body))
	}

	var result listBucketResult
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse S3 response: %w", err)
	}
	return result.Contents, nil
}

func (h *S3Handler) getObject(key string) ([]byte, error) {
	objURL := fmt.Sprintf("%s/%s/%s", h.endpoint, h.bucket, key)
	resp, err := http.Get(objURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
