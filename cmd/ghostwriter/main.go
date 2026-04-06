package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	colorReset = "\033[0m"
	colorDim   = "\033[2m"
)

// ─────────────────────────────────────────────

// Speed profiles
// ─────────────────────────────────────────────

type speedProfile struct {
	base    int // base ms per char
	jitter  int
	pauseMs int // extra pause at newlines/brackets
}

var speeds = map[string]speedProfile{
	"slow":   {80, 40, 400},
	"medium": {45, 25, 250},
	"fast":   {20, 12, 120},
	"turbo":  {8, 5, 50},
}

// ─────────────────────────────────────────────
// Snippets
// ─────────────────────────────────────────────

type snippet struct {
	lang string
	code string
}

var snippets = []snippet{
	// ── Go ──────────────────────────────────────────────────────────────────
	{lang: "go", code: `package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr    string
	timeout time.Duration
	mux     *http.ServeMux
}

func NewServer(addr string) *Server {
	return &Server{
		addr:    addr,
		timeout: 30 * time.Second,
		mux:     http.NewServeMux(),
	}
}

func (s *Server) Register(pattern string, h http.HandlerFunc) {
	s.mux.HandleFunc(pattern, h)
}

func (s *Server) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:         s.addr,
		Handler:      s.mux,
		ReadTimeout:  s.timeout,
		WriteTimeout: s.timeout,
	}
	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()
	log.Printf("listening on %s", s.addr)
	return srv.ListenAndServe()
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "{\"status\":\"ok\"}")
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := NewServer(":8080")
	s.Register("/health", healthHandler)
	if err := s.Start(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}`},
	{lang: "go", code: `package cache

import (
	"sync"
	"time"
)

type entry[V any] struct {
	val     V
	expires time.Time
}

type TTLCache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]entry[V]
	ttl   time.Duration
}

func New[K comparable, V any](ttl time.Duration) *TTLCache[K, V] {
	c := &TTLCache[K, V]{items: make(map[K]entry[V]), ttl: ttl}
	go c.evict()
	return c
}

func (c *TTLCache[K, V]) Set(k K, v V) {
	c.mu.Lock()
	c.items[k] = entry[V]{val: v, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

func (c *TTLCache[K, V]) Get(k K) (V, bool) {
	c.mu.RLock()
	e, ok := c.items[k]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expires) {
		var zero V
		return zero, false
	}
	return e.val, true
}

func (c *TTLCache[K, V]) evict() {
	ticker := time.NewTicker(c.ttl / 2)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, e := range c.items {
			if now.After(e.expires) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}`},

	// ── Python ──────────────────────────────────────────────────────────────
	{lang: "python", code: `import asyncio
import logging
from dataclasses import dataclass, field
from typing import AsyncIterator

logger = logging.getLogger(__name__)


@dataclass
class Config:
    host: str = "localhost"
    port: int = 8080
    workers: int = 4
    debug: bool = False
    allowed_origins: list[str] = field(default_factory=list)


class EventBus:
    def __init__(self) -> None:
        self._listeners: dict[str, list] = {}

    def subscribe(self, event: str):
        def decorator(fn):
            self._listeners.setdefault(event, []).append(fn)
            return fn
        return decorator

    async def emit(self, event: str, **kwargs) -> None:
        for fn in self._listeners.get(event, []):
            try:
                await fn(**kwargs)
            except Exception as exc:
                logger.exception("listener %s raised: %s", fn.__name__, exc)


async def stream_chunks(data: bytes, chunk_size: int = 4096) -> AsyncIterator[bytes]:
    offset = 0
    while offset < len(data):
        yield data[offset : offset + chunk_size]
        offset += chunk_size
        await asyncio.sleep(0)


bus = EventBus()


@bus.subscribe("user.login")
async def on_login(user_id: int, ip: str) -> None:
    logger.info("user %d logged in from %s", user_id, ip)


if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)
    asyncio.run(bus.emit("user.login", user_id=42, ip="127.0.0.1"))`},
	{lang: "python", code: `from __future__ import annotations

import hashlib
import json
import os
from pathlib import Path
from typing import Any


class ConfigStore:
    """Persistent JSON config with dot-notation access."""

    def __init__(self, path: str | Path) -> None:
        self._path = Path(path)
        self._data: dict[str, Any] = {}
        self._load()

    def _load(self) -> None:
        if self._path.exists():
            self._data = json.loads(self._path.read_text())

    def save(self) -> None:
        self._path.parent.mkdir(parents=True, exist_ok=True)
        self._path.write_text(json.dumps(self._data, indent=2))

    def get(self, key: str, default: Any = None) -> Any:
        parts = key.split(".")
        node = self._data
        for p in parts:
            if not isinstance(node, dict) or p not in node:
                return default
            node = node[p]
        return node

    def set(self, key: str, value: Any) -> None:
        parts = key.split(".")
        node = self._data
        for p in parts[:-1]:
            node = node.setdefault(p, {})
        node[parts[-1]] = value


if __name__ == "__main__":
    cfg = ConfigStore(os.path.expanduser("~/.config/myapp/settings.json"))
    cfg.set("server.host", "0.0.0.0")
    cfg.set("server.port", 9000)
    cfg.set("features.dark_mode", True)
    cfg.save()
    print(cfg.get("server.port"))  # 9000`},

	// ── TypeScript ──────────────────────────────────────────────────────────
	{lang: "typescript", code: `import { z } from "zod";

const UserSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  role: z.enum(["admin", "user", "guest"]),
  createdAt: z.coerce.date(),
  meta: z.record(z.unknown()).default({}),
});

type User = z.infer<typeof UserSchema>;

class Repository<T extends { id: string }> {
  private store = new Map<string, T>();

  upsert(item: T): void { this.store.set(item.id, item); }

  findById(id: string): T | undefined { return this.store.get(id); }

  findAll(predicate?: (item: T) => boolean): T[] {
    const items = [...this.store.values()];
    return predicate ? items.filter(predicate) : items;
  }

  delete(id: string): boolean { return this.store.delete(id); }
}

const userRepo = new Repository<User>();
const raw = {
  id: "550e8400-e29b-41d4-a716-446655440000",
  email: "dev@example.com",
  role: "admin",
  createdAt: "2026-01-01",
};
const user = UserSchema.parse(raw);
userRepo.upsert(user);
const admins = userRepo.findAll((u) => u.role === "admin");
console.log("admins:", admins.length);`},

	// ── JavaScript ──────────────────────────────────────────────────────────
	{lang: "javascript", code: `async function withRetry(fn, { retries = 3, baseDelay = 300, factor = 2 } = {}) {
  let attempt = 0;
  while (true) {
    try {
      return await fn();
    } catch (err) {
      if (attempt >= retries) throw err;
      const delay = baseDelay * Math.pow(factor, attempt) + Math.random() * 100;
      console.warn("attempt " + (attempt + 1) + " failed, retrying in " + delay.toFixed(0) + "ms");
      await new Promise((r) => setTimeout(r, delay));
      attempt++;
    }
  }
}

class EventEmitter {
  #handlers = new Map();

  on(event, handler) {
    const list = this.#handlers.get(event) ?? [];
    this.#handlers.set(event, [...list, handler]);
    return () => this.off(event, handler);
  }

  off(event, handler) {
    const list = this.#handlers.get(event) ?? [];
    this.#handlers.set(event, list.filter((h) => h !== handler));
  }

  emit(event, ...args) {
    for (const h of this.#handlers.get(event) ?? []) h(...args);
  }
}

const bus = new EventEmitter();
bus.on("data", (payload) => console.log("received:", payload));

withRetry(() => fetch("https://api.example.com/v1/items"))
  .then((r) => r.json())
  .then((data) => bus.emit("data", data))
  .catch(console.error);`},

	// ── Rust ────────────────────────────────────────────────────────────────
	{lang: "rust", code: `use std::collections::HashMap;
use std::sync::{Arc, Mutex};
use std::time::{Duration, Instant};

#[derive(Debug, Clone)]
struct CacheEntry<V> {
    value: V,
    expires: Instant,
}

pub struct Cache<K, V> {
    inner: Arc<Mutex<HashMap<K, CacheEntry<V>>>>,
    ttl: Duration,
}

impl<K: std::hash::Hash + Eq + Clone, V: Clone> Cache<K, V> {
    pub fn new(ttl: Duration) -> Self {
        Self { inner: Arc::new(Mutex::new(HashMap::new())), ttl }
    }

    pub fn insert(&self, key: K, value: V) {
        let entry = CacheEntry { value, expires: Instant::now() + self.ttl };
        self.inner.lock().unwrap().insert(key, entry);
    }

    pub fn get(&self, key: &K) -> Option<V> {
        let guard = self.inner.lock().unwrap();
        guard.get(key).and_then(|e| {
            if Instant::now() < e.expires { Some(e.value.clone()) } else { None }
        })
    }

    pub fn evict_expired(&self) {
        let now = Instant::now();
        self.inner.lock().unwrap().retain(|_, e| e.expires > now);
    }
}

fn main() {
    let cache: Cache<String, u64> = Cache::new(Duration::from_secs(10));
    cache.insert("requests".to_string(), 42);
    match cache.get(&"requests".to_string()) {
        Some(v) => println!("requests: {}", v),
        None    => println!("expired"),
    }
}`},

	// ── Bash ────────────────────────────────────────────────────────────────
	{lang: "bash", code: `#!/usr/bin/env bash
set -euo pipefail

readonly LOG_FILE="/tmp/deploy_$(date +%Y%m%d_%H%M%S).log"
readonly DEPLOY_DIR="/var/www/app"
readonly SERVICE_NAME="myapp"

log() { echo "[$(date '+%H:%M:%S')] $*" | tee -a "$LOG_FILE"; }
die() { log "ERROR: $*"; exit 1; }
require() { command -v "$1" >/dev/null 2>&1 || die "required: $1"; }

require git
require systemctl
require rsync

BRANCH="${1:-main}"
log "Deploying branch: $BRANCH"

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

log "Cloning..."
git clone --depth 1 --branch "$BRANCH" "https://github.com/org/repo.git" "$TMPDIR" >>"$LOG_FILE" 2>&1

log "Running tests..."
( cd "$TMPDIR" && ./scripts/test.sh ) >>"$LOG_FILE" 2>&1 || die "tests failed"

log "Syncing files..."
rsync -av --delete --exclude='.git' "$TMPDIR/" "$DEPLOY_DIR/" >>"$LOG_FILE" 2>&1

log "Restarting service..."
systemctl restart "$SERVICE_NAME" || die "restart failed"
systemctl is-active --quiet "$SERVICE_NAME" || die "service not running"
log "Deploy complete."`},
}

