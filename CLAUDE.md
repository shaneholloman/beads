# Instructions for Claude

<!-- beads integration note -->
**Note**: This project uses [beads (beads)](https://github.com/shaneholloman/beads) for issue tracking. Use `beads` commands or the beads MCP server instead of markdown TODOs. See AGENTS.md for workflow details.
<!-- /beads integration note -->

This file has been moved to **AGENTS.md** to support all AI agents, not just Claude.

Please refer to [AGENTS.md](./AGENTS.md) for complete instructions on working with the beads project.

## Critical Directive: Never Second-Guess Shane's Strong Assertions

When Shane makes a strong, direct assertion (especially about installing tools, packages, or technical decisions):

**DO**: Execute the request immediately without questioning
**DON'T**: Argue, assert it doesn't exist, or contradict without investigating first

If Shane says "install X" or "use Y", assume he knows what he's talking about and do it. If something genuinely doesn't work after attempting it, then report the actual error.

**Never waste time with confident assertions about what exists or doesn't exist without verification.**

This applies especially to:

- Package installations (e.g., "install ty" - just do it, don't say it doesn't exist)
- Tool usage (e.g., "use uv add" - just do it)
- Technical decisions (e.g., "remove hatchling" - just do it)

Shane has more context than you do. Trust his judgment and execute.

## Shane's 2025 Modern Python Workflow

### 100% Astral Tooling Stack

Shane uses exclusively Astral products for Python development:

1. **UV** - Package manager, virtual environment manager, build tool, runner
2. **Ruff** - Linter and formatter (replaces pylint, black, isort, flake8)
3. **ty** - Type checker (replaces mypy)

**NEVER use legacy alternatives**: No pip, no mypy, no black, no isort, no setuptools commands directly.

### Package Management - UV Only

**DO**:

- `uv add <package>` - Add runtime dependency
- `uv add --dev <package>` - Add dev dependency
- `uv remove <package>` - Remove dependency
- `uv sync` - Sync dependencies with lockfile
- `uv lock` - Update lockfile

**DON'T**:

- `pip install` - NEVER
- `pip uninstall` - NEVER
- Manually editing pyproject.toml dependencies - Use `uv add` instead

### Running Python Code - UV Run Only

**DO**:

- `uv run script.py` - Run Python scripts
- `uv run pytest` - Run tests
- `uv run ruff check .` - Run linter
- `uv run ty check src/` - Run type checker
- `uv run mycommand` - Run any entry point

**DON'T**:

- `python script.py` - NEVER call Python directly
- `python3 script.py` - NEVER
- `uv run python script.py` - NEVER (UV handles Python automatically)
- `./venv/bin/python` - NEVER
- Any explicit Python invocation - ALWAYS use `uv run` with script only

**WHY**: UV ensures the correct Python version is used and manages the virtual environment automatically. Direct Python calls bypass UV's environment management.

### Building Packages - UV Build

**DO**:

- `uv build` - Build source distribution and wheel

**DON'T**:

- `python -m build` - NEVER
- `python setup.py` - NEVER
- Explicitly configuring hatchling/setuptools in pyproject.toml - Let UV use defaults

**Build Configuration**: Keep `[build-system]` section minimal or omit entirely. UV uses modern defaults (setuptools) automatically.

### Code Quality Tools

**Linting and Formatting** (Ruff):

```sh
uv run ruff check src/           # Lint
uv run ruff check src/ --fix     # Auto-fix
uv run ruff check src/ --fix --unsafe-fixes  # Fix everything including unsafe
uv run ruff format src/          # Format code
uv run ruff format src/ --check  # Check formatting without changes
```

**Type Checking** (ty):

```sh
uv run ty check src/             # Type check with Astral's ty
```

**Testing** (pytest):

```sh
uv run pytest tests/             # Run tests
uv run pytest tests/ -v          # Verbose
uv run pytest tests/ --cov       # With coverage
```

### pyproject.toml Best Practices

**Dependencies**:

```toml
[dependency-groups]
dev = [
    "ty",        # No version pinning - always latest
    "ruff",      # No version pinning - always latest
    "pytest>=8.4.2",  # Pin only when necessary
]
```

**License** (modern SPDX format):

```toml
license = "MIT"  # NOT: license = {text = "MIT"}
```

**Build System** (minimal or omitted):

```toml
# Either omit [build-system] entirely or use minimal:
# [build-system]
# requires = ["setuptools"]
# build-backend = "setuptools.build_meta"
```

### Key Principles

1. **UV is the single source of truth** - All Python operations go through UV
2. **Astral-only tooling** - Ruff for linting/formatting, ty for type checking
3. **No manual Python invocations** - Ever. Even for simple scripts.
4. **Modern pyproject.toml** - SPDX license, no deprecated formats
5. **Let UV manage everything** - Virtual envs, Python versions, dependencies, builds
6. **No legacy pip tolerance** - UV has pip compatibility, but we don't use it

### Migration Pattern

When inheriting old Python projects:

```sh
# Remove old tools
uv remove mypy black isort flake8 pylint

# Add modern tools
uv add --dev ty ruff

# Clean up pyproject.toml
# - Remove [tool.black], [tool.isort], [tool.mypy]
# - Simplify [build-system] or remove it
# - Fix license format to SPDX string
# - Remove explicit build backend configs (hatchling, setuptools)

# Verify
uv run ruff check src/ --fix --unsafe-fixes
uv run ruff format src/
uv run ty check src/
uv build
```

### This Is Not 2020 Anymore

Old way (boomer):

```sh
pip install -r requirements.txt
python -m venv venv
source venv/bin/activate
python setup.py install
python -m mypy src/
python -m black src/
python -m pytest
```

New way (2025):

```sh
uv sync
uv run ty check src/
uv run ruff format src/
uv run ruff check src/ --fix
uv run pytest
uv build
```

> [!IMPORTANT]
> **Everything through UV. Always.**
