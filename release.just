# list release tags
tags:
    git tag -l --sort=-creatordate

# extract release notes for a tag from CHANGELOG.md
notes tag:
    #!/bin/bash
    set -euo pipefail
    VERSION="{{ tag }}"
    VERSION="${VERSION#v}"
    awk -v version="$VERSION" '
        $0 ~ "^## \\[" version "\\]" { found=1; next }
        /^## / { found=0 }
        found
    ' CHANGELOG.md | sed -e '/./,$!d' -e :a -e '/^\s*$/{ $d; N; ba; }' | cat -s

# push commits and a specific tag
push tag:
    git push
    git push origin "{{ tag }}"
