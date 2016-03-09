Transport is required.

  $ [[ -n "$TRANSPORT" ]] || (echo 'Please provide an $TRANSPORT' && exit 1)

The trap ensures background processes are killed on exit.

  $ trap 'kill $(jobs -p)' EXIT

Test code:

  $ $TESTDIR/server/server &
  $ $TESTDIR/client/client -outbound=$TRANSPORT << INPUT
  > get foo
  > set foo bar
  > get foo
  > set baz qux
  > get baz
  > INPUT
  foo = 
  foo = bar
  baz = qux