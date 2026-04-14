<div align="center">
  <img src="./static/whitebox.jpg" alt="Whitebox" />
</div>

# Whitebox

A personal AI assistant where context is explicitly controlled.

Whitebox is designed for cases where reproducibility, cost control, and understanding **what the model actually sees** matter.

---

## 1. What is Whitebox

Whitebox is a minimal runtime for working with LLMs.

It exists to remove hidden logic from agent systems:

- context is defined by files, not “magic”
- every step and cost can be observed (e.g. via Langfuse)

Core idea: **context is a product interface, not an internal detail of the model**.

---

## 2. Problem

Most AI systems behave like black boxes:

- context grows without limits
- noise leaks into the prompt
- it’s unclear why the model answered the way it did
- token usage becomes unpredictable

Result: more context ≠ better answers.

---

## 3. Core Idea

Whitebox is built around three principles:

- **Context as a single source of truth**: the model sees only what is inside `~/.whitebox/context`
- **No hidden magic**: no implicit system layers, hidden injections, or automatic decisions
- **System controls context, not LLM**: focus is managed by the system and file structure

---

## 4. Architecture

### Engine

Execution loop:

- builds prompt from Context
- calls the LLM
- checks for tool calls
- executes tools if needed and continues the loop
- returns final answer

### Context

Context is assembled from directories:

1. `minds/`
2. `memories/`
3. `skills/`
4. `tools/`
5. `sessions/` (chat history)

Order is fixed and deterministic.

### Tools

Tool calls are plain JSON returned by the model.

Example:

```json
{
  "tool": "read_file",
  "arguments": {
    "path": "notes/todo.md"
  }
}
````

### State

Execution state is simple:

* current input
* model output
* number of steps in the call chain

No hidden background agents.

### TUI

CLI interface built with Bubble Tea:

* user input
* live events (`debug`, `tool_call`, `final`)
* session history persistence

---

## 5. How it works

Execution flow:

1. User input
2. Engine loop
3. LLM call
4. If tool is needed → execute tool
5. Tool result goes back into the loop
6. Final answer is generated
7. Answer is saved in session history

---

## 6. Context System

This is the core layer of Whitebox.

After first run, the structure is created:

```text
~/.whitebox/
  context/
    minds/
    memories/
    skills/
    tools/
    sessions/
  workspace/
```

Key points:

* everything in `context/*` is included in the prompt
* anything outside does not exist for the model
* context can be explicitly cleaned, compressed, and rebuilt

Practical effect:

* less noise
* lower token cost
* more stable outputs

---

## 7. Features

* Transparent file-based context
* Explicit execution loop
* Tool calling via JSON
* Observability via events and debug mode
* Minimal architecture without hidden layers

---

## 8. Getting Started

### 1) Install

```bash
git clone https://github.com/Alice088/whitebox.git
cd whitebox
go mod download
```

### 2) Configure

```bash
cp .env.example .env
```

Set at least:

* `LLM_API_KEY`
* `SESSION_MAX_MESSAGES`
* `CALL_CHAIN_MAX`

### 3) Run

Local provider:

```bash
go run . --model <your-model> --provider local
```

API provider (e.g. DeepSeek):

```bash
go run . --model <your-model> --provider api --provider_name deepseek
```

### 4) First request

Type a message in the TUI and press Enter.

---

## 9. Examples

### Basic conversation

```text
You: Create a migration plan for a Go service.
Assistant: ...
```

### Reading a file via tool

Model may return:

```json
{
  "tool": "read_file",
  "arguments": { "path": "docs/spec.md" }
}
```

Engine reads the file from `~/.whitebox/workspace/docs/spec.md`, returns the result to the model, and requests a final answer.

### Working with context

1. Add `minds/product.md`
2. Add `skills/code-review.md`
3. Restart session

These files become part of the system prompt.

---

## 10. Philosophy

Whitebox follows simple rules:

* **Minimalism > complexity**
* **Control > automation**
* **Explicitness > magic**

This is not a “smart autonomous agent”.

It is an engineering tool where you control context, cost, and behavior.

---

## Project status

Whitebox is under active development.
If you care about reproducibility and observability in LLM systems — this is a solid foundation.
