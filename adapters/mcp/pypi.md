# Publishing mcp-beads to PyPI

This guide covers how to build and publish the mcp-beads package to the Python Package Index (PyPI).

## Prerequisites

1. **PyPI Account**: Create accounts on both:
   - Test PyPI: <https://test.pypi.org/account/register/>
   - PyPI: <https://pypi.org/account/register/>

2. **API Tokens**: Generate API tokens for authentication:
   - Test PyPI: <https://test.pypi.org/manage/account/token/>
   - PyPI: <https://pypi.org/manage/account/token/>

3. **Build Tools**: UV handles building automatically (no separate tools needed)

## Building the Package

1. **Clean previous builds** (if any):

   ```sh
   rm -rf dist/ build/ src/*.egg-info
   ```

2. **Build the distribution packages**:

   ```sh
   uv build
   ```

   This creates both:
   - `dist/mcp_beads-0.30.0-py3-none-any.whl` (wheel)
   - `dist/mcp_beads-0.30.0.tar.gz` (source distribution)

3. **Verify the build**:

   ```sh
   tar -tzf dist/mcp_beads-0.30.0.tar.gz
   ```

   Should include:
   - Source files in `src/mcp_beads/`
   - `README.md`
   - `LICENSE`
   - `pyproject.toml`

## Testing the Package

### Test on Test PyPI First

1. **Upload to Test PyPI**:

   ```sh
   uv run twine upload --repository testpypi dist/*
   ```

   When prompted, use:
   - Username: `__token__`
   - Password: Your Test PyPI API token (including the `pypi-` prefix)

2. **Install from Test PyPI**:

   ```sh
   # Test installation
   uv tool install --index-url https://test.pypi.org/simple/ mcp-beads

   # Test it works
   mcp-beads --help
   ```

3. **Verify the installation**:

   ```sh
   uv run python -c "import mcp_beads; print(mcp_beads.__version__)"
   ```

## Publishing to PyPI

Once you've verified the package works on Test PyPI:

1. **Upload to PyPI**:

   ```sh
   uv run twine upload dist/*
   ```

   Use:
   - Username: `__token__`
   - Password: Your PyPI API token

2. **Verify on PyPI**:
   - Visit <https://pypi.org/project/mcp-beads/>
   - Check that the README displays correctly
   - Verify all metadata is correct

3. **Test installation**:

   ```sh
   uv tool install mcp-beads
   mcp-beads --help
   ```

## Updating the README Installation Instructions

After publishing, users can install simply with:

```sh
uv tool install mcp-beads
```

Update the README.md to reflect this simpler installation method.

## Version Management

When releasing a new version:

1. Use the version bump script from the parent project (updates all files automatically):

   ```sh
   cd ../..
   ./scripts/bump-version.sh 0.9.5 --commit
   ```

4. Create a git tag:

   ```sh
   git tag v0.9.5
   git push origin v0.9.5
   ```

5. Clean, rebuild, and republish to PyPI

## Troubleshooting

### Package Already Exists

PyPI doesn't allow re-uploading the same version. If you need to fix something:

1. Increment the version number (even for minor fixes)
2. Rebuild and re-upload

### Missing Files in Distribution

If files are missing from the built package, create a `MANIFEST.in`:

```
include README.md
include LICENSE
recursive-include src/mcp_beads *.py
```

### Authentication Errors

- Ensure you're using `__token__` as the username (exactly)
- Token should include the `pypi-` prefix
- Check token hasn't expired

### Test PyPI vs Production

Test PyPI is completely separate from production PyPI:

- Different accounts
- Different tokens
- Different package versions (can have different versions on each)

Always test on Test PyPI first!

## Continuous Deployment (Future)

Consider setting up GitHub Actions to automate this:

1. On tag push (e.g., `v0.9.5`)
2. Run tests
3. Build package
4. Publish to PyPI

See `.github/workflows/` in the parent project for examples.

## Resources

- [Python Packaging Guide](https://packaging.python.org/tutorials/packaging-projects/)
- [PyPI Documentation](https://pypi.org/help/)
- [Twine Documentation](https://twine.readthedocs.io/)
