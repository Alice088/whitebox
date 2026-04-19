<div align="center">
  <img src="./static/whitebox.png" alt="Whitebox" />
</div>

---

A minimal FSM-based core for building fast, controlled AI agents.

Whitebox is not a ready-made assistant.  
It is a small runtime designed to build simple, specialized agents with full control.

---

## 1. What is Whitebox

Whitebox is a lightweight core built around an FSM orchestrator.

It is designed to:

- execute simple agent loops
- keep behavior predictable
- avoid hidden complexity

You use it as a base to build your own agents.

---

## 2. Problem

Most agent systems fail for one reason: uncontrolled context.

- context grows without limits  
- irrelevant data leaks into prompts  
- hidden system injections affect behavior  
- token usage becomes unpredictable  

Result:

- higher API cost  
- unstable outputs  
- hard-to-debug behavior  

---

## 3. Core Idea

Whitebox fixes this by making context explicit.

- context is loaded only from files  
- nothing is injected implicitly  
- no hidden system layers  

The only fixed injection is the output protocol.

Everything else is under your control.

---

## 4. Context System

Context is simple and transparent.

```text
~/.whitebox-name/
- context/
    - minds/
    - memories/
    - skills/
    - tools/
    - sessions/
- workspace/
````

How it works:

* files from `context/` are read and added to the prompt
* nothing outside this directory exists for the model
* no automatic memory, no hidden prompts

This gives:

* predictable inputs
* stable outputs
* controlled token usage

You decide exactly what the model sees.

---

## 5. Output Protocol

The only built-in rule is a strict output format.

The model must return JSON:

```json
{
  "type": "tool",
  "tool": "bash",
  "arguments": {
    "command": "git status"
  }
}
```

or

```json
{
  "type": "final",
  "answer": "task completed"
}
```

This removes ambiguity and keeps execution deterministic.

---

## 6. Architecture

### FSM Orchestrator

Simple states:

* Idle
* DoOne
* DoSecond
* Final
* Fail

No planning. No hidden reasoning.

---

### Engine

Execution loop:

* call LLM
* parse JSON
* execute tool
* continue

Maximum a few steps.

---

### State

Minimal:

* goal
* last result

No built-in memory system.

---

## 7. Performance

Whitebox is designed to be fast and cheap.

* written in Go
* minimal runtime overhead
* no background processes
* no heavy abstractions

Context control reduces token usage, which lowers LLM cost.

---

## 8. Observability (Langfuse)

Built-in tracing via Langfuse:

* full request trace
* inputs and outputs
* token usage
* execution flow

No hidden behavior. Everything is visible.

---

## 9. What You Get

* fast FSM execution
* strict and simple protocol
* explicit context control
* low resource usage
* low LLM cost
* no hidden logic

---

## 10. Use Case

Whitebox is a core, not a product.

You use it to build:

* CLI agents
* task-specific tools
* automation scripts
* internal assistants

Each agent is simple and focused.

---

## 11. Positioning

Whitebox can be seen as a “blaze agent” core:

* fast
* predictable
* minimal
* controllable

It trades autonomy for control and efficiency.

---

## 12. Getting Started

```bash
git clone https://github.com/Alice088/whitebox.git
cd whitebox
go mod download
```

```bash
cp .env.example .env
```

```bash
go run .
```

---

## Philosophy

* control over automation
* explicit over implicit
* simple over complex
