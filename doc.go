// Package monobank is a Go client for the monobank REST API
// (https://api.monobank.ua/docs/).
//
// It covers three authorization modes:
//
//   - Public — no auth; see [PublicAPI].
//   - Personal — a single user's token; see [PersonalAPI].
//   - Corporate / providers — service-level access via an ECDSA key pair;
//     see [CorporateAPI].
//
// # Webhooks
//
// Personal and corporate clients can subscribe a URL to receive statement
// events as signed POST requests. Verify the signature before trusting the body:
//
//	sk, _ := client.ServerKey(ctx)               // cache; refresh on rotation
//	err := sk.Verify(body, r.Header.Get("X-Sign"))
//	event, _ := monobank.ParseWebHook(body)
//
// See examples/webhook for a complete handler.
package monobank
