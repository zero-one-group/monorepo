#!/usr/bin/env bash
# Tiered coverage gate over a single Go cover profile.
#
#   coverage-gate.sh build/coverage.out
#
# Tiers (env, statement coverage %, all default to 0 = report only):
#   COVERAGE_MIN                 everything in the profile
#   COVERAGE_MIN_MODULES         <module>/modules/...            (business code)
#   COVERAGE_MIN_CRITICAL        packages listed in COVERAGE_CRITICAL_PACKAGES
#   COVERAGE_CRITICAL_PACKAGES   comma-separated path prefixes relative to the module root,
#                                e.g. "modules/order,modules/payment,internal/webhook"
#
# A tier with no statements in the profile is reported as "n/a" and passes, so the
# critical tier activates automatically once the first listed package has code.
#
# Ratchet rule: set each floor to (measured - 5) and only ever raise it. Record the
# history next to the env values in moon.yml so the next person sees the trend.
set -euo pipefail

profile="${1:?usage: coverage-gate.sh <coverprofile>}"
[ -s "$profile" ] || { echo "coverage gate: profile $profile is missing or empty (did the tests run?)"; exit 1; }
min_all="${COVERAGE_MIN:-0}"
min_modules="${COVERAGE_MIN_MODULES:-0}"
min_critical="${COVERAGE_MIN_CRITICAL:-0}"

# Go module path, so the patterns work whatever the generated project is called.
module="$(cd "$(dirname "$profile")/.." 2>/dev/null && go list -m 2>/dev/null || true)"
[ -n "$module" ] || module="$(go list -m 2>/dev/null || true)"
[ -n "$module" ] || { echo "coverage gate: cannot determine Go module (run from the module root)"; exit 1; }

critical_pattern=""
if [ -n "${COVERAGE_CRITICAL_PACKAGES:-}" ]; then
  IFS=',' read -ra pkgs <<<"$COVERAGE_CRITICAL_PACKAGES"
  for p in "${pkgs[@]}"; do
    p="${p#"${p%%[![:space:]]*}"}"; p="${p%"${p##*[![:space:]]}"}"; p="${p%/}"
    [ -n "$p" ] || continue
    critical_pattern="${critical_pattern:+$critical_pattern|}^$module/$p/"
  done
fi

# Statement-weighted coverage for files matching a regex (profile lines: file:a.b,c.d numstmts count).
tier_pct() {
  awk -v pat="$1" '
    NR > 1 && $1 ~ pat { split($0, f, " "); total += f[2]; if (f[3] + 0 > 0) covered += f[2] }
    END { if (total == 0) print "n/a"; else printf "%.1f", covered * 100 / total }
  ' "$profile"
}

fail=0
check() {
  local name="$1" pct="$2" min="$3"
  if [ "$pct" = "n/a" ]; then
    printf '  %-10s %6s   (no statements yet, minimum %s%%)\n' "$name" "$pct" "$min"
    return
  fi
  if awk -v p="$pct" -v m="$min" 'BEGIN { exit !(p + 0 < m + 0) }'; then
    printf '  %-10s %5s%%   BELOW minimum %s%%\n' "$name" "$pct" "$min"
    fail=1
  else
    printf '  %-10s %5s%%   (minimum %s%%)\n' "$name" "$pct" "$min"
  fi
}

echo "coverage gate ($profile, module $module)"
check overall  "$(tier_pct '.')"                          "$min_all"
check modules  "$(tier_pct "^$module/modules/")"          "$min_modules"
if [ -n "$critical_pattern" ]; then
  check critical "$(tier_pct "$critical_pattern")"        "$min_critical"
else
  printf '  %-10s %6s   (set COVERAGE_CRITICAL_PACKAGES to enable)\n' critical "-"
fi

[ "$fail" -eq 0 ] || { echo "coverage below minimum"; exit 1; }
