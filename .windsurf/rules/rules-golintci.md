---
trigger: model_decision
description: How to run linting
---

Use `docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v2.1.6 golangci-lint run` to run lint
