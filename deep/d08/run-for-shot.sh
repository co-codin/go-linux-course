#!/bin/bash
set -e
cd /workspace/go-linux-course/deep/d08
./d08 &
PID=$!
sleep 0.3
echo
echo "=== strace -e futex sample while contended ==="
timeout 2 strace -f -e futex -p "$PID" 2>&1 | grep -E 'FUTEX_WAIT|FUTEX_WAKE|futex\(' | head -25 || true
wait $PID 2>/dev/null || true
echo
echo "(pprof was up on :18088 during the run)"
sleep 2
