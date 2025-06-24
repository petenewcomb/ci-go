# go-licenses-ignore

A utility for managing ignore patterns for the [go-licenses](https://github.com/google/go-licenses) tool. Automatically resolves module names and supports sophisticated pattern matching for multi-module repositories.

## Installation

```bash
go install github.com/petenewcomb/ci-go/cmd/go-licenses-ignore@latest
```

## Usage

```bash
go-licenses-ignore <module-path>
```

The utility reads a `.go-licenses-ignore` configuration file from the repository root and outputs appropriate `--ignore` and `--disallowed_types` arguments for `go-licenses`.

Example usage in a script:
```bash
#!/bin/bash
MODULE_PATH="./internal/myservice"
ARGS=$(go-licenses-ignore "$MODULE_PATH")
go-licenses check $MODULE_PATH $ARGS
```

## Configuration File Format

Create a `.go-licenses-ignore` file in your repository root with the following syntax:

### Disallowed License Types

Lines starting with `!` specify license types that should be disallowed:

```
!forbidden
!restricted
!unknown
```

### Global Ignores

Package names without a prefix apply to all modules:

```
example.com/vendor/package
github.com/legacy/dependency
```

### Module-Specific Ignores

Use the `prefix:package` syntax to apply ignores only to modules matching the prefix:

```
# Single-line format
internal/services/*:github.com/problematic/dep
**/cmd/*:vendor/legacy-tool

# Multi-line format
internal/services/*:
    github.com/problematic/dep
    vendor/legacy-tool
    some.company.com/internal/utils
```

### Local Package References

Special syntax for referencing packages within the repository:

- **`.`** - Self-reference (ignore the module itself)
- **`./path`** - Relative to the current module's directory  
- **`/path`** - Relative to the repository root

Examples:
```
# Ignore all command modules themselves
**/cmd/*:.

# Ignore sibling packages relative to current module
myservice:./testutils

# Ignore packages relative to repository root
**/internal/*:/shared/deprecated
```

### Glob Patterns

Both prefixes and local package references support glob patterns:

```
# Module prefixes support globs
**/cmd/*:.
internal/services/**:.

# Local paths support globs (only directories with Go files are included)
**/cmd/*:/tools/*
internal/services/*:./*/testdata
```

### Comments

Lines starting with `#` are comments. Inline comments are also supported:

```
# This is a comment
internal/cmd/chartgen:
    golang.org/x/image/vector    # BSD-3-Clause
    github.com/ajstarks/svgo     # CC-BY-4.0
```

## Complete Example

```
# Disallowed license types
!forbidden
!restricted
!unknown

# Global ignores for all modules
legacy.company.com/deprecated

# Module-specific ignores
**/cmd/*:.                    # All commands ignore themselves
internal/services/*:./shared  # Services ignore their shared directory

# Multi-line format for complex cases
internal/cmd/chartgen:
    golang.org/x/image/vector          # BSD-3-Clause
    github.com/golang/freetype/raster  # BSD-ish FreeType License
    github.com/ajstarks/svgo           # CC-BY-4.0

# Glob expansion for repository-wide patterns
**/test/*:/tools/testgen/*
```

## Integration with CI/CD

### Pre-commit Hook

Add to your `.pre-commit-config.yaml`:

```yaml
repos:
- repo: local
  hooks:
  - id: go-licenses
    name: Check Go Licenses
    entry: bash -c 'for mod in $(find . -name go.mod -exec dirname {} \;); do go-licenses check $mod $(go-licenses-ignore $mod); done'
    language: system
    files: \.go$
```

### GitHub Actions

```yaml
- name: Check licenses
  run: |
    for module in $(find . -name go.mod -exec dirname {} \;); do
      echo "Checking licenses for $module"
      go-licenses check $module $(go-licenses-ignore $module)
    done
```

## How It Works

1. **Repository Detection**: Automatically finds the git repository root
2. **Module Resolution**: Determines the current module's path relative to the repository
3. **Pattern Matching**: Uses doublestar glob patterns to match module prefixes
4. **Local Path Resolution**: Resolves local package references to full module names
5. **Go Package Detection**: Only includes directories that contain Go source files
6. **Deduplication**: Sorts and removes duplicate package names from output

## Output Format

The utility outputs space-separated arguments ready for use with `go-licenses`:

```bash
--disallowed_types forbidden,restricted,unknown --ignore pkg1,pkg2,pkg3
```

This format makes it easy to integrate into shell scripts and CI/CD pipelines.

## Multi-Module Repository Support

The utility is designed for multi-module repositories where different modules may need different ignore patterns. The prefix matching system allows you to apply specific ignores to groups of modules while maintaining a single configuration file.

## License

This utility is part of the ci-go project and follows the same license terms.