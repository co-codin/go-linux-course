#!/bin/bash
set -e
cd /workspace/go-linux-course/deep/d01
./d01 server &
SPID=$!
sleep 0.4
echo
echo "=== attaching strace to pid $SPID while firing load ==="
echo
# capture strace for a short window overlapping the load
timeout 3 strace -f -e accept4,read,write,epoll_wait,epoll_ctl -p "$SPID" 2> /tmp/d01-strace.txt &
ST=$!
sleep 0.2
./d01 load
wait $ST 2>/dev/null || true
kill $SPID 2>/dev/null || true
wait $SPID 2>/dev/null || true
echo
echo "=== strace highlights (accept4 / epoll / read|write sample) ==="
grep -E 'accept4|epoll_ctl|epoll_wait|read\(|write\(' /tmp/d01-strace.txt | head -40
echo
echo "… total matching lines: $(grep -cE 'accept4|epoll_ctl|epoll_wait|read\(|write\(' /tmp/d01-strace.txt || echo 0)"
echo
echo "Takeaway: under load the server is mostly accept4 + epoll_wait; Ms aren't blocked per request."
sleep 10