// ─────────────────────────────────────────────
// Typer
// ─────────────────────────────────────────────

func typeChar(ch byte, sp speedProfile, rng *rand.Rand) {
	os.Stdout.Write([]byte{ch})

	d := sp.base + rng.Intn(sp.jitter*2+1) - sp.jitter
	if d < 1 {
		d = 1
	}
	switch ch {
	case '\n':
		d += sp.pauseMs + rng.Intn(sp.pauseMs/2+1)
	case '{', '}', '(', ')', ';', ':':
		d += sp.pauseMs / 3
	}
	time.Sleep(time.Duration(d) * time.Millisecond)
}

func typeSnippet(code string, sp speedProfile, rng *rand.Rand, stop <-chan struct{}) {
	for i := 0; i < len(code); i++ {
		select {
		case <-stop:
			return
		default:
		}

		ch := code[i]
		// ~2.5% typo chance
		if rng.Intn(40) == 0 {
			wrong := byte('a' + rng.Intn(26))
			typeChar(wrong, sp, rng)
			time.Sleep(80 * time.Millisecond)
			os.Stdout.Write([]byte{8, 32, 8}) // backspace
		}
		typeChar(ch, sp, rng)

		// Thinking pause before keywords
		if i+2 < len(code) {
			ahead := code[i : i+2]
			switch ahead {
			case "//", "/*", "fn", "fu", "de":
				if rng.Intn(3) == 0 {
					time.Sleep(time.Duration(500+rng.Intn(700)) * time.Millisecond)
				}
			}
		}
	}
}

