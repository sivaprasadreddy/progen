---
UseCase ID : DSK-11
Title: Preview the Generated Project on the Review Step
Depends On: [DSK-01, DSK-03, DSK-04]
Status: COMPLETED
---

## Actors

- **Progen Desktop user**: someone on the wizard's Review step who wants to see exactly which files will be generated, and what they contain, before the project is saved as a ZIP.

## Preconditions

- The user has reached the Review step of the wizard (DSK-01).
- Progen Desktop can write to a temporary directory and to `~/Downloads`.

## UseCase Details

### Goal

Let the user check the generated project's file tree and file contents inside Progen Desktop, and then save that same project, without unzipping it first.

### Trigger

The user clicks **Preview** in the Review step's action bar.

### Main success path: preview

1. The user clicks **Preview** on the Review step.
2. The project fields (name, version, Group ID, Artifact ID, package) are checked the same way Generate checks them. If they're invalid, the field errors show and the preview doesn't start.
3. While the preview is prepared, the Review step stays on screen. The Preview button reads **"Preparing preview…"** and every action on the step is disabled. Preparation can't be cancelled.
4. Preparation runs the full generation, including the code formatter, so the files are identical to the ones a ZIP would contain. It doesn't block the UI.
5. When preparation finishes, a full-window **preview page** opens with:
   - a **file tree** of every generated file and folder, hidden ones (`.gitignore`, `.github/`, `.mvn/`, …) included. The project's root and its top-level entries are shown, and subfolders start collapsed.
   - a **file viewer** showing the selected file's contents read-only. The build file (`pom.xml` for Maven, `build.gradle` for Gradle) is selected when the page opens.
6. The user expands and collapses folders and selects files to view them.

### File viewer

- Text files show with line numbers. The user can scroll, select and copy the text but can't edit it.
- These file types get syntax highlighting: Java, XML, Gradle (Groovy DSL), YAML, `.properties`, SQL, Markdown, JSON, TOML, shell scripts, HTML (including Thymeleaf templates), JTE, TypeScript, JavaScript and CSS.
- Any other text file shows as plain text, still with line numbers.
- Selecting a binary file (such as `maven-wrapper.jar`, `gradle-wrapper.jar` or an image) shows *"Binary file, preview not available"* and the file's size instead of its contents.

### Leaving the preview

- **Back to Review** returns to the Review step with every selection and field unchanged.
- **Generate** saves the previewed files as `~/Downloads/<app-name>.zip`, adding a numeric suffix if that file exists (DSK-04), without generating the project again. It then shows the Complete step.
- The preview page has no Cancel button. To cancel, the user goes back to Review and uses Cancel there (DSK-05).

### Reusing a prepared preview

- A prepared preview stays valid while the selections in every step, and the name, version, Group ID, Artifact ID and package, are unchanged.
- While it is valid:
  - clicking **Preview** again opens the preview page right away, without preparing it again;
  - clicking **Generate** on the Review step saves the prepared files as the ZIP, without generating again.
- Changing any of those values discards the preview. This applies whether the change is a field edited on Review or a choice changed in an earlier step. The next Preview or Generate generates the project again.
- Cancel (clearing the selections) and "Generate another" on the Complete step also discard the preview.

### Failure behavior

- **Preparing the preview fails** (for example, the formatter fails or the disk is full): the user stays on the Review step, the error shows in the Review step's existing error box, and nothing from the failed attempt is kept. The user can change things and click Preview again.
- **Saving the ZIP from the preview page fails**: the error shows on the preview page and the user stays there. The prepared preview is kept.

### Edge cases

- A project with no frontend, or with any of the frontends (Thymeleaf, JTE, Angular, React), previews the same way. The tree shows whatever that configuration generates.
- If no ZIP exists yet, generating from the preview names it `<app-name>.zip`. If `<app-name>.zip` already exists in `~/Downloads`, it is named `<app-name>-1.zip`, `<app-name>-2.zip` and so on.

## Business Rules

1. The preview shows exactly the files, with exactly the contents, that a ZIP generated from the same choices would contain. The code formatter's output is included.
2. A ZIP saved while a valid preview exists contains the previewed files. The project isn't generated again.
3. A preview is valid only for the choices and project fields it was prepared from. Any change to them discards it.
4. Preview has the same validation as Generate. A preview is never prepared from invalid project fields.
5. File contents in the preview are read-only. Nothing the user does in the preview changes the generated project.
6. While a preview is being prepared, no other wizard action can be used.
7. The prepared files are deleted when the preview is discarded and when the app quits. They are never left in the user's Downloads directory.
8. A failed preparation leaves no prepared preview behind.

## Acceptance Criteria

1. The Review step shows a **Preview** button in its action bar.
2. Clicking Preview with an invalid project field shows that field's error and doesn't open the preview page.
3. While the preview is being prepared, the Preview button reads "Preparing preview…" and the Review step's other actions are disabled.
4. When preparation succeeds, the preview page opens. Its file tree lists every file in the generated project, hidden files included. Only the root's top-level entries are shown, and subfolders start collapsed.
5. When the preview page opens, `pom.xml` is selected for a Maven project and `build.gradle` for a Gradle project, and its contents are shown.
6. Selecting a text file shows its contents with line numbers. The contents can be selected and copied but not edited.
7. Java, XML, Gradle, YAML, `.properties`, SQL, Markdown, JSON, TOML, shell, HTML, JTE, TypeScript, JavaScript and CSS files are syntax highlighted. Other text files show as plain text with line numbers.
8. Selecting a binary file shows "Binary file, preview not available" and the file's size.
9. Every file's contents in the preview are byte-for-byte identical to the same file in the ZIP saved from that preview.
10. **Back to Review** shows the Review step with every selection and field unchanged.
11. **Generate** on the preview page saves `~/Downloads/<app-name>.zip`, or a suffixed name if that file exists, without running generation again, and then shows the Complete step.
12. With nothing changed since the preview was prepared, going back to Review and clicking Preview opens the preview page without preparing it again. Clicking Generate on Review saves the prepared files without generating again.
13. After any selection or project field changes, the next Preview prepares the project again and shows the new files.
14. If preparation fails, the Review step stays on screen, its error box shows the error, and the preview page doesn't open.
15. If saving the ZIP from the preview page fails, the error shows on the preview page and the user stays on it.
16. The prepared files are removed after Cancel, after "Generate another", after the preview is discarded because something changed, and when the app quits.
17. Generating from the Review step without previewing behaves exactly as it does today (DSK-04).

## In Scope

- A Preview button on the Review step of Progen Desktop.
- Preparing the full generated project, formatter included, as a preview.
- A full-window preview page with a file tree, a read-only viewer with line numbers and syntax highlighting, Back to Review, and Generate.
- Reusing the prepared preview for Preview and Generate while nothing has changed, and discarding it when something changes.
- Showing binary files in the tree without their contents.
- Handling preparation and save failures as described.
- Cleaning up the prepared files.

## Out of scope

- Editing generated files in the preview.
- Searching file contents or filtering the file tree.
- Cancelling a preview while it is being prepared.
- A preview in the `progen` CLI.
- Changing the ZIP's name, location or contents.
- Previewing from steps other than Review, or comparing two previews.
