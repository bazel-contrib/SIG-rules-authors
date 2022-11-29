#!/usr/bin/env bash

set -o errexit -o pipefail -o nounset

SCRIPT_DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
DATA="$SCRIPT_DIR"/rulesets.json

if ! command -v jq &> /dev/null
then
    echo >&2 "Please install jq to use this utility"
    exit
fi

# Prompt for info we commonly need for each ruleset
read -r -p 'GitHub Repository (http://github.com/???): ' ghrepo

tmp=$(mktemp)
jq ".rulesets += [{ghrepo: \"$ghrepo\"}]" < "$DATA" > "$tmp" && mv "$tmp" "$DATA"