// ─────────────────────────────────────────────
// main
// ─────────────────────────────────────────────

func main() {
	lang := ""
	speedName := "medium"

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--lang":
			if i+1 < len(args) {
				lang = strings.ToLower(args[i+1])
				i++
			}
		case "--speed":
			if i+1 < len(args) {
				speedName = strings.ToLower(args[i+1])
				i++
			}
		case "--help", "-h":
			printHelp()
			return
		}
	}

	sp, ok := speeds[speedName]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown speed: %s (slow/medium/fast/turbo)\n", speedName)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	pool := snippets
	if lang != "" {
		var filtered []snippet
		for _, s := range snippets {
			if strings.HasPrefix(s.lang, lang) {
				filtered = append(filtered, s)
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "No snippets for lang: %s\nAvailable: go, python, typescript, javascript, rust, bash\n", lang)
			os.Exit(1)
		}
		pool = filtered
	}

	stop := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		close(stop)
		fmt.Printf("\n\n%s[ghostwriter] stopped.%s\n", colorDim, colorReset)
		os.Exit(0)
	}()

	fmt.Printf("%s[ghostwriter] running — Ctrl+C to stop%s\n\n", colorDim, colorReset)
	time.Sleep(800 * time.Millisecond)

	for {
		s := pool[rng.Intn(len(pool))]
		ext := langToExt(s.lang)
		name := randFilename(s.lang, rng)
		fmt.Printf("\n%s# %s.%s%s\n\n", colorDim, name, ext, colorReset)
		time.Sleep(time.Duration(300+rng.Intn(400)) * time.Millisecond)

		typeSnippet(s.code, sp, rng, stop)

		pause := 1200 + rng.Intn(2000)
		select {
		case <-stop:
			return
		case <-time.After(time.Duration(pause) * time.Millisecond):
		}
	}
}

