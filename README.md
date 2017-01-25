# tcat

tcat is cammand to output text with time.

# Usage

If you have the command to check high cpu usage and you want to record result of the command like "ps auxf" with time.

```
% check_if_cpu_high && ps auxf | tcat
2017-01-22 19:38:49: USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
2017-01-22 19:38:49: root         2  0.0  0.0      0     0 ?        S    Jan09   0:00 [kthreadd]
2017-01-22 19:38:49: root         3  0.0  0.0      0     0 ?        S    Jan09   1:49  \_ [ksoftirqd/0]
2017-01-22 19:38:49: root         5  0.0  0.0      0     0 ?        S<   Jan09   0:00  \_ [kworker/0:0H]

% check_if_cpu_high && ps auxf | tcat -f "%Y-%m-%d"
2017-01-22: USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
2017-01-22: root         2  0.0  0.0      0     0 ?        S    Jan09   0:00 [kthreadd]
2017-01-22: root         3  0.0  0.0      0     0 ?        S    Jan09   1:49  \_ [ksoftirqd/0]
2017-01-22: root         5  0.0  0.0      0     0 ?        S<   Jan09   0:00  \_ [kworker/0:0H]
```

You can use golang time expression(like "Mon Jan 2 15:04:05 MST 2006") or the follwoing charactor string for time expression.

* %m month
* %d day
* %H hour
* %M minute
* %S second
* %Y year
* %T %H:%M:%S
* %W name of day
* %h,%b,%B name of month

# Compatibilty of GNU cat

short options are supported.

*  -A    equivalent to -vET
*  -E    display $ at end of each line
*  -T    display TAB characters as ^I
*  -e    equivalent to -vE
*  -t    equivalent to -vT
*  -v    use ^ and M- notation, except for LFD and TAB
*  -n    number all output lines

# Author

Atsushi Kato (ktat)

# License

MIT: https://ktat.mit-license.org/2017
