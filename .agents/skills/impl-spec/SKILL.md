---
name: impl-spec
description: Implement one documented use case identified by a UC-### ID, including its tests and requirements status transitions. Use when the user asks to implement a specific use case from requirements.md; do not use to draft requirements or implement an unspecified collection of features.
disable-model-invocation: true
---

# Implement Use Case

Implement exactly one use case supplied by the user, using its description and acceptance criteria as the scope contract.

## Input and source of truth

- Require one use-case ID in `UC-###` form. Matching may be case-insensitive, but preserve the canonical ID from the document.
- Honor a user-supplied requirements path; otherwise use `docs/requirements.md`.
- Read the use case's summary row and complete detailed section. Also read directly referenced dependencies and relevant product or architecture documentation.
- If the ID does not exist, more than one ID was supplied, or the detailed section and summary disagree, stop and ask the user to resolve the problem. Do not edit files.
- Treat the documented description and acceptance criteria as the implementation contract. Do not silently expand the feature.

## Resolve questions before implementation

Inspect the relevant code and tests to determine the smallest complete implementation. Before changing any file, identify missing decisions, contradictions, technically significant ambiguity, unmet use-case dependencies, or acceptance criteria that cannot be implemented as written.

If any such question exists, ask the user concise, decision-oriented questions and wait for answers. Keep the use case's existing status unchanged while waiting. Do not choose product behavior on the user's behalf.

If the use case is already `COMPLETED`, report that fact and stop unless the user explicitly asks to reimplement or repair it. If a required dependency is not `COMPLETED`, ask whether to proceed, implement the dependency separately, or stop; never expand the current use case implicitly.

## Begin the implementation

Once the scope is unambiguous and prerequisites are resolved, make the status transition the first repository change:

1. Change the selected use case's status to `IN_PROGRESS` in both the summary table and detailed section.
2. Confirm the two entries agree before editing implementation files.
3. If the status is already `IN_PROGRESS`, keep it and continue after confirming the user wants the existing work resumed when that intent is not clear from the request.

Do not change the status of any other use case.

## Implement and verify

- Follow repository instructions and existing architecture, naming, and testing conventions.
- Preserve unrelated user changes in the worktree.
- Implement every acceptance criterion, including stated authorization, validation, failure, privacy, and state-transition behavior.
- Add or update focused automated tests that demonstrate the acceptance criteria. Prefer observable behavior over tests coupled to implementation details.
- Run the focused tests first, then the repository's broader relevant verification when practical.
- Review the final diff for scope, accidental changes, and consistency with the use case.

If implementation reveals a new product ambiguity, pause and ask the user. Leave the status as `IN_PROGRESS`. Likewise, if tests fail or required work remains, do not mark the use case complete; report the remaining work and verification evidence.

## Complete the use case

Only after all acceptance criteria are implemented and relevant verification passes:

1. Change the status from `IN_PROGRESS` to `COMPLETED` in both the summary table and detailed section.
2. Confirm both status entries match and no other use-case status changed.
3. Report the implemented behavior, tests run and their results, and the final status.

`COMPLETED` means the documented use case is fully implemented and verified, not merely that code was written.
