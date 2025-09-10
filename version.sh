#!/bin/bash

TAG=$(git describe --tags --abbrev=0 2>/dev/null)
COMMIT=$(git rev-parse --short HEAD)

# Проверяем рабочее дерево
if git diff --quiet && git diff --cached --quiet; then
    DIRTY=""
else
    DIRTY="-dirty"
fi

if [ -n "$TAG" ]; then
    VERSION="${TAG}-${COMMIT}${DIRTY}"
else
    VERSION="${COMMIT}${DIRTY}"
fi

echo "${VERSION}"