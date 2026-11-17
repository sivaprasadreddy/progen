---
name: create-spec
description: Accept a feature description or a UseCase ID from docs/requirements.md, interview the human one question at a time to clarify intent and business rules, then write reviewable behavioral requirements and acceptance criteria before implementation.
disable-model-invocation: true
---

# Behavior Specification

## Purpose
Establish what the software must do before an agent chooses how to implement it.

## Starting input
- Accept either a feature description or a UseCase ID such as `UC-001` from `docs/requirements.md`.
- When an ID is supplied, read the matching summary row and detailed use case, if present. Match IDs case-insensitively and preserve the document's canonical ID, title, and dependency IDs in the specification. Use its description and acceptance criteria as the starting requirements for the interview; ask about remaining gaps rather than repeating documented answers.
- If the requirements file is unavailable or the ID cannot be uniquely resolved, ask the human to clarify the source or ID before proceeding. Resolve conflicts between the request and documented requirements during the interview.
- When only a feature description is supplied, read `docs/requirements.md` to check whether it matches an existing use case. Reuse a matching use case's ID; if the match is ambiguous, clarify it during the interview. For a new requirement, continue without requiring an existing UseCase ID and register it after confirmation as described below.

## Workflow
1. Read the available request, issue, specification, and relevant project context.
2. Separate explicit requirements from assumptions and inferred details.
3. Identify ambiguities, conflicting rules, missing failure behavior, boundary conditions, and relevant security or data-integrity concerns.
4. Begin an adaptive interview with the most important gap in understanding. Ask exactly one question per turn, addressing a single decision or topic, then wait for the human's answer. Keep follow-ups, confirmation, and file-location or overwrite questions in separate turns too; do not bundle subquestions into one sentence, list, or tool call.
5. Use each answer to update the understanding and choose the next question. Follow up on vague answers or contradictions before moving on. Explore relevant edge cases and negative scenarios as they emerge, clearly labeling suggestions that need a business decision. Offer brief examples or choices when useful, while keeping each turn focused on one question. Ask about implementation details only when they affect expected behavior.
6. Continue until the clarity criteria below are met, using prior answers rather than repeating questions. Keep track of confirmed details and unresolved points; defer drafting or saving the specification until the interview is complete. Then summarize the understood behavior, scope, and assumptions in the conversation and ask one question inviting confirmation or correction. Wait for the response; silence is not confirmation. If corrected, resolve the remaining gaps one question at a time and confirm the revised understanding before proceeding.
7. For a new requirement, assign the next unused sequential `UC-###` ID after the highest existing ID in `docs/requirements.md`. Add a row to its use-case summary and a detailed use-case entry with the confirmed title, description, dependencies, and acceptance criteria, following the document's format. Set the new requirement's status to `NOT_IMPLEMENTED`. Use the same ID, title, and dependencies in the specification. Preserve existing entries, IDs, and implementation statuses; creating a specification for an existing use case must not add a duplicate requirement.
8. Produce a concise behavior specification using the specification format below. Describe the goal, triggers, main success paths, validation and failure behavior, and relevant boundary and edge cases under `UseCase Details`. State invariants under `Business Rules` and acceptance criteria as observable outcomes. Record any explicitly deferred, nonessential decisions under the relevant section.
9. Save the specification as `{number}-{feature-title}.md` in the feature's agreed documentation location (default to `docs/specs/`). Use a short, filesystem-safe feature title (lowercase kebab-case); do not overwrite an existing specification without checking with the human. Present the specification and any added requirements entry for human review. Do not treat unresolved assumptions as approved requirements.

## Specification format
Use the following metadata keys and headings in this order. The metadata values shown are examples: use the feature's use-case ID, title, and dependency IDs, following existing project conventions. Use `[]` when there are no dependencies. Set the initial status to `NOT_IMPLEMENTED`.

```markdown
---
UseCase ID : UC-001
Title: Initial Project Setup
Depends On: []
Status: NOT_IMPLEMENTED
---

## Actors
## Preconditions
## UseCase Details
## Business Rules
## Acceptance Criteria
## In Scope
## Out of scope
```

Populate every section with the confirmed behavior; use `None` when a section does not apply.

## Clarity criteria
The interview is complete when the human has established:
- The actors, their goals, and the problem the feature addresses.
- The scope, exclusions, triggers, and preconditions.
- The main success paths and observable outcomes.
- The business rules, validation, failure behavior, and boundaries that materially affect correctness.
- Enough detail to state testable acceptance criteria without inventing requirements.

Use these as coverage criteria, not a fixed questionnaire. An answer may resolve several areas at once. Resolve material uncertainty before declaring clarity; record explicitly deferred, nonessential decisions as open questions. Once the criteria and confirmation are satisfied, write the specification instead of continuing the interview.

## Rules
- The human owns business intent and final decisions about expected behavior.
- Do not silently invent business rules or convert technical preferences into business requirements.
- Prefer observable behavior over prescribing classes, methods, frameworks, or internal design.
- Reuse existing project conventions and documented decisions; call out conflicts rather than silently overriding them.
- Keep requirements testable, specific, and concise.
- Do not begin feature implementation until the behavior and material ambiguities are approved.
