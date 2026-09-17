package visualtest

import "github.com/chromedp/chromedp"

// pollBool polls a raw JavaScript predicate until it evaluates to true,
// capturing the outcome into held. Wrapping the predicate in Boolean(...) and
// capturing into a *bool is the only safe bool shape: a bare bool-valued
// predicate polled into a *string ALWAYS fails json unmarshal and the caller
// retries forever (the e2e lesson behind TODO #193) — this helper makes that
// mistake impossible to write. A timeout surfaces as the returned action's
// error, so err == nil implies held == true.
func pollBool(predicate string, held *bool, opts ...chromedp.PollOption) chromedp.Action {
	return chromedp.Poll("Boolean("+predicate+")", held, opts...)
}

// pollTrue is pollBool for callers that only need "did it hold" — the outcome
// lives inside the action, the error carries the verdict.
func pollTrue(predicate string, opts ...chromedp.PollOption) chromedp.Action {
	var held bool

	return pollBool(predicate, &held, opts...)
}

// pollText polls a JavaScript expression that must evaluate to a STRING until
// it returns a non-empty value, capturing it into dest (an empty result keeps
// polling, so "innerText of a region that does not exist yet" cannot
// false-pass). Never pass a boolean-valued expression here — wrap the
// condition in pollTrue instead.
func pollText(expression string, dest *string, opts ...chromedp.PollOption) chromedp.Action {
	return chromedp.Poll(expression, dest, opts...)
}
