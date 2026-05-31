---
mode: agent
description: Implement a feature end to end in this repo
---

Implement the requested feature in this repository with minimal, targeted changes.

Inputs:
- Feature request: ${input:feature_request}
- Files/modules likely involved (optional): ${input:candidate_files}
- Constraints (optional): ${input:constraints}

Process:
1. Inspect the relevant code paths and summarize impact.
2. Implement the smallest safe change that satisfies the request.
3. Add or update tests for behavior changes.
4. Run `make test` and report outcomes.
5. Summarize changed files and any follow-up recommendations.

Output:
- Concise implementation summary
- Test results
- Risks or assumptions
