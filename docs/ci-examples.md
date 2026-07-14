# CI/CD Integration Examples

## GitHub Actions

### Basic Mission Archive

```yaml
name: Archive Mission

on:
  push:
    branches: [main]

jobs:
  archive:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      
      - name: Build Spacecraft
        run: make build
      
      - name: Archive Mission
        run: |
          ./scripts/spacecraft archive --ci > archive-result.json
        env:
          SPACECRAFT_MISSION: ${{ github.event.head_commit.message }}
      
      - name: Upload Archive Result
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: archive-result
          path: archive-result.json
```

### Deploy with Hooks

```yaml
name: Deploy with Hooks

on:
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Spacecraft
        run: make build
      
      - name: Pre-deploy Hook
        run: |
          echo "Running pre-deploy checks..."
          # Add your pre-deploy logic here
      
      - name: Deploy
        run: |
          ./scripts/spacecraft archive --ci
        env:
          SPACECRAFT_MISSION: ${{ github.event.inputs.mission_id }}
      
      - name: Post-deploy Hook
        if: success()
        run: |
          echo "Deployment successful!"
          # Add your post-deploy logic here (notifications, monitoring, etc.)
```

### CI-Friendly Output Parsing

```yaml
name: CI Pipeline

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build
        run: make build
      
      - name: Run Tests
        run: make test
      
      - name: Archive with JSON Output
        id: archive
        run: |
          ./scripts/spacecraft archive --ci > result.json
          echo "success=$(jq -r .success result.json)" >> $GITHUB_OUTPUT
          echo "mission_id=$(jq -r .missionId result.json)" >> $GITHUB_OUTPUT
      
      - name: Check Result
        if: steps.archive.outputs.success == 'true'
        run: |
          echo "Mission ${{ steps.archive.outputs.mission_id }} archived successfully"
```

## GitLab CI

```yaml
stages:
  - test
  - deploy

test:
  stage: test
  script:
    - make build
    - make test

deploy:
  stage: deploy
  script:
    - make build
    - ./scripts/spacecraft archive --ci > archive.json
    - cat archive.json
  only:
    - main
  artifacts:
    paths:
      - archive.json
```

## Hook Configuration

Create `.space/hooks.json` to define deploy hooks:

```json
{
  "hooks": [
    {
      "event": "deploy.before",
      "label": "pre-deploy-check",
      "command": "echo 'Running pre-deploy checks...' && ./scripts/validate.sh",
      "blocking": true,
      "timeout": 60
    },
    {
      "event": "deploy.after",
      "label": "post-deploy-notify",
      "command": "echo 'Deployment complete!' && ./scripts/notify.sh",
      "blocking": false,
      "timeout": 30
    }
  ]
}
```

## Environment Variables

- `SPACECRAFT_MISSION`: Override mission selector
- `SPACECRAFT_ROOT`: Override project root
- `SPACECRAFT_SPACE_DIR`: Override .space directory location

## JSON Output Schema

The `--ci` flag produces JSON output with this schema:

```json
{
  "success": true,
  "missionId": "M07ABC123",
  "archiveDir": ".space/archive/M07ABC123"
}
```

On failure:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```
