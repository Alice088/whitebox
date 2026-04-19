<div align="center">
  <img src="./static/whitebox.png" alt="Whitebox" />
</div>

---

A minimal FSM-based core for building fast, controlled agents for OrcAI.

Whitebox is not a standalone assistant.  
It is a runtime designed to work as an execution unit inside the OrcAI orchestration system.

---

## 1. What is Whitebox

Whitebox is a lightweight agent core built around an FSM orchestrator.

It is designed to:

- execute tasks received from OrcAI
- keep behavior predictable
- avoid hidden complexity

Each instance acts as a worker that processes tasks of a specific type.

---

## 2. Problem

In orchestration systems, the main issue is not intelligence, but control.

Typical agents:

- accumulate uncontrolled context  
- mix responsibilities  
- behave inconsistently  
- consume unpredictable tokens  

Result:

- high API cost  
- unstable execution  
- hard-to-debug pipelines  

---

## 3. Core Idea

Whitebox is designed for OrcAI’s model:

- tasks come from a broker (`task.created`)  
- agents are stateless workers  
- behavior is strictly controlled  

Key principles:

- explicit context only  
- no hidden injections  
- deterministic execution  

The only built-in injection is the output protocol.

---

## 4. Context System

Context is simple and file-based.

```text
~/.whitebox-name/
  context/
    minds/
    memories/
    skills/
    tools/
    sessions/
  workspace/
````

How it works:

* all files in `context/` are loaded into the prompt
* nothing outside this directory exists for the model
* no implicit memory or hidden layers

This ensures:

* predictable inputs
* stable outputs
* controlled token usage

Context is fully owned by the system.

---

## 5. Output Protocol

The only built-in rule is a strict LLM output protocol.

The model must respond in valid JSON:

```json
{
  "type": "tool",
  "tool": "bash",
  "arguments": {
    "command": "git status"
  }
}
````

or

```json
{
  "type": "final",
  "answer": "task completed"
}
```

This protocol is used by the engine to:

* parse model output deterministically
* decide whether to execute a tool or finish
* keep execution predictable

It is not related to OrcAI routing.
OrcAI operates on task/event level, while this protocol is strictly between the LLM and the local engine.

---

## 6. Architecture

### FSM Orchestrator

States:

* Idle
* DoOne
* DoSecond
* Final
* Fail

No planning. No hidden reasoning.

---

### Engine

Execution loop:

* receives task input
* calls LLM
* parses JSON
* executes tool
* publishes results back

Integrated with broker:

* consumes `task.created`
* emits `task.logs`, `task.result`, `task.error`

---

### State

Minimal runtime state:

* goal
* last result

No persistent memory by default.

---

## 7. Performance

Whitebox is optimized for orchestration workloads.

* written in Go
* low memory usage
* fast startup and execution
* no background overhead

Strict context control reduces token usage and cost.

---

## 8. Observability (Langfuse)

Built-in support for Langfuse:

* full trace per task
* LLM inputs and outputs
* token usage
* execution steps

This makes every task observable inside OrcAI.

---

## 9. What You Get

* FSM-based execution worker
* strict tool/final protocol
* explicit context control
* NATS-compatible task processing
* low resource usage
* predictable behavior

---

## 10. Use Case

Whitebox is used inside OrcAI to run agents.

Each agent:

* subscribes to `task.created.*`
* filters by `type`
* processes only its tasks
* publishes results back

Typical usage:

* git agents
* file system agents
* code execution agents
* domain-specific workers

---

## 11. Positioning

Whitebox is a “blaze agent” core for OrcAI:

* fast
* minimal
* controllable
* cheap to run

It is designed for execution, not autonomy.

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

* execution over reasoning
* control over automation
* explicit over implicit
* simple over complex


