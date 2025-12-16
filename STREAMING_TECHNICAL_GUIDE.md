# Technical Guide: Streaming with itemSchema and itemEncoding

🌍 **English** • [Português (Brasil)](STREAMING_TECHNICAL_GUIDE_pt.md) • [Español](STREAMING_TECHNICAL_GUIDE_es.md)

**Document Version:** 1.0  
**Last Updated:** December 15, 2025  
**OpenAPI Version:** 3.2.0+

---

## Table of Contents

1. [Overview](#overview)
2. [ItemSchema: Streaming Response Schemas](#itemschema-streaming-response-schemas)
3. [ItemEncoding: Per-Item Encoding Rules](#itemencoding-per-item-encoding-rules)
4. [Server-Side Implementation](#server-side-implementation)
5. [Client-Side Consumption](#client-side-consumption)
6. [Content Type Strategies](#content-type-strategies)
7. [Performance Optimization](#performance-optimization)
8. [Testing & Validation](#testing--validation)
9. [Troubleshooting](#troubleshooting)
10. [Real-World Examples](#real-world-examples)

---

## Overview

OpenAPI 3.2.0 introduces **itemSchema** and **itemEncoding** to properly document streaming responses where data arrives incrementally rather than as a single payload.

### Key Concepts

| Concept | Purpose | Use Case |
|---------|---------|----------|
| **itemSchema** | Schema for **each item** in stream | SSE, NDJSON, chunked responses |
| **itemEncoding** | Encoding rules **per item** | Custom serialization, compression |
| **schema** | Schema for **entire response** | Traditional single-object responses |

### When to Use

| Response Type | Use | OpenAPI Field |
|---------------|-----|---------------|
| Single object | `{"id": 1, "name": "..."}` | `schema` |
| Array | `[{...}, {...}]` | `schema` with `type: array` |
| **Stream** | `{...}\n{...}\n{...}` | **`itemSchema`** |
| **SSE** | `data: {...}\n\ndata: {...}` | **`itemSchema`** |

---

## ItemSchema: Streaming Response Schemas

### Basic Syntax

**nexs-swag annotation:**
```go
// @Success 200 {stream} ModelName "Description"
// @Produce text/event-stream
```

**Generated OpenAPI:**
```yaml
responses:
  '200':
    description: Description
    content:
      text/event-stream:
        itemSchema:
          $ref: '#/components/schemas/ModelName'
```

### Difference from schema

**Traditional (schema):**
```yaml
responses:
  '200':
    content:
      application/json:
        schema:  # Entire response body
          type: array
          items:
            $ref: '#/components/schemas/Event'
```
**Implication:** Client expects `[Event, Event, Event]` in **one response**.

**Streaming (itemSchema):**
```yaml
responses:
  '200':
    content:
      text/event-stream:
        itemSchema:  # Each streamed chunk
          $ref: '#/components/schemas/Event'
```
**Implication:** Client receives `Event`, `Event`, `Event` **over time**.

---

### Supported Content Types

| Content-Type | Description | itemSchema Support |
|--------------|-------------|-------------------|
| `text/event-stream` | Server-Sent Events (SSE) | ✅ Primary use case |
| `application/x-ndjson` | Newline Delimited JSON | ✅ Supported |
| `application/stream+json` | JSON streaming | ✅ Supported |
| `multipart/mixed` | Multipart streams | ✅ With itemEncoding |
| `application/octet-stream` | Binary streams | ⚠️ Use with encoding |

---

### Example: Server-Sent Events (SSE)

**Go Implementation:**
```go
// Event represents a real-time system event
type Event struct {
    ID        string    `json:"id" example:"evt_123"`
    Type      string    `json:"type" example:"user.login"`
    Timestamp time.Time `json:"timestamp"`
    Data      any       `json:"data"`
}

// StreamEvents sends real-time events to the client
// @Summary      Stream system events
// @Description  Receive real-time events via Server-Sent Events (SSE)
// @Tags         events
// @Produce      text/event-stream
// @Success      200 {stream} Event "Continuous stream of events"
// @Router       /events/stream [get]
func StreamEvents(c *gin.Context) {
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

    flusher, ok := c.Writer.(http.Flusher)
    if !ok {
        c.JSON(500, gin.H{"error": "Streaming unsupported"})
        return
    }

    eventChan := make(chan Event)
    defer close(eventChan)

    // Subscribe to event bus
    eventBus.Subscribe(eventChan)
    defer eventBus.Unsubscribe(eventChan)

    for {
        select {
        case event := <-eventChan:
            // SSE format: data: {json}\n\n
            data, _ := json.Marshal(event)
            fmt.Fprintf(c.Writer, "data: %s\n\n", data)
            flusher.Flush()

        case <-c.Request.Context().Done():
            // Client disconnected
            return
        }
    }
}
```

**Generated OpenAPI:**
```yaml
/events/stream:
  get:
    summary: Stream system events
    description: Receive real-time events via Server-Sent Events (SSE)
    tags: [events]
    responses:
      '200':
        description: Continuous stream of events
        content:
          text/event-stream:
            itemSchema:
              $ref: '#/components/schemas/Event'

components:
  schemas:
    Event:
      type: object
      properties:
        id:
          type: string
          example: evt_123
        type:
          type: string
          example: user.login
        timestamp:
          type: string
          format: date-time
        data: {}
```

**Client Implementation (TypeScript):**
```typescript
const eventSource = new EventSource('/events/stream');

eventSource.onmessage = (event: MessageEvent) => {
  const data: Event = JSON.parse(event.data);
  console.log('Received event:', data.type, data.id);
};

eventSource.onerror = () => {
  console.error('SSE connection error');
  eventSource.close();
};
```

---

### Example: Newline Delimited JSON (NDJSON)

**Go Implementation:**
```go
// LogEntry represents a log entry
type LogEntry struct {
    Level     string    `json:"level" example:"info"`
    Message   string    `json:"message" example:"Request processed"`
    Timestamp time.Time `json:"timestamp"`
    TraceID   string    `json:"trace_id" example:"abc123"`
}

// StreamLogs streams logs in NDJSON format
// @Summary      Stream application logs
// @Description  Stream logs as newline-delimited JSON
// @Tags         logs
// @Produce      application/x-ndjson
// @Success      200 {stream} LogEntry "Continuous log stream"
// @Router       /logs/stream [get]
func StreamLogs(c *gin.Context) {
    c.Header("Content-Type", "application/x-ndjson")
    c.Header("X-Content-Type-Options", "nosniff")

    encoder := json.NewEncoder(c.Writer)
    flusher := c.Writer.(http.Flusher)

    logChan := make(chan LogEntry)
    defer close(logChan)

    logger.Subscribe(logChan)
    defer logger.Unsubscribe(logChan)

    for {
        select {
        case entry := <-logChan:
            // NDJSON: one JSON object per line
            if err := encoder.Encode(entry); err != nil {
                return
            }
            flusher.Flush()

        case <-c.Request.Context().Done():
            return
        }
    }
}
```

**Client Implementation (Node.js):**
```javascript
const response = await fetch('/logs/stream');
const reader = response.body.getReader();
const decoder = new TextDecoder();
let buffer = '';

while (true) {
  const {done, value} = await reader.read();
  if (done) break;

  buffer += decoder.decode(value, {stream: true});
  const lines = buffer.split('\n');
  buffer = lines.pop(); // Keep incomplete line

  for (const line of lines) {
    if (line.trim()) {
      const entry = JSON.parse(line);
      console.log(`[${entry.level}] ${entry.message}`);
    }
  }
}
```

---

## ItemEncoding: Per-Item Encoding Rules

### What is itemEncoding?

While **itemSchema** defines the structure, **itemEncoding** defines **how each item is serialized** when transmitted.

### Syntax

**nexs-swag annotation:**
```go
// @Success 200 {stream} Event "Description"
// @Header 200 {string} Content-Type "application/x-ndjson"
// Note: itemEncoding configured in overrides file (see below)
```

**OpenAPI (manual configuration):**
```yaml
responses:
  '200':
    content:
      application/x-ndjson:
        itemSchema:
          $ref: '#/components/schemas/Event'
        itemEncoding:
          contentType:
            contentType: application/json
          headers:
            X-Item-ID:
              schema:
                type: string
```

### Use Cases

| Scenario | itemEncoding Configuration |
|----------|---------------------------|
| JSON items with headers | `headers: {X-Item-ID: {...}}` |
| Compressed items | `contentType: application/gzip` |
| Custom serialization | `contentType: application/x-custom` |
| Multipart boundaries | `headers: {Content-Disposition: {...}}` |

---

### Example: Compressed Stream

**Go Implementation:**
```go
// CompressedEvent represents an event that will be gzipped
type CompressedEvent struct {
    Data      []byte    `json:"data"`
    Timestamp time.Time `json:"timestamp"`
}

// StreamCompressedEvents streams gzip-compressed events
// @Summary      Stream compressed events
// @Description  Events compressed per-item with gzip
// @Tags         events
// @Produce      application/x-ndjson
// @Success      200 {stream} CompressedEvent "Gzip-compressed event stream"
// @Router       /events/compressed [get]
func StreamCompressedEvents(c *gin.Context) {
    c.Header("Content-Type", "application/x-ndjson")

    flusher := c.Writer.(http.Flusher)
    eventChan := make(chan CompressedEvent)

    for event := range eventChan {
        // Serialize to JSON
        data, _ := json.Marshal(event)

        // Compress with gzip
        var buf bytes.Buffer
        gzipWriter := gzip.NewWriter(&buf)
        gzipWriter.Write(data)
        gzipWriter.Close()

        // Send compressed chunk with metadata
        fmt.Fprintf(c.Writer, "Content-Encoding: gzip\r\n")
        fmt.Fprintf(c.Writer, "Content-Length: %d\r\n\r\n", buf.Len())
        c.Writer.Write(buf.Bytes())
        c.Writer.Write([]byte("\r\n"))
        flusher.Flush()
    }
}
```

**OpenAPI Configuration (.swaggo overrides):**
```yaml
paths:
  /events/compressed:
    get:
      responses:
        '200':
          content:
            application/x-ndjson:
              itemSchema:
                $ref: '#/components/schemas/CompressedEvent'
              itemEncoding:
                contentType:
                  contentType: application/gzip
                  headers:
                    Content-Encoding:
                      schema:
                        type: string
                        enum: [gzip]
```

---

### Example: Multipart Stream with Headers

**Go Implementation:**
```go
// FileChunk represents a chunk of a file being streamed
type FileChunk struct {
    PartNumber int    `json:"part_number"`
    Data       []byte `json:"data"`
    Checksum   string `json:"checksum"`
}

// StreamFileChunks streams file in multipart chunks
// @Summary      Stream file chunks
// @Description  Stream file in multipart/mixed format with per-chunk headers
// @Tags         files
// @Produce      multipart/mixed
// @Success      200 {stream} FileChunk "File chunk stream"
// @Router       /files/{id}/stream [get]
func StreamFileChunks(c *gin.Context) {
    boundary := "chunk-boundary-12345"
    c.Header("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", boundary))

    writer := multipart.NewWriter(c.Writer)
    writer.SetBoundary(boundary)
    defer writer.Close()

    flusher := c.Writer.(http.Flusher)

    for i, chunk := range getFileChunks(c.Param("id")) {
        // Create part with headers
        header := textproto.MIMEHeader{}
        header.Set("Content-Type", "application/json")
        header.Set("X-Part-Number", strconv.Itoa(i))
        header.Set("X-Checksum", chunk.Checksum)

        part, _ := writer.CreatePart(header)
        json.NewEncoder(part).Encode(chunk)
        flusher.Flush()
    }
}
```

**OpenAPI Configuration:**
```yaml
itemEncoding:
  headers:
    X-Part-Number:
      description: Sequential part number
      schema:
        type: integer
    X-Checksum:
      description: MD5 checksum of chunk
      schema:
        type: string
        pattern: '^[a-f0-9]{32}$'
```

---

## Server-Side Implementation

### Best Practices

#### 1. Always Set Proper Headers

```go
func setupStreamHeaders(c *gin.Context) {
    c.Header("Content-Type", "text/event-stream")
    c.Header("Cache-Control", "no-cache")
    c.Header("Connection", "keep-alive")
    c.Header("X-Accel-Buffering", "no")        // Nginx
    c.Header("X-Proxy-Buffering", "no")        // Generic proxy
    c.Header("Transfer-Encoding", "chunked")
}
```

#### 2. Handle Client Disconnections

```go
func StreamWithCleanup(c *gin.Context) {
    ctx, cancel := context.WithCancel(c.Request.Context())
    defer cancel()

    eventChan := make(chan Event)
    defer close(eventChan)

    go func() {
        <-ctx.Done()
        // Cleanup resources
        unsubscribe(eventChan)
    }()

    for {
        select {
        case event := <-eventChan:
            // Send event
        case <-ctx.Done():
            return
        }
    }
}
```

#### 3. Implement Backpressure

```go
func StreamWithBackpressure(c *gin.Context) {
    eventChan := make(chan Event, 100) // Buffered channel
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case event := <-eventChan:
            sendEvent(c, event)
        case <-ticker.C:
            // Send heartbeat to detect disconnections
            fmt.Fprintf(c.Writer, ": heartbeat\n\n")
            c.Writer.(http.Flusher).Flush()
        case <-c.Request.Context().Done():
            return
        }
    }
}
```

#### 4. Use Structured Logging

```go
func StreamWithLogging(c *gin.Context) {
    logger := log.With("client_id", c.ClientIP(), "stream", "events")
    logger.Info("Stream started")
    defer logger.Info("Stream ended")

    sentCount := 0
    for event := range eventChan {
        if err := sendEvent(c, event); err != nil {
            logger.Error("Send failed", "error", err, "sent_count", sentCount)
            return
        }
        sentCount++
    }
}
```

---

## Client-Side Consumption

### JavaScript/TypeScript

#### EventSource (SSE)

```typescript
interface Event {
  id: string;
  type: string;
  timestamp: string;
  data: any;
}

function connectEventStream(url: string): EventSource {
  const eventSource = new EventSource(url);

  eventSource.onopen = () => {
    console.log('Stream connected');
  };

  eventSource.onmessage = (event: MessageEvent) => {
    const data: Event = JSON.parse(event.data);
    handleEvent(data);
  };

  eventSource.onerror = (error) => {
    console.error('Stream error:', error);
    eventSource.close();
  };

  return eventSource;
}

// Auto-reconnect wrapper
class ReconnectingEventSource {
  private eventSource: EventSource | null = null;
  private reconnectDelay = 1000;
  private maxReconnectDelay = 30000;

  constructor(private url: string) {
    this.connect();
  }

  private connect() {
    this.eventSource = new EventSource(this.url);

    this.eventSource.onerror = () => {
      this.eventSource?.close();
      setTimeout(() => this.connect(), this.reconnectDelay);
      this.reconnectDelay = Math.min(
        this.reconnectDelay * 2,
        this.maxReconnectDelay
      );
    };

    this.eventSource.onopen = () => {
      this.reconnectDelay = 1000; // Reset on successful connect
    };
  }

  close() {
    this.eventSource?.close();
  }
}
```

#### Fetch API (NDJSON)

```typescript
async function consumeNDJSONStream(url: string) {
  const response = await fetch(url);
  if (!response.body) throw new Error('No response body');

  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = '';

  try {
    while (true) {
      const {done, value} = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, {stream: true});
      const lines = buffer.split('\n');
      buffer = lines.pop() || ''; // Keep incomplete line

      for (const line of lines) {
        if (line.trim()) {
          const item = JSON.parse(line);
          await processItem(item);
        }
      }
    }
  } finally {
    reader.releaseLock();
  }
}
```

### Python

#### SSE Client

```python
import requests
import json
from typing import Generator, Dict, Any

def consume_sse_stream(url: str) -> Generator[Dict[str, Any], None, None]:
    """Consume Server-Sent Events stream"""
    with requests.get(url, stream=True) as response:
        response.raise_for_status()
        
        for line in response.iter_lines(decode_unicode=True):
            if line.startswith('data:'):
                data = line[5:].strip()  # Remove 'data:' prefix
                if data:
                    yield json.loads(data)

# Usage
for event in consume_sse_stream('http://localhost:8080/events/stream'):
    print(f"Event: {event['type']} at {event['timestamp']}")
```

#### NDJSON Client

```python
def consume_ndjson_stream(url: str) -> Generator[Dict[str, Any], None, None]:
    """Consume Newline Delimited JSON stream"""
    with requests.get(url, stream=True) as response:
        response.raise_for_status()
        
        for line in response.iter_lines(decode_unicode=True):
            if line:
                yield json.loads(line)

# With error handling and retries
from tenacity import retry, stop_after_attempt, wait_exponential

@retry(stop=stop_after_attempt(3), wait=wait_exponential(multiplier=1))
def consume_stream_with_retry(url: str):
    for item in consume_ndjson_stream(url):
        process_item(item)
```

### Go Client

```go
// SSE Client
func ConsumeSSEStream(url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    reader := bufio.NewReader(resp.Body)
    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return err
        }

        if strings.HasPrefix(line, "data:") {
            data := strings.TrimPrefix(line, "data:")
            data = strings.TrimSpace(data)
            
            var event Event
            if err := json.Unmarshal([]byte(data), &event); err != nil {
                log.Printf("Parse error: %v", err)
                continue
            }
            
            handleEvent(event)
        }
    }
}

// NDJSON Client
func ConsumeNDJSONStream(url string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    decoder := json.NewDecoder(resp.Body)
    for decoder.More() {
        var item LogEntry
        if err := decoder.Decode(&item); err != nil {
            return err
        }
        processItem(item)
    }
    return nil
}
```

---

## Content Type Strategies

### Choosing the Right Content Type

| Content-Type | Best For | Pros | Cons |
|--------------|----------|------|------|
| `text/event-stream` | Web browsers, real-time UI | Native browser support, reconnect | SSE format overhead |
| `application/x-ndjson` | Log streaming, data pipelines | Simple parsing, efficient | No browser API |
| `application/stream+json` | API-to-API streaming | Clean format | Less tooling support |
| `multipart/mixed` | File uploads, mixed content | Headers per part | Complex parsing |

### Implementation Matrix

```go
func StreamHandler(c *gin.Context) {
    accept := c.GetHeader("Accept")
    
    switch {
    case strings.Contains(accept, "text/event-stream"):
        streamSSE(c)
    case strings.Contains(accept, "application/x-ndjson"):
        streamNDJSON(c)
    case strings.Contains(accept, "application/stream+json"):
        streamJSON(c)
    default:
        c.JSON(400, gin.H{"error": "Unsupported content type"})
    }
}
```

---

## Performance Optimization

### 1. Connection Pooling

```go
var eventPool = sync.Pool{
    New: func() interface{} {
        return new(Event)
    },
}

func StreamOptimized(c *gin.Context) {
    for {
        event := eventPool.Get().(*Event)
        populateEvent(event)
        
        sendEvent(c, event)
        
        // Reset and return to pool
        *event = Event{}
        eventPool.Put(event)
    }
}
```

### 2. Batch Flushing

```go
func StreamBatched(c *gin.Context) {
    buffer := make([]Event, 0, 10)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()

    for {
        select {
        case event := <-eventChan:
            buffer = append(buffer, event)
            if len(buffer) >= 10 {
                sendBatch(c, buffer)
                buffer = buffer[:0]
            }
        case <-ticker.C:
            if len(buffer) > 0 {
                sendBatch(c, buffer)
                buffer = buffer[:0]
            }
        }
    }
}
```

### 3. Compression

```go
func StreamCompressed(c *gin.Context) {
    c.Header("Content-Encoding", "gzip")
    
    gzipWriter := gzip.NewWriter(c.Writer)
    defer gzipWriter.Close()
    
    encoder := json.NewEncoder(gzipWriter)
    for event := range eventChan {
        encoder.Encode(event)
        gzipWriter.Flush()
        c.Writer.(http.Flusher).Flush()
    }
}
```

---

## Testing & Validation

### Unit Tests

```go
func TestStreamResponse(t *testing.T) {
    router := gin.New()
    router.GET("/stream", StreamEvents)

    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/stream", nil)
    
    done := make(chan bool)
    go func() {
        router.ServeHTTP(w, req)
        done <- true
    }()

    // Send test events
    eventChan <- Event{ID: "1", Type: "test"}
    
    time.Sleep(100 * time.Millisecond)
    
    assert.Contains(t, w.Body.String(), `"id":"1"`)
    assert.Contains(t, w.Body.String(), `"type":"test"`)
}
```

### Integration Tests

```bash
# Test SSE stream with curl
curl -N http://localhost:8080/events/stream

# Test with timeout
timeout 5 curl -N http://localhost:8080/events/stream | head -10

# Validate JSON structure per event
curl -N http://localhost:8080/logs/stream | \
  while IFS= read -r line; do
    echo "$line" | jq -e . > /dev/null || echo "Invalid JSON: $line"
  done
```

---

## Troubleshooting

### Problem: Stream Buffering

**Symptoms:** Events arrive in batches instead of real-time

**Solutions:**
```go
// 1. Disable proxy buffering
c.Header("X-Accel-Buffering", "no")

// 2. nginx configuration
location /stream {
    proxy_buffering off;
    proxy_cache off;
    proxy_set_header Connection '';
    chunked_transfer_encoding on;
}

// 3. Use explicit flush
flusher, _ := c.Writer.(http.Flusher)
flusher.Flush()
```

---

### Problem: Client Disconnection Not Detected

**Solution:**
```go
func StreamWithHeartbeat(c *gin.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case event := <-eventChan:
            if err := sendEvent(c, event); err != nil {
                return // Connection closed
            }
        case <-ticker.C:
            // SSE comment = heartbeat
            if _, err := fmt.Fprintf(c.Writer, ": ping\n\n"); err != nil {
                return
            }
            c.Writer.(http.Flusher).Flush()
        }
    }
}
```

---

### Problem: Memory Leak from Unclosed Streams

**Solution:**
```go
type StreamManager struct {
    mu      sync.RWMutex
    streams map[string]chan Event
}

func (sm *StreamManager) Register(id string) chan Event {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    ch := make(chan Event, 100)
    sm.streams[id] = ch
    return ch
}

func (sm *StreamManager) Unregister(id string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    if ch, ok := sm.streams[id]; ok {
        close(ch)
        delete(sm.streams, id)
    }
}

func StreamWithCleanup(c *gin.Context) {
    clientID := uuid.New().String()
    eventChan := streamManager.Register(clientID)
    defer streamManager.Unregister(clientID)
    
    // ... streaming logic
}
```

---

## Real-World Examples

### Example 1: Real-Time Analytics Dashboard

```go
type MetricUpdate struct {
    Metric    string    `json:"metric"`
    Value     float64   `json:"value"`
    Timestamp time.Time `json:"timestamp"`
    Tags      map[string]string `json:"tags"`
}

// @Summary Stream real-time metrics
// @Produce text/event-stream
// @Success 200 {stream} MetricUpdate
// @Router /metrics/stream [get]
func StreamMetrics(c *gin.Context) {
    setupStreamHeaders(c)
    
    metricsChan := metrics.Subscribe()
    defer metrics.Unsubscribe(metricsChan)
    
    for metric := range metricsChan {
        data, _ := json.Marshal(metric)
        fmt.Fprintf(c.Writer, "event: metric\ndata: %s\n\n", data)
        c.Writer.(http.Flusher).Flush()
    }
}
```

### Example 2: Progressive AI Response

```go
type AIChunk struct {
    Delta     string `json:"delta"`      // Text chunk
    TokenID   int    `json:"token_id"`
    IsLast    bool   `json:"is_last"`
}

// @Summary Stream AI completion
// @Produce application/x-ndjson
// @Success 200 {stream} AIChunk
// @Router /ai/complete [post]
func StreamAICompletion(c *gin.Context) {
    c.Header("Content-Type", "application/x-ndjson")
    
    encoder := json.NewEncoder(c.Writer)
    flusher := c.Writer.(http.Flusher)
    
    for chunk := range aiModel.StreamCompletion(prompt) {
        encoder.Encode(chunk)
        flusher.Flush()
    }
}
```

---

## Additional Resources

- [OpenAPI 3.2.0 Specification](https://spec.openapis.org/oas/v3.2.0.html)
- [Server-Sent Events Standard](https://html.spec.whatwg.org/multipage/server-sent-events.html)
- [NDJSON Specification](http://ndjson.org/)
- [HTTP Chunked Transfer Encoding (RFC 7230)](https://datatracker.ietf.org/doc/html/rfc7230#section-4.1)

---

**Questions?** Open an issue at [github.com/fsvxavier/nexs-swag/issues](https://github.com/fsvxavier/nexs-swag/issues)
