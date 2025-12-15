# Example 23: Recursive Parsing

## Overview

Demonstrates **recursive parsing** with `--parseInternal` and `--exclude` flags working together with `--parseDependency`.

## Project Structure

```
23-recursive-parsing/
├── main.go                    # Main package with @title and routes
├── internal/
│   ├── handlers/
│   │   └── user.go           # Internal user handlers (POST, PUT)
│   └── models/
│       └── product.go        # Internal product models (GET, LIST)
├── pkg/
│   └── utils/
│       └── helpers.go        # Public utilities (GET)
├── config/
│   └── settings.go           # Config endpoints (excluded with --exclude)
└── run.sh                    # Test script
```

## Key Flags

### `--parseInternal`
- **Default:** `false`
- **Purpose:** Parse directories named `internal/` recursively
- **Impact:** Includes all routes defined in `internal/handlers/` and `internal/models/`

### `--exclude <pattern>`
- **Default:** `""`
- **Purpose:** Exclude directories/files from parsing
- **Patterns:** Supports exact names, paths, wildcards (`*`)
- **Examples:** 
  - `--exclude config` (excludes `config/` directory)
  - `--exclude ./config` (same, with relative path)
  - `--exclude vendor,config` (multiple patterns)

### `--parseDependency` (`--pd`)
- **Default:** `false`
- **Purpose:** Parse external dependencies in `vendor/` or `go.mod`
- **Works with:** `--parseDependencyLevel` (`--pdl`)

## Tests

### Test 1: Without `--parseInternal`

```bash
nexs-swag init --output ./docs-no-internal --ov 3.1
```

**Expected endpoints:**
- `/users/{id}` (from `main.go`)
- `/pkg/helpers` (from `pkg/utils/helpers.go`)
- `/health` (from `pkg/utils/helpers.go`)
- `/config/settings` (from `config/settings.go`)
- `/config/database` (from `config/settings.go`)

**NOT included:**
- `/internal/users` ❌
- `/internal/users/{id}` ❌
- `/internal/products` ❌
- `/internal/products/{id}` ❌

---

### Test 2: With `--parseInternal`

```bash
nexs-swag init --output ./docs-with-internal --ov 3.1 --parseInternal
```

**Expected endpoints (all):**
- `/users/{id}` ✅
- `/internal/users` ✅
- `/internal/users/{id}` ✅
- `/internal/products` ✅
- `/internal/products/{id}` ✅
- `/pkg/helpers` ✅
- `/health` ✅
- `/config/settings` ✅
- `/config/database` ✅

---

### Test 3: With `--parseInternal` + `--exclude config`

```bash
nexs-swag init --output ./docs-exclude-config --ov 3.1 --parseInternal --exclude config
```

**Expected endpoints:**
- `/users/{id}` ✅
- `/internal/users` ✅
- `/internal/users/{id}` ✅
- `/internal/products` ✅
- `/internal/products/{id}` ✅
- `/pkg/helpers` ✅
- `/health` ✅

**NOT included (excluded):**
- `/config/settings` ❌
- `/config/database` ❌

---

### Test 4: Full Command (Corrected Syntax)

```bash
nexs-swag init \
  --output ./docs \
  --ov 3.1 \
  --pd \
  --pdl 3 \
  --parseInternal \
  --validate \
  --exclude config
```

**Flags explained:**
- `--ov 3.1` → OpenAPI version 3.1
- `--pd` → Parse dependencies (external packages) - boolean flag, NO value needed
- `--pdl 3` → Parse dependency level 3 (all: models + operations + schemas)
- `--parseInternal` → Include `internal/` directories - boolean flag, NO value needed
- `--validate` → Validate OpenAPI spec after generation
- `--exclude config` → Exclude `config/` directory

**Expected result:** Same as Test 3 (all internal routes, no config routes)

---

## ⚠️ Common Syntax Errors

### ❌ WRONG Syntax

```bash
# DO NOT use "true" with boolean flags
nexs-swag init --parseInternal true --pd true   ❌

# DO NOT use "./" prefix unnecessarily  
nexs-swag init --exclude ./config               ❌
```

**Why it fails:**
- `--parseInternal true` → CLI interprets "true" as the NEXT argument (like a file), not as the flag's value
- `--pd true` → Same issue, "true" becomes a positional argument
- Boolean flags in CLI libraries (like urfave/cli) are **presence-based**, not value-based

### ✅ CORRECT Syntax

```bash
# Boolean flags without values
nexs-swag init --parseInternal --pd             ✅

# Exclude can use ./ but it's optional
nexs-swag init --exclude config                 ✅
nexs-swag init --exclude ./config               ✅ (both work)
```

## Running the Example

```bash
cd examples/23-recursive-parsing
chmod +x run.sh
./run.sh
```

## Key Learnings

1. **`--parseInternal` enables recursive parsing of `internal/` directories**
   - Without it, `internal/` is completely skipped
   - With it, all subdirectories are parsed recursively

2. **`--exclude` works independently of `--parseInternal`**
   - You can parse `internal/` while excluding other directories
   - Patterns match directory names, relative paths, and wildcards

3. **`--parseDependency` is for external dependencies, not project structure**
   - Use `--pd` + `--pdl` to parse vendor packages or go.mod dependencies
   - Does NOT affect parsing of your own project directories

4. **All flags can be combined**
   - `--parseInternal` + `--exclude` → Parse internal, exclude specific dirs
   - `--parseInternal` + `--pd` → Parse both internal and dependencies
   - All three together work seamlessly

## Common Issues

### Issue: "Internal routes not showing up"
**Solution:** Add `--parseInternal` flag

### Issue: "Unwanted config routes in docs"
**Solution:** Use `--exclude config` or `--exclude ./config`

### Issue: "External package types not resolved"
**Solution:** Add `--pd true --pdl 1` (or higher)

## Best Practices

1. **Always use `--parseInternal` for microservices**
   - Most Go projects use `internal/` for implementation details
   - You need this flag to document internal APIs

2. **Exclude sensitive or non-API directories**
   - `--exclude config,scripts,migrations`
   - Keeps documentation focused on actual API routes

3. **Use `--pdl 3` only when necessary**
   - Level 1 (models) is usually sufficient
   - Level 3 can slow down parsing significantly

4. **Combine with `--validate` in CI/CD**
   - Ensures generated OpenAPI spec is valid
   - Catches documentation errors early

## Related Examples

- [08-parse-internal](../08-parse-internal) - Basic `--parseInternal` usage
- [09-parse-dependency](../09-parse-dependency) - Basic `--parseDependency` usage
- [10-dependency-level](../10-dependency-level) - `--parseDependencyLevel` details
- [06-exclude](../06-exclude) - `--exclude` patterns and wildcards
