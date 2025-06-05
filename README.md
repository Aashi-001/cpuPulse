# cpuPulse

- this allows you to run any command and note its cpu and memory usage
- particularly useful for burst processes (when you have to note these values but the command runs so fast that you cannot note pid and check using top etc.)
- can handle keyboard interrupts as well.
- still in version 1
- next versions aim to log, plot cpu usages over time based on requirements.

## update 1
- version 2 built -> logs cpu and memory data of each sample in a csv file.
- plots cpu and memory usage vs samples and saves it as a png file.