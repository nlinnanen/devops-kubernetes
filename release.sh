#!/usr/bin/env bash
set -euo pipefail

arg="${1:-}"
if [[ -z "$arg" ]]; then
  echo "usage: $0 <release-id>"
  exit 1
fi

# Adjust these if you want different formatting
tag="${arg}"
title="Exercise ${tag}"

readme="Readme.md"
if [[ ! -f "$readme" ]]; then
  echo "ERROR: ${readme} not found"
  exit 1
fi

url="https://github.com/nlinnanen/devops-kubernetes/tree/${arg}/todo_app"

printf "\n- [${arg}](${url})" >> "$readme"

git add "$readme"
git commit -m "chore(release): ${tag}"
git push

# Create GitHub release (requires gh)
# --notes can be swapped for --generate-notes if you prefer
gh release create "$tag" --title "$title"

copy_to_clipboard() {
  if command -v pbcopy >/dev/null 2>&1; then
    printf "%s" "$1" | pbcopy
  elif command -v wl-copy >/dev/null 2>&1; then
    printf "%s" "$1" | wl-copy
  elif command -v xclip >/dev/null 2>&1; then
    printf "%s" "$1" | xclip -selection clipboard
  elif command -v xsel >/dev/null 2>&1; then
    printf "%s" "$1" | xsel --clipboard --input
  else
    return 1
  fi
}

if copy_to_clipboard "$url"; then
  echo "release: $url"
  echo "copied to clipboard"
else
  echo "release: $url"
  echo "clipboard tool not found (pbcopy/wl-copy/xclip/xsel)."
fi
