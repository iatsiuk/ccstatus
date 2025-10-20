# ccstatus

Real-time context usage monitor for Claude Code status line.

## What is this?

`ccstatus` is a lightweight utility that displays your current Claude Code session's token usage directly in the status line. It helps you track how much of your context window is being used during conversations.

## Features

- Real-time context usage display: `[ctx: 58173/200000 29.1%] claude-sonnet-4-5-20250929`
- Color-coded warnings based on usage level:
  - Green (0-60%) - plenty of space
  - Yellow (60-80%) - approaching limit
  - Red (80-100%) - context nearly full
- Fast and lightweight (single binary, <3MB)
- Zero dependencies
- Handles empty/new sessions gracefully

## Algorithm

Context calculation (corrected after codex-expert review):
```
current_context = input_tokens + cache_read_input_tokens
percentage = (current_context / model_limit) * 100
```

**Why this formula:**
- `input_tokens`: all non-cached tokens (system prompts, user requests, tool calls)
- `cache_read_input_tokens`: all cached tokens being read from cache
- `cache_creation_input_tokens`: NOT included (represents cache write operations, not context usage)

**Model-specific limits:**
- Claude 3 Opus/Sonnet/Haiku: 200K tokens
- Claude Sonnet 4/4.5: 200K tokens
- Claude 2: 100K tokens
- Claude Instant: 100K tokens
- Unknown models: 200K default

Based on extended reasoning analysis by codex-expert, fixing issues found in ccstatusline#100.

## Installation

### 1. Build from source

```bash
git clone <repository-url> ccstatus
cd ccstatus
go build -o ccstatus
```

### 2. Configure Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "/absolute/path/to/ccstatus/ccstatus"
  }
}
```

Replace `/absolute/path/to/ccstatus/ccstatus` with the actual path to your built binary.

### 3. Restart Claude Code

The status line will appear at the bottom of your Claude Code interface.

## Usage

Once configured, ccstatus runs automatically. No manual interaction needed.

**Example output:**
```
[ctx: 98882/200000 49.4%] claude-sonnet-4-5-20250929
```

- `98882/200000` - current tokens / maximum tokens
- `49.4%` - percentage of context used
- `claude-sonnet-4-5-20250929` - model identifier

## Development

### Project structure

```
ccstatus/
├── main.go                    # entry point
├── internal/
│   ├── parser/               # JSONL transcript parsing
│   │   ├── parser.go
│   │   └── parser_test.go
│   ├── calculator/           # context calculation
│   │   ├── calculator.go
│   │   └── calculator_test.go
│   └── formatter/            # colored output formatting
│       ├── formatter.go
│       └── formatter_test.go
└── testdata/                 # test fixtures
```

### Running tests

```bash
go test ./... -v
```

### Test coverage

100% coverage for all packages:
- parser: 100%
- calculator: 100%
- formatter: 100%

## How it works

1. **Claude Code invokes ccstatus** and passes session info via stdin:
   ```json
   {
     "session_id": "af99e13e-377a-4064-ae40-3987bc91cdee",
     "cwd": "/path/to/project",
     "model": {
       "id": "claude-sonnet-4-5-20250929",
       "display_name": "Sonnet 4.5"
     },
     "transcript_path": "/path/to/session.jsonl"
   }
   ```

2. **ccstatus reads the transcript** JSONL file and finds the last message with usage data

3. **Extracts token counts** from the usage field:
   ```json
   {
     "message": {
       "role": "assistant",
       "usage": {
         "input_tokens": 9,
         "cache_read_input_tokens": 58164,
         "cache_creation_input_tokens": 1097
       }
     }
   }
   ```

4. **Calculates context usage:**
   - Formula: `input_tokens + cache_read_input_tokens`
   - Example: `9 + 58164 = 58173 tokens`
   - Percentage: `58173 / 200000 = 29.1%`

5. **Outputs formatted result** to stdout with ANSI color codes

6. **Claude Code displays** the output in its status line

## License

MIT
