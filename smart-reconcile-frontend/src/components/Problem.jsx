import useReveal from '../useReveal'

const problems = [
  {
    source: 'M-Pesa',
    desc: (
      <>
        Returns <span className="highlight">TransID</span> and{' '}
        <span className="highlight">MSISDN</span> in a proprietary JSON envelope. Timestamps
        have <span className="highlight">no timezone</span>, which breaks naive date matching
        across markets.
      </>
    ),
  },
  {
    source: 'Stripe',
    desc: (
      <>
        Amounts are in the <span className="highlight">lowest currency denomination</span>.
        Webhooks can fire twice under load, creating silent duplicate charges in your ledger.
      </>
    ),
  },
  {
    source: 'Bank CSV',
    desc: (
      <>
        Exports have <span className="highlight">inconsistent date formats</span>, missing
        reference numbers, and occasional blank or malformed rows that cause silent parse
        failures.
      </>
    ),
  },
  {
    source: 'PayPal',
    desc: (
      <>
        Settlement reports <span className="highlight">lag 24–48 hours</span>.
        Cross-referencing delayed settlements against real-time internal records produces
        false mismatches on every run.
      </>
    ),
  },
]

export default function Problem() {
  const ref = useReveal()
  return (
    <section className="section reveal" ref={ref}>
      <div className="section-num">01 — The Problem</div>
      <h2>Providers don't agree with each other. Or themselves.</h2>
      <p className="section-lead">
        Every payment provider returns data in a different format. None of them are wrong —
        they just don't speak the same language. Without a proper abstraction layer, your
        reconciliation logic absorbs all that complexity.
      </p>
      <div className="problem-list">
        {problems.map((p) => (
          <div className="problem-item" key={p.source}>
            <span className="problem-source">{p.source}</span>
            <span className="problem-desc">{p.desc}</span>
          </div>
        ))}
      </div>
    </section>
  )
}
