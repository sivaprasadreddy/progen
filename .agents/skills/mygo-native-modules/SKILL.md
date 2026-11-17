---
name: mygo-native-modules
description: Write and change native modules of MyGo apps, Go interfaces implemented in Swift (macOS frameworks, packages and views), Rust or Zig (crates and C libraries, every platform). Use when a MyGo app needs an API Go lacks, or when a package has a file ending in _native.go or a swift/, rust/ or zig/ directory beside it.
---

# MyGo native modules

A native module is a Go package whose interfaces Swift, Rust or Zig
implement. It is experimental: after upgrading MyGo, run `mygo native
check` and follow the release notes, as the generated code and commands
may change. Read `docs/native-modules.md` of the app's MyGo (`go list -m -f
'{{.Dir}}' github.com/egoist/mygo`) for the whole contract; this is the
loop.

## Make one

```sh
mygo native new <name>                  # modules/<name>: the Go interface, swift/ in Swift
mygo native new -lang rust,zig <name>   # or in Rust and Zig, for every platform
mygo native add rust [dir]              # implement an existing module in another language
```

Pick Swift for what macOS offers only in Swift (its frameworks and views),
and Rust or Zig for a crate or a C library on every platform. A module may
have several: macOS uses Swift when there is one, other platforms Rust,
then Zig.

## The loop

1. Change the Go interface in `<name>.go`: the contract. `//mygo:native`
   marks an entry point (the package gets `New<T>()`), `//mygo:mainthread`
   runs methods on the main thread, embedding `native.View` makes the
   implementation an `NSView` that `ui.NativeView(c, v)` shows (Swift
   only), and `native.NewError("code")` declares errors the implementation
   returns as `<Module>Error`.
2. Run `mygo native check [-lang <language>] [dir]`. It regenerates
   `<name>_native.go` and each language's generated file
   (`swift/Sources/<Module>/<Module>.mygo.swift`, `rust/src/mygo.rs`,
   `zig/src/mygo.zig`: never edit those), builds the implementation and
   vets Go, and prints every problem as `file:line: message` with what to
   write instead. Fix them all; repeat until it says the module is fine.
3. Implement it:
   - Swift, in `swift/Sources/<Module>/<Module>.swift`: a factory
     `make<T>()` per entry point and classes conforming to the generated
     protocols.
   - Rust, in `rust/src/lib.rs` (with `mod mygo; use mygo::*;`): a
     factory `make_<t>() -> Result<Arc<dyn T>, Error>` and types
     implementing the generated traits.
   - Zig, in `zig/src/<name>.zig` (with `const mygo = @import("mygo.zig");`):
     a factory `make<T>() mygo.Error!mygo.T`, returning `mygo.T.of(ptr)`
     of a struct with a `pub fn` per method, taking `arena: Allocator`
     after `self`. Results may live in the arena; what outlives the call
     uses `mygo.allocator` and is freed in `deinit`.
   The compiler names any method missing.
4. Try a method without the app:
   `mygo native call [dir] <Interface>.<Method> '{"param": value}'`
   prints the result as JSON, a line per value of a stream, and the calls
   of callbacks.
5. Test from Go: `go test` in the module runs the implementation once
   `mygo dev` or `mygo native check` built it;
   `mygo native build -target host -lang <language> [dir]` builds another
   of its languages to test that one. Test code that takes the interface
   with a fake of it.

## What crosses

- Parameters, results and fields: numbers, `bool`, `string`, `[]byte`,
  `time.Time`, `time.Duration`, slices, `map[string]T`, pointers, structs
  and enums (`type T string` with constants) of the package.
- Methods return nothing, a value, an error, both, `iter.Seq[T]` or
  `iter.Seq2[T, error]` (taking a `context.Context` first), or an
  `io.Closer` that stops what they started.
- `context.Context` first makes the Swift method `async`, and runs Rust's
  and Zig's on a thread with a `Context` saying when Go canceled it.
  Function parameters are closures the implementation may keep: without
  results they run later, in order, on Go's main thread.
- An interface the module returns is the implementation's (an object,
  closed with `io.Closer`); one it only takes is Go's, which it calls.
- `*native.Buffer` shares memory without copying. `*mygo.Window`,
  `ui.Rect` anchors, views, textures and other modules' interfaces are
  Swift only.

Keep each language's own types (enums with payloads, generics, overloads)
inside the implementation. Put usage descriptions in `swift/Info.plist`
and entitlements in `swift/Entitlements.plist`; dependencies in
`swift/Package.swift` (keeping its product dynamic and named after the
module), `rust/Cargo.toml` (keeping `crate-type = ["cdylib"]`) or
`zig/build.zig.zon`.
