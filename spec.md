# Tick Spec v0.0.1

> **Tick** - A dead-simple, line-oriented task syntax for personal use.

---

## Design Goals

1. **One line = one task.** No nesting, no indentation rules.
2. **Human first.** Readable at 2 AM without a cheat sheet.
3. **Parser second.** A 20-line script should be able to extract everything.
4. **Order-agnostic metadata.** Tokens before `::` can be shuffled without breaking the parse.
5. **Inline tagging.** `#tags` and `@mentions` can live naturally inside the title.

---

## File Structure

A `.tick` file is a plain text file encoded in UTF-8. It contains:

- **Comment lines** - start with `;`
- **Task lines** - one per task
- **Blank lines** - ignored by the parser

### Example File

```text
;tick:0.0.1

; quick brain dump - no project needed
- :: Buy oat milk
- :: Reply to landlord about lease renewal

; work stuff
/ =2h +groww #android :: Fix trace batching crash on cold start
- =1d p1 +groww #backend :: Publish error-code table for orders API ; review with @saheb first
- +groww #infra #blocked :: Spike new log drain for staging ; waiting on cost estimate
x =3h +groww #fno :: Hook P&L dashboard to ledger API ; pagination spec from @saheb

; side project
/ =30m +blog :: Draft post on Tick format
~ +blog :: Add comments system ; killed scope - not worth the spam fight
```

---

## Task Line Grammar

```
<status> [<priority>] [=<duration>] [+<project>] [<head-tags>...] :: <title> [; <note>]
```

### Status (required, position 0)

The very first character of every task line.

| Char | State        | Meaning                               |
|------|--------------|---------------------------------------|
| `-`  | **Todo**     | Not started yet                       |
| `/`  | **Active**   | In progress / currently working on it |
| `x`  | **Done**     | Completed                             |
| `~`   | **Dropped**  | Won't do; kept for the record         |

Four states, no cheat sheet. `~` is dropped only — if you need deferred, blocked, or on-hold, use a `#tag` in the header (it's just as grep-able and doesn't burn a status slot).

### Priority (optional, pre-`::` only)

A space-delimited token placed anywhere before `::`.

| Token | Level    |
|-------|----------|
| `p1`  | High     |
| `p2`  | Medium   |
| `p3`  | Low      |

If omitted, priority is considered unset (default). Lower number = higher urgency, same convention as P1/P2/P3 incident levels.

### Duration (optional, pre-`::` only)

A space-delimited token prefixed with `=`. Format: `=<number><unit>`.

| Example | Meaning    |
|---------|------------|
| `=15m`  | 15 minutes |
| `=2h`   | 2 hours    |
| `=1d`   | 1 day      |
| `=1w`   | 1 week     |

**Unit suffixes:**

| Suffix | Unit    |
|--------|---------|
| `m`    | minutes |
| `h`    | hours   |
| `d`    | days    |
| `w`    | weeks   |

No months - use `=4w` if you need that scale. This avoids the `m` vs `M` case-sensitivity trap.

Use duration for **estimates**, **logged time**, or **time-boxing** - the format does not distinguish. Interpretation is up to the user and their tooling.

### Project (optional, pre-`::` only)

Prefixed with `+`. At most one per task line.

```
+groww
+personal
+website
```

Tasks without a `+project` are valid - they're simply untagged. No need for a junk-drawer `+misc`.

### Tags (optional, anywhere)

Prefixed with `#`. Can appear:
- In the header (before `::`) for quick scanning
- Inline inside the `<title>` or `<note>` for natural language flow

```
#applens #bugs #backend
```

The parser extracts **all** `#words` from the entire line.

### Mentions (optional, anywhere)

Prefixed with `@`. Same rules as tags - header or inline.

```
@sahebg @devops @backend-team
```

### Delimiters

| Symbol | Role                                           |
|--------|------------------------------------------------|
| `::`   | **Hard separator.** Splits machine metadata (left) from human narrative (right). Required. |
| `;`    | **Note separator.** Everything after `;` in the right-hand side is treated as an inline comment / note. Optional. |

### Title & Note

- **Title**: The human-readable description of the task. Lives between `::` and `;` (or end-of-line if no `;`).
- **Note**: Optional context, links, or reminders. Lives after `;` to end-of-line.

---

## Links

Raw URLs are written as-is. The parser auto-detects `http://` and `https://` patterns inside the title or note.

```text
- =1h +docs :: Review deployment guide ; https://internal.dev/docs/deploy
```

> **Caveat:** URLs containing `#fragments` may be picked up as `#tags` by a naive parser. For personal use, either accept the false positive or avoid fragment URLs in task lines.

---

## Comments

Lines beginning with `;` are comments and are ignored by the parser.

```text
;tick:0.0.1
; 2026-05-15
; dumped from brain this morning
```

A `;` appearing after `::` does **not** start a comment - it starts the inline note.

---

## Parsing Rules

1. **Split on `::`** (first occurrence only). Left = metadata, Right = narrative.
2. **Split narrative on `;`** (first occurrence only). Left = title, Right = note.
3. **Tokenize metadata** by whitespace. Extract structural tokens by prefix/pattern (`+`, `p1`/`p2`/`p3`, `=`).
4. **Scan the entire line** with simple regex for `#tags` and `@mentions`.
5. **Scan the entire line** for `https?://` URLs.
6. **Ignore** blank lines and lines starting with `;`.

### Pseudocode

```
if line is empty or starts with ";": skip

status  = line[0]
rest    = line[1:]

meta, narrative = rest.split("::", 1)
title, note     = narrative.split(";", 1) if ";" in narrative else (narrative, "")

tokens = meta.split()
project  = first token starting with "+" if present, else null
priority = first token matching /^p[1-3]$/ if present, else null
duration = first token starting with "=" if present, else null

tags     = all "#word" matches in line
mentions = all "@word" matches in line
urls     = all "https?://..." matches in line
```

---

## Constraints & Edge Cases

| Situation | Rule |
|-----------|------|
| Missing `::` | Invalid task line. Parser may skip or treat entire line as title. |
| No `+project` | Valid. Task is simply untagged to any project. |
| Multiple `+project` tokens | Undefined. First one wins. Don't do it. |
| Multiple `=duration` tokens | Undefined. First one wins. Don't do it. |
| `p1` inside title (e.g., "update p1 docs") | Fine. Only tokens before `::` are scanned for priority. |
| `=` inside title (e.g., "API = broken") | Fine. Only tokens before `::` are scanned for duration. |
| `#` or `@` inside URLs | Naive parsers may over-capture. Acceptable for personal use. |
| Unicode in title / tags | Allowed. Tags and mentions match `\w+` (locale-dependent) or equivalent. |

---

## Version

`tick:0.0.1`

Bump patch for clerical tweaks; minor for additive changes (new optional token types); major for breaking grammar changes.

---

## Changelog

- **v0.0.1** — First numbered Tick release: flat one-line tasks, 4 statuses (`-`, `/`, `x`, `~`), `p1`/`p2`/`p3` priority, `=` duration (m/h/d/w), optional `+project`, `#tags` / `@mentions` anywhere, URLs, `tick:<version>` header comment.

---

*Keep it flat. Keep it fast.*
