// Package wire is the transport-agnostic wiring contract shared by the HTMX
// and Datastar integrations. One typed Action describes a client-initiated
// hypermedia exchange (method, URL, triggering event, target region);
// Attributes renders it as the attribute dialect of the configured transport.
//
// # The one-object rule
//
// All Datastar fetch options travel in ONE options object literal —
// {selector: '…', mode: '…', contentType: 'form'} — never several. The pinned
// runtime's action dispatcher invokes the action with spread arguments and
// reads exactly one options object, so a second object argument is silently
// ignored (the earlier multi-object shape dropped form encoding whenever a
// Selector was also set). datastarActionExpr upholds the rule; any new option
// must merge into the same literal.
//
// # URL templates
//
// URLs are static strings, but components that substitute per-render values
// (forms.Calendar's month navigation) clone the consumer's Action and rewrite
// {year}/{month} placeholders in the URL before rendering — see
// forms/calendar_nav.go. Consumers can adopt the same convention for their
// own per-item wiring: keep the Action static in props, expand placeholders
// at render time. Placeholders are a CONVENTION, not parsed by this package —
// an unexpanded "{year}" renders verbatim.
//
// # The common subset
//
// The contract covers only the dialects' common subset (ADR-0036), extended
// twice by explicit ADR: ContentType/DebounceMS (ADR-0038) and the typed
// interval/reveal trigger language (ADR-0043). Transport-specific machinery
// stays in the htmx and datastar modules. See docs/transport-wiring.md for
// the full consumer guide.
//
// # Servers
//
// Handler turns one endpoint into a both-transports endpoint: Datastar
// callers get response-header targeting, htmx and plain callers pass through.
// DecodeForm decodes the request body (or query) into a typed struct for
// either dialect.
package wire
