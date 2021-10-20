#!/usr/bin/env bash
# Uses GitHub API to gather updates on the rulesets
# - show how active they are, recent releases, etc
# - TODO: what other data can we gather...

set -o errexit -o pipefail -o nounset

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
GH_TOKEN=${GH_TOKEN:-error you must set a github token}

function debug() {
    # echo $@
    true
}

function ghapi() {
    curl 2>/dev/null \
        -u "$USER:$GH_TOKEN" \
        -H "Accept: application/vnd.github.v3+json" \
        "https://api.github.com/$1"
}

# Make the function visible to subshells
export -f ghapi

node "$SCRIPT_DIR/rulesets.js" | \
  jq -r '.rulesets[]|[.ghrepo] | @tsv' |
  tr '\t' ',' |
  while IFS=$',' read -r ghrepo; do
    out="$SCRIPT_DIR/$ghrepo"
    mkdir -p "$out"
    ghapi "repos/$ghrepo" > "$out"/repo.json
    ghapi "repos/$ghrepo/stats/participation" | jq .all > "$out"/participation.json
    ghapi "repos/$ghrepo/community/profile" > "$out"/community_profile.json
    # rules_jvm_external has some tags starting with "2018"
    # so we filter those out here by looking only for tags with a dot
    latest_tag=$(ghapi repos/"$ghrepo"/tags |\
        jq -r '[.[] | select(.name | contains("."))] | .[0].name')
    
    ghapi "repos/$ghrepo/commits/$latest_tag" | jq -r .commit > "$out"/latest_tag.json
done

date > "$SCRIPT_DIR/last_updated.txt"
