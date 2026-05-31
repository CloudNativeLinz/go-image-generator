---
mode: agent
description: Debug a failing test, command, or runtime issue
---

Debug and fix a failure in this repository.

Inputs:
- Symptom or error: ${input:error}
- Repro command (optional): ${input:repro_command}
- Suspected file (optional): ${input:suspected_file}

Process:
1. Reproduce the issue (or reason from available logs/output).
2. Identify the root cause and explain it briefly.
3. Apply a minimal fix.
4. Add or adjust tests when relevant.
5. Run focused checks (and broader checks if needed).

Output:
- Root cause
- Fix summary
- Verification steps and results
- Remaining risks (if any)
