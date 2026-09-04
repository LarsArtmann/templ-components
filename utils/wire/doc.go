// Package wire is the transport-agnostic wiring contract shared by the HTMX
// and Datastar integrations. One typed Action describes a client-initiated
// hypermedia exchange (method, URL, triggering event, target region);
// Attributes renders it as the attribute dialect of the configured transport.
//
// The contract composes with every component in the library today: the
// rendered attributes spread into BaseProps.Attrs, and components that opt in
// (display.Button) take an Action directly via a Wire field. On the server,
// Handler turns one endpoint into a both-transports endpoint: Datastar
// callers get response-header targeting, htmx and plain callers pass through.
//
// The contract covers only the dialects' common subset (ADR-0036);
// transport-specific machinery stays in the htmx and datastar modules. See
// docs/transport-wiring.md for the full consumer guide.
package wire
