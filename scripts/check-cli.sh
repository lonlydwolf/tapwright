#!/usr/bin/env bash

set -euo pipefail

tmp_dir="$(mktemp -d)" || exit 1
trap 'rm -rf -- "$tmp_dir"' EXIT

test_number=0
total_tests=9
export TAPWRIGHT_TAPS_ROOT="$tmp_dir/root"
mkdir -p "$TAPWRIGHT_TAPS_ROOT"

if ! go build -o "$tmp_dir/tapwright" .; then
  printf 'FAIL: build failed.\n' >&2
  exit 1
fi

program_status=0
"$tmp_dir/tapwright" >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright\`"

printf 'Usage: tapwright <command>\n' >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stdout" ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stderr"; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

program_status=0
"$tmp_dir/tapwright" inspect >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright inspect\`"

printf 'Usage: tapwright inspect <argument>\n' >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stdout" ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stderr"; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

program_status=0
"$tmp_dir/tapwright" inspect arg1 arg2 >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright inspect arg1 arg2\`"

printf 'Usage: tapwright inspect <argument>\n' >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stdout" ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stderr"; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

program_status=0
"$tmp_dir/tapwright" inspect invalid_arg >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright inspect invalid_arg\`"

printf 'Error: invalid scheme "invalid_arg": must be https://\n' >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stdout" ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stderr"; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

program_status=0
"$tmp_dir/tapwright" unknown >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright unknown\`"

printf 'Unknown command: unknown\n' >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stdout" ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stderr"; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

program_status=0
"$tmp_dir/tapwright" inspect https://github.com/lonlydwolf/tapwright >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright inspect https://github.com/lonlydwolf/tapwright\`"

printf '{"owner":"lonlydwolf","repository":"tapwright","canonical_url":"https://github.com/lonlydwolf/tapwright"}\n' >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 0 ]; then
  printf 'FAIL: expected exit status 0; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stderr" ]; then
  printf 'FAIL: expected empty stderr.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stdout"; then
  printf 'FAIL: stdout not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

program_status=0
"$tmp_dir/tapwright" create alice tap >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright create alice tap\`"

printf '%s/alice/homebrew-tap\n' "$TAPWRIGHT_TAPS_ROOT" >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 0 ]; then
  printf 'FAIL: expected exit status 0; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stderr" ]; then
  printf 'FAIL: expected empty stderr.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stdout"; then
  printf 'FAIL: stdout not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

program_status=0
"$tmp_dir/tapwright" create alice homebrew-tap >"$tmp_dir/tapwright.stdout" 2>"$tmp_dir/tapwright.stderr" || program_status=$?
test_number=$((test_number + 1))

echo "Testing: \`tapwright create alice homebrew-tap\`"

printf "Error: invalid tap name \"homebrew-tap\": tap name must not start with \`homebrew-\`\n" >"$tmp_dir/tapwright.expected"

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stdout" ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s "$tmp_dir/tapwright.expected" "$tmp_dir/tapwright.stderr"; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

missing_root="$tmp_dir/nonexistent"

program_status=0
TAPWRIGHT_TAPS_ROOT="$missing_root" \
  "$tmp_dir/tapwright" create alice tap \
  >"$tmp_dir/tapwright.stdout" \
  2>"$tmp_dir/tapwright.stderr" || program_status=$?

test_number=$((test_number + 1))

echo "Testing: \`tapwright create alice tap\` with nonexistent taps root"

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s "$tmp_dir/tapwright.stdout" ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! grep -Fq 'Error: taps root unavailable' "$tmp_dir/tapwright.stderr"; then
  printf 'FAIL: stderr not matching expected error.\n' >&2
  exit 1
fi

if [ -e "$missing_root" ] || [ -L "$missing_root" ]; then
  printf 'FAIL: expected taps root to remain nonexistent.\n' >&2
  exit 1
fi

echo "PASS $test_number/$total_tests"

exit 0
