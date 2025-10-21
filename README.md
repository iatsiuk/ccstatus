# ccstatus

Real-time context usage monitor for Claude Code status line.

## What is this?

`ccstatus` is a lightweight utility that displays your current Claude Code session's token usage directly in the status line. It helps you track how much of your context window is being used during conversations.

## Known Issues

Status line may not appear immediately on Claude Code startup. It will display after the first message exchange or when pressing Shift+Tab, as Claude Code updates the status line only when conversation messages update.

## Installation

### 1. Build from source

```bash
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
