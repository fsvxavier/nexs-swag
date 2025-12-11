# Example 06 - Exclude Patterns

🌍 **English** • [Português (Brasil)](README_pt.md) • [Español](README_es.md)

Demonstrates how to exclude directories and files from parsing.

## Flag

```bash
--exclude pattern1,pattern2,pattern3
```

## Usage

```bash
# Exclude one directory
nexs-swag init --exclude mock

# Exclude multiple
nexs-swag init --exclude mock,testdata,vendor

# Exclude with wildcards
nexs-swag init --exclude "*.test.go,*_mock.go"
```

## Automatic Exclusions

Always excluded (no need to specify):
- `vendor/` - Dependencies
- `testdata/` - Test data
- `docs/` - Generated documentation
- `.git/` - Git repository
- `*_test.go` - Test files

## Example Structure

```
06-exclude/
├── main.go           # ✅ Will be parsed
├── main_test.go      # ❌ Excluded (test)
├── mock/
│   └── mock.go       # ❌ Excluded (with flag)
└── testdata/
    └── data.go       # ❌ Excluded (automatic)
```

## How to Run

```bash
./run.sh
```

## Use Cases

- **mock:** Mocking code for tests
- **testdata:** Fixtures and test data
- **vendor:** Dependencies (if using vendor)
- **examples:** Example code
- **internal:** Internal packages (use --parseInternal to include)
