declare-option str kakkaktempprefixorsomething /tmp/kakkakkak

declare-option -hidden str kakkakchatoutfifo
declare-option -hidden str kakkakchatinfifo

declare-option str kakkaklockfile /tmp/kakkakkak.lock
declare-option str kakkakprogram kakkak
declare-option bool kakkakstarted false

define-command opengpt -hidden -override %{
  try %{
    buffer chatkakkak
  } catch %{
    eval -try-client tools %{
      edit -fifo %opt{kakkakchatoutfifo} chatkakkak
    }
  }
}

define-command gpt -override -params 0.. %{
  evaluate-commands %sh{
    if [ "$kak_opt_kakkakstarted" = false ]; then
      echo start-kakkak
    fi
  }

  nop %sh{
    if [ $(($(printf %s "${kak_selection}" | wc -m))) -gt 1 ]; then
      echo "$@ " "$(printf '%s' "${kak_selection}" | tr '\n' ' ')" > $kak_opt_kakkakchatinfifo
    else
      echo "$@" > $kak_opt_kakkakchatinfifo
    fi
  }

  opengpt
}

define-command kakkakreifywith -override -params 1..2 %{
  evaluate-commands %sh{
    printf %s "define-command start-kakkak -override -params 0 %{
      eval %sh{
        infifo=\$(mktemp -u \"\${kak_opt_kakkaktempprefixorsomething}XXXXXXXX\")
        mkfifo \$infifo
        echo \"set-option global kakkakchatinfifo \$infifo\"
        outfifo=\$(mktemp -u \"\${kak_opt_kakkaktempprefixorsomething}XXXXXXXX\")
        mkfifo \$outfifo
        echo \"set-option global kakkakchatoutfifo \$outfifo\"
        (eval GEMINI_API_KEY=$1 GEMINI_DEBUG=$2 \$kak_opt_kakkakprogram \$infifo \$outfifo 2>&1 & ) > /dev/null 2>&1 < /dev/null
        echo \"set-option global kakkakstarted true\"
      }
    }"
  }
}