func langToExt(lang string) string {
	switch lang {
	case "go":
		return "go"
	case "python":
		return "py"
	case "typescript":
		return "ts"
	case "javascript":
		return "js"
	case "rust":
		return "rs"
	case "bash":
		return "sh"
	}
	return "txt"
}

var filenameParts = map[string][]string{
	"go":         {"server", "cache", "handler", "router", "middleware", "config", "store", "client", "worker"},
	"python":     {"main", "utils", "models", "service", "config", "client", "helpers", "parser"},
	"typescript": {"types", "api", "store", "hooks", "utils", "schema", "client"},
	"javascript": {"index", "utils", "api", "events", "helpers", "client"},
	"rust":       {"lib", "cache", "server", "client", "handler", "config"},
	"bash":       {"deploy", "setup", "install", "build", "release", "backup"},
}

func randFilename(lang string, rng *rand.Rand) string {
	parts, ok := filenameParts[lang]
	if !ok {
		return "main"
	}
	return parts[rng.Intn(len(parts))]
}

func printHelp() {
	fmt.Printf("\nghostwriter — пишет код сам по себе\n\n")
	fmt.Printf("  ghostwriter                    случайный язык\n")
	fmt.Printf("  ghostwriter --lang go          только Go\n")
	fmt.Printf("  ghostwriter --lang python      только Python\n")
	fmt.Printf("  ghostwriter --lang typescript  только TypeScript\n")
	fmt.Printf("  ghostwriter --lang javascript  только JavaScript\n")
	fmt.Printf("  ghostwriter --lang rust        только Rust\n")
	fmt.Printf("  ghostwriter --lang bash        только Bash\n")
	fmt.Printf("  ghostwriter --speed fast       slow/medium/fast/turbo\n")
	fmt.Printf("  Ctrl+C                         остановить\n\n")
}
