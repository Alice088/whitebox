```
================================================================================
                                coreClaw
================================================================================

An AI coding assistant focused on context control and observability.

This is a chatbot, but not in the usual sense. The goal is not just to generate
answers, but to explicitly control what the model sees and how it behaves.

--------------------------------------------------------------------------------
PROBLEM
--------------------------------------------------------------------------------

Most modern AI agents do not control context properly.

- Context grows uncontrollably
- Irrelevant data is injected into prompts
- It is unclear what the model actually sees
- Results are hard to reproduce
- Token usage and cost become unpredictable

The result is a black box system:
it produces output, but you cannot understand or control it.

--------------------------------------------------------------------------------
SOLUTION
--------------------------------------------------------------------------------

coreClaw is built around explicit context management.

- You decide what goes into the prompt
- Context is structured and bounded
- Every request is traceable
- Model behavior becomes predictable

Not a "smart agent", but a controlled system.

--------------------------------------------------------------------------------
FEATURES (PLANNED)
--------------------------------------------------------------------------------

- Explicit context manager
- Token limits and prioritization
- Manual file injection into context
- TUI interface
- File tools (read / write)
- Logging and tracing via Langfuse

--------------------------------------------------------------------------------
ARCHITECTURE (HIGH LEVEL)
--------------------------------------------------------------------------------

Pipeline:

user -> orchestrator -> context manager -> prompt builder -> LLM
     -> tools -> response -> logging

Core components:

- Model layer
- Prompt builder
- Context manager
- Tools
- Observability

--------------------------------------------------------------------------------
TECH STACK
--------------------------------------------------------------------------------

- Go
- Bubbletea (TUI)
- DeepSeek / OpenAI-compatible APIs
- SQLite / local storage
- Langfuse

--------------------------------------------------------------------------------
CURRENT STATUS
--------------------------------------------------------------------------------

Idea + early implementation.

Next step:
minimal CLI with prompt control and logging.

--------------------------------------------------------------------------------
ROADMAP (SHORT)
--------------------------------------------------------------------------------

1. MVP
   - CLI
   - Basic LLM call
   - Logging

2. Context
   - Token limits
   - Prioritization
   - Trimming

3. Improvements
   - Memory
   - Stability

4. Expansion
   - Plugins
   - Model abstraction

--------------------------------------------------------------------------------
PHILOSOPHY
--------------------------------------------------------------------------------

coreClaw is not about building a better AI.

It is about control:

- Control of input
- Control of context
- Control of cost
- Control of behavior

If you do not control the context,
you do not control the system.
================================================================================
```
