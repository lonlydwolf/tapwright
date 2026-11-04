#!/bash/sh

if ! go build -o /tmp/tapwright .; then
  printf 'FAIL: build failed.\n' >&2
  exit 1
fi

/tmp/tapwright >/tmp/tapwright.stdout 2>/tmp/tapwright.stderr
program_status=$?

echo "Testing: \`tapwright\`"

printf 'Usage: tapwright <command>\n' >/tmp/tapwright.expected

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s /tmp/tapwright.stdout ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s /tmp/tapwright.expected /tmp/tapwright.stderr; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS 1/5"

/tmp/tapwright inspect >/tmp/tapwright.stdout 2>/tmp/tapwright.stderr
program_status=$?

echo "Testing: \`tapwright inspect\`"

printf 'Usage: tapwright inspect <argument>\n' >/tmp/tapwright.expected

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s /tmp/tapwright.stdout ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s /tmp/tapwright.expected /tmp/tapwright.stderr; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS 2/5"

/tmp/tapwright inspect arg1 arg2 >/tmp/tapwright.stdout 2>/tmp/tapwright.stderr
program_status=$?

echo "Testing: \`tapwright inspect arg1 arg2\`"

printf 'Usage: tapwright inspect <argument>\n' >/tmp/tapwright.expected

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s /tmp/tapwright.stdout ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s /tmp/tapwright.expected /tmp/tapwright.stderr; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS 3/5"

/tmp/tapwright inspect invalid_arg >/tmp/tapwright.stdout 2>/tmp/tapwright.stderr
program_status=$?

echo "Testing: \`tapwright inspect invalid_arg\`"

printf 'Error: invalid scheme "": must be https\n' >/tmp/tapwright.expected

if [ "$program_status" -ne 1 ]; then
  printf 'FAIL: expected exit status 1; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s /tmp/tapwright.stdout ]; then
  printf 'FAIL: expected empty stdout.\n' >&2
  exit 1
fi

if ! cmp -s /tmp/tapwright.expected /tmp/tapwright.stderr; then
  printf 'FAIL: stderr not matching expected output \n' >&2
  exit 1
fi

echo "PASS 4/5"

/tmp/tapwright inspect https://github.com/lonlydwolf/tapwright >/tmp/tapwright.stdout 2>/tmp/tapwright.stderr
program_status=$?

echo "Testing: \`tapwright inspect https://github.com/lonlydwolf/tapwright\`"

printf '{"owner":"lonlydwolf","repository":"tapwright","canonical_url":"https://github.com/lonlydwolf/tapwright"}\n' >/tmp/tapwright.expected

if [ "$program_status" -ne 0 ]; then
  printf 'FAIL: expected exit status 0; got %s\n' "$program_status" >&2
  exit 1
fi

if [ -s /tmp/tapwright.stderr ]; then
  printf 'FAIL: expected empty stderr.\n' >&2
  exit 1
fi

if ! cmp -s /tmp/tapwright.expected /tmp/tapwright.stdout; then
  printf 'FAIL: stdout not matching expected output \n' >&2
  exit 1
fi

echo "PASS 5/5"
exit 0
