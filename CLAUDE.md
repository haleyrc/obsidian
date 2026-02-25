# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Go CLI for managing personal Obsidian vaults. Currently implements a `movies` command that integrates with The Movie Database (TMDB) API to generate Obsidian notes with metadata, cast links, genre links, and poster images.

## Build and Run

The commands in this project contact live services and should only be built and run by the user.

## Dependencies

- `github.com/haleyrc/lib` — shared library (TMDB client wrapper, CLI framework)
- `gopkg.in/yaml.v3` — YAML frontmatter encoding

## Architecture

**CLI layer** (`cmd/movies/`) uses a generic `CLI[App]` from `lib/cli`. `main.go` wires up dependencies into an `App` struct (carrying a `TMDBClient`) and registers subcommands. `new.go` implements the `new` subcommand handler.

**Domain layer** (`internal/movies/`) contains:
- `Movie`/`MovieNote` structs and template rendering via embedded `movie.tmpl`
- `Vault` for filesystem operations (writing notes, downloading posters)
- `TMDBClient` interface for dependency injection

**Integration layer** (`internal/tmdb/`) wraps the lib TMDB client. `internal/http/` provides a download utility. `internal/obsidian/` has Obsidian-specific helpers (filename sanitization, `Link` type).

**Vault directory structure** (expected to exist on disk):
```
vault-root/
├── movies/        # Movie markdown notes
├── actors/        # Actor markdown notes
├── genres/        # Genre markdown notes
├── attachments/   # Downloaded poster images
└── reviews/       # Review embed references
```

**Cast filtering**: actors are included if order < 10 OR popularity >= 1.0 OR already known in the vault.
