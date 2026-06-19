export default function Hero() {
  return (
    <section className="hero">
      <div className="hero-eyebrow">Financial Infrastructure</div>
      <h1>
        Your ledger is off.
        <br />
        You just don't know
        <br />
        <em>why.</em>
      </h1>
      <p className="hero-sub">
        smart-reconcile is a backend platform that ingests transactions from multiple payment
        providers, normalizes everything into a unified format, and runs a reconciliation
        engine that detects discrepancies automatically — with a full audit trail explaining
        every decision.
      </p>
      <div className="hero-meta">
        <div className="meta-item">
          <span className="label">Status</span>
          <span className="value red">Under active development</span>
        </div>
        <div className="meta-item">
          <span className="label">Built for</span>
          <span className="value">Teams that can't afford to be wrong</span>
        </div>
      </div>
    </section>
  )
}
