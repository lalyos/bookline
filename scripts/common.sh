VERSION=$(git describe --abbrev=0 --tags)
GIT_REV=$(git rev-parse --short HEAD)$([[ -z $(git status --porcelain) ]] || echo "-dirty")
