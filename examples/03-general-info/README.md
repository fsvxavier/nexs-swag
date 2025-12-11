# Example 03 - General Info File

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates using `--generalInfo` to specify which file contains the API's general annotations.

## Problem

When you have multiple Go files, the parser may find general info annotations (@title, @version) in several places, causing conflicts.

## Solution

Use `--generalInfo` to specify exactly which file contains the general info:

```bash
nexs-swag init --generalInfo main.go
```

## Structure

```
03-general-info/
├── main.go       # ✅ HAS @title, @version, @host, etc
├── products.go   # ❌ Only product endpoints
├── orders.go     # ❌ Only order endpoints
└── run.sh
```

## Rule

- **General Info File:** Must have @title, @version, @host, @BasePath
- **Other Files:** Should have ONLY endpoints (@Router, @Summary, etc)

## How to Run

```bash
chmod +x run.sh
./run.sh
```

## Comparison

### Without --generalInfo
```bash
nexs-swag init --dir .
# May generate error if @title found in multiple files
```

### With --generalInfo
```bash
nexs-swag init --dir . --generalInfo main.go
# ✅ Correct: only main.go is parsed for general info
# ✅ products.go and orders.go provide only endpoints
```

## Benefits

1. **Avoids conflicts:** Single location for API info
2. **Faster:** Parser doesn't need to check all files for general info
3. **Organization:** Separates concerns (general info vs endpoints)
